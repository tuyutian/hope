package products

import (
	"context"
	"fmt"
	"time"

	"backend/internal/domain/entity/jobs"
	productEntity "backend/internal/domain/entity/products"
	cartRepo "backend/internal/domain/repo/carts"
	jobRepo "backend/internal/domain/repo/jobs"
	productRepo "backend/internal/domain/repo/products"
	"backend/internal/domain/repo/shopifys"
	userRepo "backend/internal/domain/repo/users"
	"backend/internal/providers"
	"backend/pkg/logger"
)

type ProductService struct {
	productRepo        productRepo.ProductRepository
	jobProductRepo     jobRepo.ProductRepository
	VariantRepo        productRepo.VariantRepository
	userRepo           userRepo.UserRepository
	cartSettingRepo    cartRepo.CartSettingRepository
	asynqRepo          jobRepo.AsynqRepository
	productGraphqlRepo shopifys.ProductGraphqlRepository
}

func NewProductService(
	repos *providers.Repositories) *ProductService {
	return &ProductService{
		productRepo:        repos.ProductRepo,
		VariantRepo:        repos.VariantRepo,    // 添加初始化
		jobProductRepo:     repos.JobProductRepo, // 添加初始化
		userRepo:           repos.UserRepo,
		cartSettingRepo:    repos.CartSettingRepo,
		asynqRepo:          repos.AsyncRepo,
		productGraphqlRepo: repos.ProductGraphqlRepo,
	}
}

func (p *ProductService) UploadProduct(ctx context.Context, req *productEntity.ProductReq) error {
	// 查询产品
	product, err := p.productRepo.First(ctx, req.UserID)

	if err != nil {
		logger.Error(ctx, "upload-product-db异常", "Err:", err.Error())
		return err
	}

	var shopifyProductId int64
	var userProductId int64
	if product == nil {
		// 创建产品
		userProductId, err = p.productRepo.CreateProduct(ctx, &productEntity.UserProduct{
			UserID:      req.UserID,
			Title:       "Protectify",
			Vendor:      "Protectify",
			Tags:        "Protectify",
			Description: "Protectify",
			Option1:     "Title",
			ImageUrl:    "https://img0.baidu.com/it/u=1868319137,207061070&fm=253&fmt=auto?w=1431&h=800",
		})

		// 创建变体
		var variants []*productEntity.UserVariant

		price := 1.0
		for i := 0; i < 100; i++ {
			variant := &productEntity.UserVariant{
				// 根据你的 productEntity.UserVariants 结构体，填充字段
				UserID:        req.UserID,
				UserProductId: userProductId,
				SkuName:       fmt.Sprintf("%.2f", price),
				Sku1:          fmt.Sprintf("ZIK%d", i),
				Price:         price,
			}
			price = price + 1.01
			variants = append(variants, variant)
		}

		// 调用 DAO 创建方法
		if err = p.VariantRepo.CreateVariants(ctx, variants); err != nil {
			// 处理错误
			logger.Error(ctx, "upload-product-db(2)异常", "Err:", err.Error())
			return err
		}

		shopifyProductId = 0
	} else {
		userProductId = product.Id

		if product.ProductId == 0 {
			shopifyProductId = 0
		} else {
			shopifyProductId = product.ProductId
		}
	}

	jobId, err := p.jobProductRepo.Create(ctx, &jobs.JobProduct{
		UserID:        req.UserID,
		UserProductId: userProductId,
	})
	if err != nil {
		logger.Error(ctx, "upload-product-db(3)异常", "Err:", err.Error())
		return err
	}

	_, err = p.asynqRepo.NewProductTask(ctx, jobId, userProductId, shopifyProductId)
	if err != nil {
		logger.Error(ctx, "upload-product-db推送asynq异常", "Err:", err.Error())
		return err
	}

	return nil
}

func (p *ProductService) ProductUpdate(ctx context.Context, req productEntity.ProductWebHookReq) error {
	// 动到我们的保险产品 触发更新产品
	go func(req productEntity.ProductWebHookReq) {
		// 给每个协程独立超时控制，比如10秒
		newCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// 保证防止panic导致进程挂掉
		defer func() {
			if r := recover(); r != nil {
				logger.Error(ctx, "Product 更新协程 panic:", fmt.Sprintf("%v", r))
			}
		}()

		uid, err := p.userRepo.GetUserIDByShop(newCtx, req.AppId, req.Shop)
		if err != nil {
			logger.Error(ctx, "Product 更新协程 获取UID失败:"+err.Error())
			return
		}

		if uid > 0 {
			productId := p.productRepo.ExistsByProductID(newCtx, uid, req.ProductId)
			if productId > 0 {
				_, err := p.asynqRepo.ProductWebhookUpdateTask(ctx, uid, productId)
				if err != nil {
					logger.Error(ctx, "ProductUpdate 推送产品队列失败:", err.Error())
					return
				}
			}
		}

		if uid > 0 {
			product, _ := p.productRepo.First(ctx, uid)
			if product == nil {
				return
			}
		}

	}(req) // **参数传进来，避免闭包引用外层变量**
	return nil
}

func (p *ProductService) ProductDel(ctx context.Context, req productEntity.ProductWebHookReq) error {
	// **立刻返回**
	go func(req productEntity.ProductWebHookReq) {
		// 给每个协程独立超时控制，比如10秒
		newCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// 保证防止panic导致进程挂掉
		defer func() {
			if r := recover(); r != nil {
				logger.Error(ctx, "Product 删除协程 panic:", fmt.Sprintf("%v", r))
			}
		}()

		//查出 UID
		uid, err := p.userRepo.GetUserIDByShop(newCtx, req.AppId, req.Shop)
		if err != nil {
			logger.Error(ctx, "Product 删除协程 获取UID失败:", err.Error())
			return
		}

		if uid > 0 {
			productId := p.productRepo.ExistsByProductID(newCtx, uid, req.ProductId)
			if productId > 0 {
				_, err := p.asynqRepo.DelProductTask(ctx, uid, req.ProductId, 0)
				if err != nil {
					logger.Error(ctx, "删除通知 del_product_queue 推送队列失败:", err.Error())
					return
				}
			}
		}
	}(req) // **参数传进来，避免闭包引用外层变量**

	return nil
}
