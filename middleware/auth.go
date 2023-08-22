package middleware

import (
	"strings"

	"github.com/alrasyidin/bwa-backer-startup/pkg/helper"
	"github.com/alrasyidin/bwa-backer-startup/pkg/tokenization"
	"github.com/alrasyidin/bwa-backer-startup/service"
	"github.com/gin-gonic/gin"
)

var (
	HeaderAuthorizationType = "Bearer"
	AuthorizationUserKey    = "currentUser"
)

func AuthMiddlware(userService service.IUserService, token tokenization.Generator) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")

		if header == "" {
			helper.AbortResponse(c, "token not provided", nil)
			return
		}

		fields := strings.Fields(header)
		if len(fields) < 2 {
			helper.AbortResponse(c, "invalid format token", nil)
			return
		}

		tokenType := fields[0]
		if tokenType != HeaderAuthorizationType {
			helper.AbortResponse(c, "unsupported token", nil)
			return
		}
		tokenValue := fields[1]

		claims, err := token.ValidateToken(tokenValue)
		if err != nil {
			helper.AbortResponse(c, err.Error(), nil)
			return
		}

		user, err := userService.GetUserByID(claims.UserID)

		if err != nil {
			helper.AbortResponse(c, "user not found", nil)
			return
		}

		c.Set(AuthorizationUserKey, user)

		c.Next()
	}
}
