package handler

import (
	"errors"
	"fmt"
	"time"

	"github.com/alrasyidin/bwa-backer-startup/config"
	"github.com/alrasyidin/bwa-backer-startup/db/models"
	"github.com/alrasyidin/bwa-backer-startup/pkg/constant"
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
	config         *config.Config
}

func NewUserHandler(userService service.IUserService, tokenGenerator tokenization.Generator, config *config.Config) *UserHandler {
	return &UserHandler{userService, tokenGenerator, config}
}

// Register				 	godoc
// @Summary      		Register
// @Description  		Register a new user
// @Tags         		Authentication
// @Param        		user   body     dto.RegisterUserRequest  true  "Register Request Body"
// @Accept       		json
// @Produce      		json
// @Success      		200  {object}  helper.Response
// @Failure      		400  {object}  helper.Response
// @Failure      		404  {object}  helper.Response
// @Failure      		500  {object}  helper.Response
// @Router       		/users/register [post]
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

// Login				 		godoc
// @Summary      		Login
// @Description  		Login a user
// @Tags         		Authentication
// @Param        		user   body     dto.LoginUserRequest  true  "Login Request Body"
// @Accept       		json
// @Produce      		json
// @Success      		200  {object}  helper.Response
// @Failure      		400  {object}  helper.Response
// @Failure      		404  {object}  helper.Response
// @Failure      		500  {object}  helper.Response
// @Router       		/users/session [post]
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
	token, err := handler.tokenGenerator.GenerateToken(user.ID, time.Duration(handler.config.AccessTokenDuration))
	if err != nil {
		helper.BadRequestResponse(c, "Register account failed", nil, err.Error())
		return
	}
	response := dto.FormatUser(user, token)

	helper.SuccessResponse(c, "Login success", response)
}

// CheckEmail				godoc
// @Summary      		Check email
// @Description  		Check email a new user
// @Tags         		Authentication
// @Param        		user   body     dto.EmailCheckRequest  true  "Check email Request Body"
// @Accept       		json
// @Produce      		json
// @Success      		200  {object}  helper.Response
// @Failure      		400  {object}  helper.Response
// @Failure      		404  {object}  helper.Response
// @Failure      		500  {object}  helper.Response
// @Router       		/users/email-check [post]
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

// UploadAvatar			godoc
// @Summary      		Upload avatar
// @Description  		Upload avatar
// @Tags         		Authentication
// @Param        		avatar   formData     file  true  "avatar file image"
// @Accept       		json
// @Produce      		json
// @Security      	ApiKeyAuth
// @Success      		200  {object}  helper.Response
// @Failure      		400  {object}  helper.Response
// @Failure      		404  {object}  helper.Response
// @Failure      		500  {object}  helper.Response
// @Router       		/users/avatar [post]
func (handler *UserHandler) UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		helper.BadRequestResponse(c, "Failed to upload image", gin.H{"is_uploaded": false}, nil)
		return
	}
	currentUser := c.MustGet(constant.AuthorizationUserKey).(models.User)
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

// CurrentUser			godoc
// @Summary      		Current user
// @Description  		Current user
// @Tags         		Authentication
// @Accept       		json
// @Produce      		json
// @Security      	ApiKeyAuth
// @Success      		200  {object}  helper.Response
// @Failure      		400  {object}  helper.Response
// @Failure      		404  {object}  helper.Response
// @Failure      		500  {object}  helper.Response
// @Router       		/users/me [get]
func (handler *UserHandler) Me(c *gin.Context) {
	currentUser := helper.GetCurrentUser(c, "")
	response := dto.FormatUser(currentUser, "")
	helper.SuccessResponse(c, "Avatar successfully uploaded", response)
}
