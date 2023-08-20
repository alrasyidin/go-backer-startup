package tokenization

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaim struct {
	jwt.RegisteredClaims
	UserID int `json:"user_id"`
}

// Constructor for Claim
func NewCustomClaim(userID int, expire time.Duration) *CustomClaim {
	return &CustomClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
		},
		UserID: userID,
	}
}
