package users

import (
	"context"

	"backend/internal/domain/entity/users"
	userRepo "backend/internal/domain/repo/users"
	"backend/internal/providers"
)

type UserService struct {
	userRepo userRepo.UserRepository
}

func NewUserService(repos *providers.Repositories) *UserService {
	return &UserService{userRepo: repos.UserRepo}
}

func (s *UserService) GetLoginUserFromID(ctx context.Context, id int64) (interface{}, error) {
	return nil, nil
}

func (s *UserService) GetLoginUserFromShop(ctx context.Context, shop string) (*users.User, error) {
	user, err := s.userRepo.GetByShop(ctx, shop)
	return user, err
}

func (s *UserService) GetLoginAdminFromID(ctx context.Context, id int64) (interface{}, error) {
	return nil, nil
}
