package tokenization

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("Token is invalid")
	ErrExpiredToken = errors.New("Token is expired")
)

type Generator interface {
	GenerateToken(userID int, expire time.Duration) (string, error)
	ValidateToken(encodedToken string) (*CustomClaim, error)
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

func (g *JWTGenerator) ValidateToken(encodedToken string) (*CustomClaim, error) {
	token, err := jwt.ParseWithClaims(encodedToken, &CustomClaim{}, func(t *jwt.Token) (interface{}, error) {

		_, ok := t.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, ErrInvalidToken
		}

		return []byte(g.SecretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*CustomClaim)

	if !ok {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
