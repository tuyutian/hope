package users

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"

	"backend/internal/domain/entity/users"
	ur "backend/internal/domain/repo/users"
)

var _ ur.UserCacheRepository = (*userCacheImpl)(nil)

type userCacheImpl struct {
	redisClient redis.UniversalClient
	userRepo    ur.UserRepository
}

// NewUserCacheRepository NewUserCacheRepo 用户缓存资源
func NewUserCacheRepository(redisClient redis.UniversalClient, userRepo ur.UserRepository) ur.UserCacheRepository {
	return &userCacheImpl{redisClient: redisClient, userRepo: userRepo}
}

// randomSecond 设置随机值，防止缓存雪崩
func (u *userCacheImpl) randomSecond() time.Duration {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return time.Duration(r.Intn(30)) * time.Second
}

// Set 将 users.User 写入缓存
func (u *userCacheImpl) Set(ctx context.Context, id int64, user *users.User, ttl time.Duration) error {
	// 添加前缀
	key := fmt.Sprintf("user:%d", id)

	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	// 写入缓存
	return u.redisClient.Set(ctx, key, string(data), ttl+u.randomSecond()).Err()
}

// Get 从缓存中根据 id 获取 users.User
func (u *userCacheImpl) Get(ctx context.Context, id int64) (*users.User, error) {
	// 添加前缀
	key := fmt.Sprintf("user:%d", id)

	result, err := u.redisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	user := &users.User{}
	if err := json.Unmarshal([]byte(result), user); err != nil {
		return nil, err
	}
	return user, nil
}

// SetByShop 将 users.User 写入缓存
func (u *userCacheImpl) SetByShop(ctx context.Context, appId string, shop string, user *users.User, ttl time.Duration) error {
	// 添加前缀
	key := fmt.Sprintf("app:%s:user:%s", appId, shop)
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	// 写入缓存
	return u.redisClient.Set(ctx, key, string(data), ttl+u.randomSecond()).Err()
}

// GetByShop 从缓存中根据 id 获取 users.User
func (u *userCacheImpl) GetByShop(ctx context.Context, appId string, shop string) (*users.User, error) {
	// 添加前缀
	key := fmt.Sprintf("app:%s:user:%s", appId, shop)
	result, err := u.redisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	user := &users.User{}
	if err := json.Unmarshal([]byte(result), user); err != nil {
		return nil, err
	}
	return user, nil
}
