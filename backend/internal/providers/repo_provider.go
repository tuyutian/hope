package providers

import (
	"github.com/redis/go-redis/v9"
	"xorm.io/xorm"

	"backend/internal/domain/repo/jobs"
	jwtRepo "backend/internal/domain/repo/jwtauth"
	"backend/internal/domain/repo/orders"
	"backend/internal/domain/repo/users"
	"backend/internal/infras/config"
	"backend/internal/infras/jwtauth"
	"backend/internal/interfaces/persistence/job"
	"backend/internal/interfaces/persistence/order"
	"backend/internal/interfaces/persistence/user"
	"backend/pkg/crypto/bcrypt"
	"backend/pkg/jwt"
)

// Repositories 这个providers层可以根据实际情况看是否要添加
// 资源列表
type Repositories struct {
	UserRepo         users.UserRepository
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
		userRepo,
		orderRepo,
		jobOrderRepo,
		orderSummaryRepo,
		jwtRepository,
		aesCrypto,
	}

	return r
}
