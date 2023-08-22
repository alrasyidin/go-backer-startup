package handler

import (
	"errors"
	"fmt"
	"time"

	"github.com/alrasyidin/bwa-backer-startup/db/models"
	"github.com/alrasyidin/bwa-backer-startup/middleware"
	"github.com/alrasyidin/bwa-backer-startup/pkg/helper"
	"github.com/alrasyidin/bwa-backer-startup/pkg/tokenization"
	"github.com/alrasyidin/bwa-backer-startup/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/alrasyidin/bwa-backer-startup/handler/user/dto"
)

type UserHandler struct {
	service        service.IUserService
	tokenGenerator tokenization.Generator
}

func NewUserHandler(userService service.IUserService, tokenGenerator tokenization.Generator) *UserHandler {
	return &UserHandler{userService, tokenGenerator}
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
		helper.BadRequestResponse(c, "Register account failed", nil, nil)
		return
	}
	user, err := handler.service.Register(input)

	if err != nil {
		helper.BadRequestResponse(c, "Register account failed", nil, err.Error())
		return
	}

	token, err := handler.tokenGenerator.GenerateToken(user.ID, time.Hour)
	if err != nil {
		helper.BadRequestResponse(c, "Register account failed", nil, err.Error())
		return
	}
	response := dto.FormatUser(user, token)

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
		helper.BadRequestResponse(c, "Login failed", nil, nil)
		return
	}
	user, err := handler.service.Login(input)

	if err != nil {
		helper.BadRequestResponse(c, "Login failed", nil, nil)
		return
	}
	token, err := handler.tokenGenerator.GenerateToken(user.ID, time.Hour)
	if err != nil {
		helper.BadRequestResponse(c, "Register account failed", nil, err.Error())
		return
	}
	response := dto.FormatUser(user, token)

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
		helper.BadRequestResponse(c, "Check email failed", nil, nil)
		return
	}

	isEmailAvailable, err := handler.service.IsEmailAvailable(input.Email)

	data := gin.H{
		"is_available": isEmailAvailable,
	}
	if err != nil {
		helper.BadRequestResponse(c, err.Error(), data, nil)
		return
	}
	helper.SuccessResponse(c, "Email available", data)

}

func (handler *UserHandler) UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		helper.BadRequestResponse(c, "Failed to upload image", gin.H{"is_uploaded": false}, nil)
		return
	}
	currentUser := c.MustGet(middleware.AuthorizationUserKey).(models.User)
	userID := currentUser.ID
	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)
	err = c.SaveUploadedFile(file, path)

	if err != nil {
		helper.BadRequestResponse(c, "Failed to upload image", gin.H{"is_uploaded": false}, nil)
		return
	}

	_, err = handler.service.UploadAvatar(userID, path)
	if err != nil {
		helper.BadRequestResponse(c, "Failed to upload image", gin.H{"is_uploaded": false}, nil)
		return
	}

	helper.SuccessResponse(c, "Avatar successfully uploaded", gin.H{"is_uploaded": true})
}
