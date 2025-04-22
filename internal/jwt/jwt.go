package myjwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/grigory222/avito-backend-trainee/config"
)

type Claims struct {
	// Встроенные стандартные claims
	jwt.RegisteredClaims

	// Кастомные claims
	//UserID string `json:"sub"` - встроенный
	Role string `json:"role"`
}

type Provider struct {
	secret []byte
}

func NewJwtProvider(cfg config.JWT) *Provider {
	return &Provider{[]byte(cfg.Secret)}
}

//func (provider *Provider) GenerateToken(user models.User, role string) (string, error) {
//
//}

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
