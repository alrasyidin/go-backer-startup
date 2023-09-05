package helper

import (
	"github.com/alrasyidin/bwa-backer-startup/db/models"
	"github.com/gin-gonic/gin"

	"github.com/alrasyidin/bwa-backer-startup/pkg/constant"
)

func GetCurrentUser(c *gin.Context, message string) models.User {
	currentUser, ok := c.MustGet(constant.AuthorizationUserKey).(models.User)

	if message == "" {
		message = "user not valid"
	}
	if !ok {
		InternalServerResponse(c, message, nil, nil)
		return models.User{}
	}
	return currentUser
}
