package users

import (
	"context"

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

func (s *UserService) GetLoginUserFromShop(ctx context.Context, shop string) (interface{}, error) {
	return nil, nil
}

func (s *UserService) GetLoginAdminFromID(ctx context.Context, id int64) (interface{}, error) {
	return nil, nil
}
