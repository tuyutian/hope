package jwt

import (
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestNew(t *testing.T) {
	secretKey := "a015b2c021844fdda940415df9e44cf2"
	jwtManager := New(secretKey,
		WithAccessExpiration(360*time.Second),
		WithRefreshExpiration(720*time.Second),
		WithSigningMethod(jwt.SigningMethodHS256),
		WithRefreshTokenValidate(func(claims *CustomClaims) bool { return true }),
	)

	// generate token
	accessToken, refreshToken, err := jwtManager.GenerateToken(
		BizClaims{
			UserID: 100,
		},
		CustomClaimsWithRegisteredClaims(jwt.RegisteredClaims{
			Issuer:  "admin",
			Subject: "jwttest",
		}),
	)
	if err != nil {
		t.Fatal(err)
	}

	log.Print("accessToken: ", accessToken)
	log.Print("refreshToken: ", refreshToken)

	// verify token
	claims, err := jwtManager.VerifyToken(accessToken)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("claims: %s", func() string {
		data, _ := json.Marshal(claims)
		return string(data)
	}())

	// refresh token
	accessToken, refreshToken, err = jwtManager.RefreshToken(accessToken, refreshToken)
	if err != nil {
		t.Fatal("jwt refresh error:", err)
	}

	log.Print("newAccessToken: ", accessToken)
	log.Print("newRefreshToken: ", refreshToken)
}
