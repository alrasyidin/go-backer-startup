package handler

import (
	"errors"

	"github.com/alrasyidin/bwa-backer-startup/pkg/helper"
	"github.com/alrasyidin/bwa-backer-startup/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/alrasyidin/bwa-backer-startup/handler/user/dto"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService}
}

func (handler *UserHandler) Register(c *gin.Context) {
	var input dto.RegisterUserRequest

	err := c.ShouldBindJSON(&input)
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			helper.FailedValidationResponse(c, "Register account failed, please check your input", ve)
			return
		}
		helper.BadRequestResponse(c, "Register account failed", nil)
		return
	}
	user, err := handler.service.Register(input)

	if err != nil {
		helper.BadRequestResponse(c, "Register account failed", nil)
		return
	}

	response := dto.FormatUser(user, "tokentokentokentokentokentoken")

	helper.SuccessResponse(c, "Account successfully register", response)
}

func (handler *UserHandler) Login(c *gin.Context) {
	var input dto.LoginUserRequest

	err := c.ShouldBindJSON(&input)
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			helper.FailedValidationResponse(c, "Login failed", ve)
			return
		}
		helper.BadRequestResponse(c, "Login failed", nil)
		return
	}
	user, err := handler.service.Login(input)

	if err != nil {
		helper.BadRequestResponse(c, "Login failed", nil)
		return
	}

	response := dto.FormatUser(user, "tokentokentokentokentokentoken")

	helper.SuccessResponse(c, "Login success", response)
}

func (handler *UserHandler) CheckEmailAvailability(c *gin.Context) {
	var input dto.EmailCheckRequest

	err := c.ShouldBindJSON(&input)
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			helper.FailedValidationResponse(c, "Check email failed", ve)
			return
		}
		helper.BadRequestResponse(c, "Check email failed", nil)
		return
	}

	isEmailAvailable, err := handler.service.IsEmailAvailable(input)

	data := gin.H{
		"is_available": isEmailAvailable,
	}
	if err != nil {
		helper.BadRequestResponse(c, err.Error(), data)
		return
	}
	helper.SuccessResponse(c, "Email available", data)

}
