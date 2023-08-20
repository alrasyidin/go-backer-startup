package tokenization

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Generator interface {
	GenerateToken(userID int, expire time.Duration) (string, error)
}

type JWTGenerator struct {
	SecretKey string
}

func NewJWTGenerator(secretKey string) Generator {
	return &JWTGenerator{
		SecretKey: secretKey,
	}
}

func (g *JWTGenerator) GenerateToken(userID int, expire time.Duration) (string, error) {
	claims := NewCustomClaim(userID, expire)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtToken, err := token.SignedString([]byte(g.SecretKey))
	if err != nil {
		return jwtToken, err
	}

	return jwtToken, nil
}
