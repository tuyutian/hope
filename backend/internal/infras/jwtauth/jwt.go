package jwtauth

import (
	"context"
	"fmt"

	"backend/internal/domain/repo/jwtauth"
	"backend/pkg/crypto/bcrypt"
	"backend/pkg/jwt"
)

var _ jwtauth.JWTRepository = (*jwtRepoImpl)(nil)

type jwtRepoImpl struct {
	secret    string
	manager   *jwt.JwtManager
	aesCrypto bcrypt.BCrypto
}

// NewJWTRepository 创建 JWT 资源
func NewJWTRepository(secret string, manager *jwt.JwtManager, aes bcrypt.BCrypto) jwtauth.JWTRepository {
	return &jwtRepoImpl{
		secret:    secret,
		manager:   manager,
		aesCrypto: aes,
	}
}

// Verify 校验 token
func (j *jwtRepoImpl) Verify(ctx context.Context, token string) (*jwt.BizClaims, error) {
	return j.Parse(token)
}

func (j *jwtRepoImpl) GenerateToken(ctx context.Context, claims jwt.BizClaims) (string, string, error) {
	token, refreshToken, err := j.manager.GenerateToken(claims)
	if err != nil {
		return "", "", err
	}
	return token, refreshToken, nil
}

// Parse 将jwt的claims字段，解密
func (j *jwtRepoImpl) Parse(payload string) (*jwt.BizClaims, error) {
	// 解析 token
	claims, err := j.manager.VerifyToken(payload)
	if err != nil {
		return nil, err
	}
	return &claims.BizClaims, err
}

func (j *jwtRepoImpl) cryptoTokenParse(token string) (*jwt.BizClaims, error) {
	// 解密 token
	payload, err := j.aesCrypto.Decrypt(token)
	if err != nil {
		return nil, err
	}
	// 解析 token
	claims, err := j.manager.VerifyToken(payload)
	if err != nil {
		return nil, err
	}
	return &claims.BizClaims, err
}

// UpdateSecret 更新 JWT 的 secret 和 manager
func (j *jwtRepoImpl) UpdateSecret(secret string, manager *jwt.JwtManager) {
	fmt.Printf("Updating JWT secret from %s to %s\n", j.secret[:8]+"...", secret[:8]+"...")
	j.secret = secret
	j.manager = manager
}
