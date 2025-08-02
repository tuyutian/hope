package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hibiken/asynq"

	"backend/internal/domain/entity/jobs"
	productEntity "backend/internal/domain/entity/products"
	shopifyEntity "backend/internal/domain/entity/shopifys"
	"backend/internal/domain/repo/carts"
	jobRepo "backend/internal/domain/repo/jobs"
	"backend/internal/domain/repo/products"
	shopifyRepo "backend/internal/domain/repo/shopifys"
	"backend/internal/domain/repo/users"
	"backend/internal/infras/shopify_graphql"
	"backend/internal/providers"
	"backend/pkg/logger"
	"backend/pkg/utils"
)

type ProductService struct {
	productRepo        products.ProductRepository
	variantRepo        products.VariantRepository
	userRepo           users.UserRepository
	cartSettingRepo    carts.CartSettingRepository
	jobProductRepo     jobRepo.ProductRepository
	shopifyRepo        shopifyRepo.ShopifyRepository
	productGraphqlRepo shopifyRepo.ProductGraphqlRepository
}

func NewProductService(repos *providers.Repositories) *ProductService {
	return &ProductService{
		productRepo:        repos.ProductRepo,
		variantRepo:        repos.VariantRepo,
		userRepo:           repos.UserRepo,
		cartSettingRepo:    repos.CartSettingRepo,
		jobProductRepo:     repos.JobProductRepo,
		shopifyRepo:        repos.ShopifyRepo,
		productGraphqlRepo: repos.ProductGraphqlRepo,
	}
}

func (p *ProductService) DelProduct(ctx context.Context, t *asynq.Task) error {
	var payload jobs.DelProductPayload

	defer func() {
		if r := recover(); r != nil {
			logger.Error(ctx, "del_product_queue:", r)
		}
	}()

	logger.Info(ctx, "del_product_queue", "我正在消费删除配置")

	err := json.Unmarshal(t.Payload(), &payload)
	if err != nil {
		logger.Error(ctx, "del_product_queue:payload 反序列化失败", err)
		return nil
	}

	uid := payload.UserID
	productId := payload.ProductId
	//delType := payload.DelType
	// 未来删除还需要逻辑承载

	if err := p.cartSettingRepo.CloseCart(ctx, uid); err != nil {
		logger.Info(ctx, "Product  删除配置失败：", productId, err.Error())
	}
	if err := p.productRepo.DelShopifyProduct(ctx, uid); err != nil {
		logger.Info(ctx, "Product 清空主表失败：", productId, err.Error())
	}
	if err := p.variantRepo.DelShopifyVariant(ctx, uid); err != nil {
		logger.Info(ctx, "Product 清空变体失败：", productId, err.Error())
	}

	logger.Info(ctx, "Product 删除成功：", productId)
	return nil
}

func (p *ProductService) UploadProduct(ctx context.Context, t *asynq.Task) error {
	var payload jobs.ProductPayload

	defer func() {
		if r := recover(); r != nil {
			_ = p.fail(ctx, payload.JobId, "panic捕获", fmt.Errorf("%v", r))
		}
	}()

	logger.Info(ctx, "product_queue"+"我正在消费产品队列")

	err := json.Unmarshal(t.Payload(), &payload)
	if err != nil {
		return p.fail(ctx, payload.JobId, "payload 反序列化失败", err)
	}

	job, err := p.jobProductRepo.First(ctx, payload.JobId)
	if err != nil || job == nil {
		return p.fail(ctx, payload.JobId, "job日志不存在", err)
	}

	if job.IsSuccess == 1 {
		return p.skip(ctx, job.Id, "任务已完成，跳过")
	}

	uid := job.UserID

	err = p.jobProductRepo.UpdateJobTime(ctx, job.Id)

	if err != nil {
		return p.fail(ctx, job.Id, "更新Job时间失败", err)
	}

	user, err := p.userRepo.Get(ctx, uid)
	if err != nil || user == nil {
		return p.fail(ctx, job.Id, "查询用户信息失败", err)
	}
	product, err := p.productRepo.FirstProductByID(ctx, payload.UserProductId, uid)
	if err != nil || product == nil {
		return p.fail(ctx, job.Id, "查询产品信息失败", err)
	}

	variants, err := p.variantRepo.FindID(ctx, payload.UserProductId)
	if err != nil || variants == nil {
		return p.fail(ctx, job.Id, "查询变体失败", err)
	}
	productId := payload.ShopifyProductId
	shopName, _ := utils.GetShopName(user.Shop)
	client := shopify_graphql.NewGraphqlClient(shopName, user.AccessToken)
	p.productGraphqlRepo.WithClient(client)
	if productId == 0 {
		var variantId int64
		// 创建产品
		shopifyProductResp, err := p.productGraphqlRepo.CreateProductWithMedia(ctx, shopifyEntity.ProductCreateInput{
			Title:           product.Title,
			Status:          "ACTIVE",
			DescriptionHtml: product.Description,
			ProductType:     product.ProductType,
			Tags:            strings.Split(product.Tags, ","),
			Vendor:          product.Vendor,
			ProductOptions: []shopifyEntity.ProductOptionInput{
				{
					Name:     "Protectify",
					Position: 0,
					Values: []shopifyEntity.OptionValueCreateInput{
						{Name: "In-0"},
					},
				},
			},
		}, []shopifyEntity.CreateMediaInput{
			{
				Alt:              "image-1",
				OriginalSource:   product.ImageUrl,
				MediaContentType: "IMAGE",
			},
		}) // 通过 client 调用方法

		if err != nil {
			return p.fail(ctx, job.Id, "上传产品到 Shopify 失败", err)
		}
		shopifyProduct := shopifyProductResp.ProductCreate.Product
		productId = utils.GetIdFromShopifyGraphqlId(shopifyProduct.ID)
		variantId = utils.GetIdFromShopifyGraphqlId(shopifyProduct.Variants.Nodes[0].ID)
		// 删除默认变体
		err = p.productGraphqlRepo.DeleteVariant(ctx, productId, variantId)

		if err != nil {
			return p.fail(ctx, job.Id, "删除变体失败", err)
		}

		// 创建变体
		var variantCreateInput []*shopifyEntity.VariantCreateInput
		variantDbId := make(map[string]int64)

		for _, item := range variants {
			variantDbId[item.SkuName] = item.Id
			gqlVariant := &shopifyEntity.VariantCreateInput{
				Price: strconv.FormatFloat(item.Price, 'f', 2, 64),
				OptionValues: []shopifyEntity.VariantOptionValues{
					{Name: item.Sku1, OptionName: "Title"},
				},

				InventoryItem: shopifyEntity.InventoryItemInput{
					SKU:     item.SkuName,
					Tracked: false,
				},
			}
			variantCreateInput = append(variantCreateInput, gqlVariant)
		}

		gqlVariants, err := p.productGraphqlRepo.CreateVariants(ctx, productId, variantCreateInput) // 通过 client 调用方法

		if err != nil {
			return p.fail(ctx, job.Id, "创建变体失败", err)
		}

		// 填补数据库里的产品ID和变体ID
		if err := p.productRepo.UpdateProduct(ctx, payload.UserProductId, uid, &productEntity.UserProduct{
			ProductId:   productId,
			Status:      1,
			PublishTime: time.Now().Unix(),
		}); err != nil {
			return p.fail(ctx, job.Id, "更新产品信息失败", err)
		}

		for _, item := range gqlVariants {
			sku := item["sku"].(string)             // 确保 item["sku"] 是 string 类型
			dbVariantId, exists := variantDbId[sku] // 使用 sku 来获取存储的 Id
			if !exists {
				continue
			}

			_ = p.variantRepo.UpdateVariants(ctx, dbVariantId, uid, &productEntity.UserVariant{
				ProductId:   productId,
				VariantId:   utils.GetIdFromShopifyGraphqlId(item["id"].(string)),
				InventoryId: utils.GetIdFromShopifyGraphqlId(item["inventory_id"].(string)),
			})
		}

	} else {
		// 修改Shopify 产品及变体
		err = p.productGraphqlRepo.UpdateProduct(ctx, payload.ShopifyProductId, shopifyEntity.ProductUpdateInput{
			Status: "ACTIVE",
			Title:  product.Title,
			//DescriptionHtml: product.Description,
			ProductType: product.ProductType,
			//Tags:            strings.Split(product.Tags, ","),
			//Vendor: product.Vendor,
		}, []shopifyEntity.CreateMediaInput{
			//{
			//	Alt:              "image-1",
			//	MediaContentType: "IMAGE",
			//	OriginalSource:   product.ImageUrl,
			//},
		})
		if err != nil {
			return p.fail(ctx, job.Id, "修改Shopify产品失败", err)
		}

	}

	logger.Warn(ctx, "开始上传产品:", productId, "商店ID：", user.PublishId)

	if productId != 0 {
		err = p.productGraphqlRepo.PublishProduct(ctx, productId, user.PublishId)
		if err != nil {
			logger.Error(ctx, fmt.Sprintf("上传产品失败:%d 商店ID：%d error:%s", productId, user.PublishId, err.Error()))
		}
	}

	return p.ok(ctx, job.Id)
}

func (p *ProductService) HandleShopifyProduct(ctx context.Context, t *asynq.Task) error {
	var payload jobs.ShopifyProductPayload

	defer func() {
		if r := recover(); r != nil {
			logger.Error(ctx, "shopify_product_queue:", r)
		}
	}()

	logger.Info(ctx, "shopify_product_queue", "我正在消费同步shopify产品队列")

	err := json.Unmarshal(t.Payload(), &payload)
	if err != nil {
		logger.Error(ctx, "shopify_product_queue:payload 反序列化失败", err)
		return nil
	}

	uid := payload.UserID

	user, err := p.userRepo.Get(ctx, uid)
	if err != nil || user == nil {
		logger.Error(ctx, "shopify_product_queue:查询用户信息失败", err)
		return nil
	}
	shopName, _ := utils.GetShopName(user.Shop)
	client := shopify_graphql.NewGraphqlClient(shopName, user.AccessToken)
	p.productGraphqlRepo.WithClient(client)
	product, err := p.productRepo.FirstProductByID(ctx, payload.UserProductId, uid)
	if err != nil || product == nil {
		logger.Error(ctx, "shopify_product_queue:查询产品信息失败", err)
		return nil
	}

	variants, err := p.variantRepo.FindID(ctx, payload.UserProductId)
	if err != nil || variants == nil {
		logger.Error(ctx, "shopify_product_queue:查询变体失败", err)
		return nil
	}

	productId := product.ProductId

	// 继续同步Shopify 产品及变体
	err = p.productGraphqlRepo.UpdateProduct(ctx, productId, shopifyEntity.ProductUpdateInput{
		Title:           product.Title,
		Status:          "ACTIVE",
		DescriptionHtml: product.Description,
		ProductType:     product.ProductType,
		Tags:            strings.Split(product.Tags, ","),
		Vendor:          product.Vendor,
	}, []shopifyEntity.CreateMediaInput{})

	if err != nil {
		logger.Error(ctx, "shopify_product_queue:修改Shopify产品失败", err)
		return nil
	}

	// 创建变体
	var variantCreateInput []*shopifyEntity.VariantUpdateInput

	for _, item := range variants {
		if item.VariantId == 0 {
			continue
		}

		gqlVariant := &shopifyEntity.VariantUpdateInput{
			// 此处使用 Shopify 的全局唯一标识符，例如 "gid://shopify/ProductVariant/<id>"
			Id:    fmt.Sprintf("gid://shopify/Product/%d", item.VariantId),
			Price: strconv.FormatFloat(item.Price, 'f', 2, 64),
			OptionValues: []shopifyEntity.VariantOptionValues{
				{Name: item.Sku1, OptionName: "Title"},
			},

			InventoryItem: shopifyEntity.InventoryItemInput{
				SKU:     item.SkuName,
				Tracked: false,
			},
		}
		variantCreateInput = append(variantCreateInput, gqlVariant)
	}

	if len(variantCreateInput) == 0 {
		return nil
	}

	err = p.productGraphqlRepo.UpdateVariants(ctx, productId, variantCreateInput) // 通过 client 调用方法

	if err != nil {
		logger.Error(ctx, "shopify_product_queue:修改Shopify变体失败", err)
		return nil
	}

	return nil
}

// 这里要抽出来 失败和成功的逻辑 共用 解耦
func (p *ProductService) ok(ctx context.Context, jobID int64) error {
	_ = p.jobProductRepo.UpdateStatus(ctx, jobID, 1) // 3 表示失败
	return nil
}

func (p *ProductService) fail(ctx context.Context, jobID int64, msg string, err error) error {
	logger.Error(ctx, fmt.Sprintf("2product_queue:%d msg:%s error:%s", jobID, msg, err.Error()))
	_ = p.jobProductRepo.UpdateStatus(ctx, jobID, 3) // 3 表示失败
	return nil
}

func (p *ProductService) skip(ctx context.Context, jobID int64, msg string) error {
	logger.Info(ctx, fmt.Sprintf("2product_queue:%d msg:%s", jobID, msg))
	_ = p.jobProductRepo.UpdateStatus(ctx, jobID, 2) // 2 表示跳过
	return nil
}
