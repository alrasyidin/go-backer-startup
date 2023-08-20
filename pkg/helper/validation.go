package helper

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func MsgForTag(tag string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	}
	return ""
}

func FormatValidationError(errors validator.ValidationErrors) gin.H {
	errorMessages := make(gin.H, len(errors))
	for _, v := range errors {
		field := strings.ToLower(v.Field())
		tag := v.Tag()

		if MsgForTag(tag) != "" {
			errorMessages[field] = MsgForTag(tag)
		} else {
			errorMessages[field] = v.Error()
		}
	}

	errorMessages = gin.H{
		"errors": errorMessages,
	}

	return errorMessages
}
