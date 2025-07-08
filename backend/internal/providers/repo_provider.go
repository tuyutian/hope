package providers

import (
	"github.com/redis/go-redis/v9"
	"xorm.io/xorm"

	"backend/internal/domain/repo"
	"backend/internal/domain/repo/apps"
	"backend/internal/domain/repo/billings"
	"backend/internal/domain/repo/carts"
	"backend/internal/domain/repo/jobs"
	jwtRepo "backend/internal/domain/repo/jwtauth"
	"backend/internal/domain/repo/orders"
	"backend/internal/domain/repo/products"
	"backend/internal/domain/repo/shopifys"
	"backend/internal/domain/repo/users"
	"backend/internal/infras/cache"
	userCacheRepo "backend/internal/infras/cache/users"
	"backend/internal/infras/config"
	"backend/internal/infras/jwtauth"
	"backend/internal/infras/shopify"
	shopifyBillingRepo "backend/internal/infras/shopify_graphql/billings"
	shopifyOrderRepo "backend/internal/infras/shopify_graphql/orders"
	shopifyProductRepo "backend/internal/infras/shopify_graphql/products"
	shopifyShopRepo "backend/internal/infras/shopify_graphql/shops"
	"backend/internal/interfaces/persistence/app"
	"backend/internal/interfaces/persistence/billing"
	"backend/internal/interfaces/persistence/cart"
	"backend/internal/interfaces/persistence/job"
	"backend/internal/interfaces/persistence/order"
	"backend/internal/interfaces/persistence/product"
	"backend/internal/interfaces/persistence/user"
	"backend/pkg/crypto/bcrypt"
	"backend/pkg/jwt"
)

// Repositories 这个providers层可以根据实际情况看是否要添加
// 资源列表
type Repositories struct {
	ShopifyRepos
	TableRepos
	CacheRepos
	ThirdPartRepos
	AsyncRepo jobs.AsynqRepository
}

type TableRepos struct {
	OrderSummaryRepo     orders.OrderSummaryRepository
	UserRepo             users.UserRepository
	AppAuthRepo          users.AppAuthRepository
	OrderInfoRep         orders.OrderInfoRepository
	ProductRepo          products.ProductRepository
	CartSettingRepo      carts.CartSettingRepository
	VariantRepo          products.VariantRepository
	OrderRepo            orders.OrderRepository
	JobOrderRepo         jobs.OrderRepository
	JobProductRepo       jobs.ProductRepository
	UserSubscriptionRepo billings.UserSubscriptionRepository
	AppRepo              apps.AppRepository
}

type CacheRepos struct {
	CacheRepo     repo.CacheRepository
	UserCacheRepo users.UserCacheRepository
}

type ThirdPartRepos struct {
	AesCrypto     bcrypt.BCrypto
	JwtRepo       jwtRepo.JWTRepository
	AliyunOssRepo repo.AliyunOSSRepository
}

type ShopifyRepos struct {
	ShopifyRepo             shopifys.ShopifyRepository
	ProductGraphqlRepo      shopifys.ProductGraphqlRepository
	ShopGraphqlRepo         shopifys.ShopGraphqlRepository
	OrderGraphqlRepo        shopifys.OrderGraphqlRepository
	SubscriptionGraphqlRepo shopifys.SubscriptionGraphqlRepository
	UsageChargeGraphqlRepo  shopifys.UsageChargeGraphqlRepository
}

// NewRepositories 创建 Repositories
func NewRepositories(db *xorm.Engine, redisClient redis.UniversalClient, appConf *config.AppConfig, opts ...Option) *Repositories {
	tableRepos := NewTableRepos(db, redisClient)
	cacheRepos := NewCacheRepos(redisClient, tableRepos.UserRepo)
	thirdPartRepos := NewThirdPartRepos(appConf)
	shopifyRepos := NewShopifyRepos(&appConf.Shopify)
	r := &Repositories{
		TableRepos:     tableRepos,
		CacheRepos:     cacheRepos,
		ThirdPartRepos: thirdPartRepos,
		ShopifyRepos:   shopifyRepos,
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func NewCacheRepos(redisClient redis.UniversalClient, userRepo users.UserRepository) CacheRepos {
	cacheRepo := cache.NewCacheRepository(redisClient)
	uCacheRepo := userCacheRepo.NewUserCacheRepository(redisClient, userRepo)
	return CacheRepos{
		CacheRepo:     cacheRepo,
		UserCacheRepo: uCacheRepo,
	}
}

func NewTableRepos(db *xorm.Engine, redisClient redis.UniversalClient) TableRepos {
	userRepo := user.NewUserRepository(db)
	orderRepo := order.NewOrderRepository(db)
	jobOrderRepo := job.NewOrderRepository(db)
	jobProductRepo := job.NewProductRepository(db)
	orderInfoRepo := order.NewOrderInfoRepository(db)
	productRepo := product.NewProductRepository(db)
	variantRepo := product.NewVariantRepository(db)
	cartSettingRepo := cart.NewCartSettingRepository(db)
	orderSummaryRepo := order.NewOrderSummaryRepository(db)
	appRepo := app.NewAppRepository(db, redisClient)
	appAuthRepo := user.NewAppAuthRepository(db)
	userSubscriptionRepo := billing.NewUserSubscriptionRepository(db)
	return TableRepos{
		UserRepo:             userRepo,
		OrderRepo:            orderRepo,
		JobOrderRepo:         jobOrderRepo,
		OrderSummaryRepo:     orderSummaryRepo,
		JobProductRepo:       jobProductRepo,
		OrderInfoRep:         orderInfoRepo,
		ProductRepo:          productRepo,
		VariantRepo:          variantRepo,
		AppAuthRepo:          appAuthRepo,
		CartSettingRepo:      cartSettingRepo,
		AppRepo:              appRepo,
		UserSubscriptionRepo: userSubscriptionRepo,
	}
}

func NewThirdPartRepos(appConf *config.AppConfig) ThirdPartRepos {
	aesCrypto := bcrypt.NewAesBCrypto(appConf.Crypto.AES.Key, appConf.Crypto.AES.IV)
	// JwtManager
	jwtManager := jwt.New(
		appConf.JWT.SecretKey,
		jwt.WithAccessExpiration(appConf.JWT.AccessExpiration),
		jwt.WithRefreshExpiration(appConf.JWT.RefreshExpiration),
	)
	jwtRepository := jwtauth.NewJWTRepository(appConf.JWT.SecretKey, jwtManager, aesCrypto)
	return ThirdPartRepos{
		JwtRepo:   jwtRepository,
		AesCrypto: aesCrypto,
	}
}

func NewShopifyRepos(shopifyConf *config.Shopify) ShopifyRepos {
	shopifyRepos := shopify.NewShopifyRepository(shopifyConf)
	shopGraphqlRepo := shopifyShopRepo.NewShopGraphqlRepository()
	productGraphqlRepo := shopifyProductRepo.NewProductGraphqlRepository()
	orderGraphqlRepo := shopifyOrderRepo.NewOrderGraphqlRepository()
	subscriptionGraphqlRepo := shopifyBillingRepo.NewSubscriptionGraphqlRepository()
	usageChargeGraphqlRepo := shopifyBillingRepo.NewUsageChargeGraphqlRepository()
	return ShopifyRepos{
		ShopifyRepo:             shopifyRepos,
		ProductGraphqlRepo:      productGraphqlRepo,
		ShopGraphqlRepo:         shopGraphqlRepo,
		OrderGraphqlRepo:        orderGraphqlRepo,
		SubscriptionGraphqlRepo: subscriptionGraphqlRepo,
		UsageChargeGraphqlRepo:  usageChargeGraphqlRepo,
	}
}
