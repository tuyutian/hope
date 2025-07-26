package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Option JwtManager option
type Option func(m *JwtManager)

// WithAccessExpiration accessToken expiration time
func WithAccessExpiration(exp time.Duration) Option {
	return func(m *JwtManager) {
		m.accessExpiration = exp
	}
}

// WithRefreshExpiration refreshToken expiration time
func WithRefreshExpiration(exp time.Duration) Option {
	return func(m *JwtManager) {
		m.refreshExpiration = exp
	}
}

// WithRefreshTokenClaims token refresh mechanism
func WithRefreshTokenClaims(fn RefreshTokenClaims) Option {
	return func(m *JwtManager) {
		m.refreshTokenClaims = fn
	}
}

// WithRefreshTokenValidate set refresh token validator
func WithRefreshTokenValidate(r RefreshTokenValidate) Option {
	return func(m *JwtManager) {
		m.refreshTokenValidate = r
	}
}

// WithSigningMethod default: jwt.SigningMethodES256
func WithSigningMethod(signingMethod jwt.SigningMethod) Option {
	return func(m *JwtManager) {
		m.signingMethod = signingMethod
	}
}

// CustomClaimsOption CustomClaims option
type CustomClaimsOption func(c *CustomClaims)

// CustomClaimsWithBizClaims .
func CustomClaimsWithBizClaims(bizClaims BizClaims) CustomClaimsOption {
	return func(c *CustomClaims) {
		c.BizClaims = bizClaims
	}
}

// CustomClaimsWithRegisteredClaims .
func CustomClaimsWithRegisteredClaims(regClaims jwt.RegisteredClaims) CustomClaimsOption {
	return func(c *CustomClaims) {
		c.RegisteredClaims = regClaims
	}
}
