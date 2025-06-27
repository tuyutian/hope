package providers

import (
	"github.com/redis/go-redis/v9"
	"xorm.io/xorm"

	"backend/internal/domain/repo/carts"
	"backend/internal/domain/repo/jobs"
	jwtRepo "backend/internal/domain/repo/jwtauth"
	"backend/internal/domain/repo/orders"
	"backend/internal/domain/repo/products"
	"backend/internal/domain/repo/users"
	"backend/internal/infras/config"
	"backend/internal/infras/jwtauth"
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
	UserRepo         users.UserRepository
	OrderInfoRep     orders.OrderInfoRepository
	ProductRepo      products.ProductRepository
	CartSettingRepo  carts.CartSettingRepository
	VariantRepo      products.VariantRepository
	OrderRepo        orders.OrderRepository
	JobOrderRepo     jobs.OrderRepository
	OrderSummaryRepo orders.OrderSummaryRepository
	JwtRepo          jwtRepo.JWTRepository
	AesCrypto        bcrypt.BCrypto
}

// NewRepositories 创建 Repositories
func NewRepositories(db *xorm.Engine, redisClient redis.UniversalClient, appConf *config.AppConfig) *Repositories {
	userRepo := user.NewUserRepository(db)
	orderRepo := order.NewOrderRepository(db)
	jobOrderRepo := job.NewOrderRepository(db)
	orderInfoRepo := order.NewOrderInfoRepository(db)
	productRepo := product.NewProductRepository(db)
	variantRepo := product.NewVariantRepository(db)
	cartSettingRepo := cart.NewCartSettingRepository(db)
	orderSummaryRepo := order.NewOrderSummaryRepository(db)

	aesCrypto := bcrypt.NewAesBCrypto(appConf.Crypto.AES.Key, appConf.Crypto.AES.IV)
	// JwtManager
	jwtManager := jwt.New(
		appConf.JWT.SecretKey,
		jwt.WithAccessExpiration(appConf.JWT.AccessExpiration),
		jwt.WithRefreshExpiration(appConf.JWT.RefreshExpiration),
	)
	jwtRepository := jwtauth.NewJWTRepository(appConf.JWT.SecretKey, jwtManager, aesCrypto)

	r := &Repositories{
		UserRepo:         userRepo,
		OrderRepo:        orderRepo,
		JobOrderRepo:     jobOrderRepo,
		OrderSummaryRepo: orderSummaryRepo,
		JwtRepo:          jwtRepository,
		AesCrypto:        aesCrypto,
		OrderInfoRep:     orderInfoRepo,
		ProductRepo:      productRepo,
		VariantRepo:      variantRepo,
		CartSettingRepo:  cartSettingRepo,
	}

	return r
}
