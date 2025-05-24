package myjwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/grigory222/avito-backend-trainee/config"
	"time"
)

const TimeToLive = 24 * time.Hour

type Claims struct {
	jwt.RegisteredClaims

	Role string `json:"role"`
}

type Provider struct {
	secret []byte
}

func NewJwtProvider(cfg config.JWT) *Provider {
	return &Provider{[]byte(cfg.Secret)}
}

func (provider *Provider) GenerateTokenWithRole(role string) (string, error) {
	now := time.Now()
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(TimeToLive)),
		},
		Role: role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(provider.secret)
}

func (provider *Provider) VerifyToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		// Убедимся, что используется HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return provider.secret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
