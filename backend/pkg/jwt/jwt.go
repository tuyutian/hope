package jwt

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	// ErrTokenExpired token is expired
	ErrTokenExpired = errors.New("token is expired")
	// ErrTokenMalformed token is malformed
	ErrTokenMalformed = errors.New("token is malformed")
	// ErrTokenNotValidYet token is not valid yet
	ErrTokenNotValidYet = errors.New("token is not valid yet")
	// ErrTokenSignatureInvalid = token signature is invalid
	ErrTokenSignatureInvalid = errors.New("token signature is invalid")
)

// RefreshTokenValidate determine whether the token is allowed to refresh
type RefreshTokenValidate func(claims *CustomClaims) bool

// RefreshTokenClaims when refreshing the token, verify the validity of the token
type RefreshTokenClaims func(accessToken string, refreshToken string) (*CustomClaims, error)

// JwtManager for jwt auth
type JwtManager struct {
	secretKey         []byte
	accessExpiration  time.Duration
	refreshExpiration time.Duration
	signingMethod     jwt.SigningMethod

	refreshTokenValidate RefreshTokenValidate
	refreshTokenClaims   RefreshTokenClaims
}

// BizClaims BizClaims business claims
type BizClaims struct {
	Dest    string `json:"dest"`               // <shop-name.myshopify.com>
	UserID  int64  `json:"user_id,omitempty"`  // <user ID>
	AdminID int64  `json:"admin_id,omitempty"` // <admin ID>
	Jti     string `json:"jti,omitempty"`      // <random UUID>
	Sid     string `json:"sid,omitempty"`      // <session ID>
	Sig     string `json:"sig,omitempty"`      // <signature>
}

// CustomClaims .
type CustomClaims struct {
	jwt.RegisteredClaims

	BizClaims
}

// New return *JwtManager
func New(secretKey string, opts ...Option) *JwtManager {
	m := &JwtManager{
		secretKey:         []byte(secretKey),
		accessExpiration:  24 * time.Hour,
		refreshExpiration: 25 * time.Hour,
		signingMethod:     jwt.SigningMethodHS256,
	}

	for _, opt := range opts {
		opt(m)
	}

	// set default claims validator
	if m.refreshTokenValidate == nil {
		m.refreshTokenValidate = m.defaultRefreshTokenValidate
	}

	if m.refreshTokenClaims == nil {
		m.refreshTokenClaims = m.defaultRefreshTokenClaims
	}

	return m
}

// GenerateToken return accessToken and refreshToken
func (m *JwtManager) GenerateToken(bizClaim BizClaims, opts ...CustomClaimsOption) (accessToken string, refreshToken string, err error) {
	claims := &CustomClaims{
		BizClaims: bizClaim,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.accessExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),         // Issuance time
			NotBefore: jwt.NewNumericDate(time.Now()),         // Effective time
			ID:        strconv.FormatInt(bizClaim.UserID, 10), // JWT id
		},
	}

	for _, opt := range opts {
		opt(claims)
	}

	// generate access token
	accessToken, err = m.generateToken(claims)
	if err != nil {
		return "", "", err
	}

	// generate refresh token
	refreshToken, err = m.generateToken(jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.refreshExpiration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),         // Issuance time
		NotBefore: jwt.NewNumericDate(time.Now()),         // Effective time
		ID:        strconv.FormatInt(bizClaim.UserID, 10), // JWT id
	})
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// generateToken for generate token with claims
func (m *JwtManager) generateToken(claims jwt.Claims) (string, error) {
	return jwt.NewWithClaims(m.signingMethod, claims).SignedString(m.secretKey)
}

// VerifyToken validate access token
func (m *JwtManager) VerifyToken(accessToken string) (*CustomClaims, error) {
	claims := &CustomClaims{}
	jwtToken, err := jwt.ParseWithClaims(accessToken, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(jwt.SigningMethod); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return m.secretKey, nil
	})

	if err != nil {
		return nil, parseValidationError(err)
	}

	if !jwtToken.Valid {
		return nil, errors.New("invalid accessToken")
	}

	return claims, nil
}

// RefreshToken for jwt refresh token
func (m *JwtManager) RefreshToken(accessToken string, refreshToken string, opts ...CustomClaimsOption) (newAccessToken string, newRefreshToken string, err error) {
	claims, err := m.refreshTokenClaims(accessToken, refreshToken)
	if err != nil {
		return "", "", err
	}

	if m.refreshTokenValidate(claims) {
		return m.GenerateToken(claims.BizClaims, opts...)
	}

	return accessToken, refreshToken, nil
}

// defaultRefreshTokenClaims default RefreshTokenClaims
func (m *JwtManager) defaultRefreshTokenClaims(accessToken string, refreshToken string) (*CustomClaims, error) {
	// verify refresh token
	if _, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return m.secretKey, nil
	}); err != nil {
		return nil, err
	}

	// verify access token
	return m.VerifyToken(accessToken)
}

// defaultRefreshTokenValidate default RefreshTokenValidate
func (m *JwtManager) defaultRefreshTokenValidate(claims *CustomClaims) bool {
	return time.Until(claims.ExpiresAt.Time) < m.accessExpiration>>1
}

// parseValidationError parse jwt error
func parseValidationError(err error) error {
	if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
		return ErrTokenSignatureInvalid
	}

	if errors.Is(err, jwt.ErrTokenExpired) {
		return ErrTokenExpired
	}

	if errors.Is(err, jwt.ErrTokenNotValidYet) {
		return ErrTokenNotValidYet
	}

	if errors.Is(err, jwt.ErrTokenMalformed) {
		return ErrTokenMalformed
	}

	// others error
	return fmt.Errorf("token validation failed,error:%v", err)
}
