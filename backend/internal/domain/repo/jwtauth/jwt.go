package jwtauth

import (
	"context"

	"backend/pkg/jwt"
)

// JWTRepository 认证资源
type JWTRepository interface {
	Verify(ctx context.Context, token string) (*jwt.BizClaims, error)
	GenerateToken(ctx context.Context, claims jwt.BizClaims) (string, string, error)
}
