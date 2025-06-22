package providers

import (
	"github.com/redis/go-redis/v9"
	"xorm.io/xorm"
)

// Repositories 这个providers层可以根据实际情况看是否要添加
// 资源列表
type Repositories struct {
}

// NewRepositories 创建 Repositories
func NewRepositories(db *xorm.Engine, redisClient redis.UniversalClient) *Repositories {

	r := &Repositories{}

	return r
}
