package helper

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/alrasyidin/bwa-backer-startup/pkg/common"
)

type Meta struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Response model info
type Response struct {
	Meta       Meta               `json:"meta"`
	Pagination *common.Pagination `json:"pagination,omitempty"`
	Data       any                `json:"data,omitempty"`
	Error      any                `json:"error,omitempty"`
}

func apiResponse(ctx *gin.Context, message, status string, code int, data any, errorMessage any, pagination *common.Pagination) Response {
	meta := Meta{
		Status:  status,
		Code:    code,
		Message: message,
	}

	response := Response{
		Meta:       meta,
		Data:       data,
		Error:      errorMessage,
		Pagination: pagination,
	}
	return response
}

func SuccessResponse(ctx *gin.Context, message string, data any) {
	response := apiResponse(ctx, message, "success", http.StatusOK, data, nil, nil)
	ctx.JSON(response.Meta.Code, response)
}

func SuccessResponseWithPagination(ctx *gin.Context, message string, data any, pagination *common.Pagination) {
	response := apiResponse(ctx, message, "success", http.StatusOK, data, nil, pagination)
	ctx.JSON(response.Meta.Code, response)
}

func BadRequestResponse(ctx *gin.Context, message string, data any, errorMessage any) {
	response := apiResponse(ctx, message, "error", http.StatusBadRequest, data, errorMessage, nil)
	ctx.JSON(response.Meta.Code, response)
}

func FailedValidationResponse(ctx *gin.Context, message string, errors validator.ValidationErrors) {
	errorMessages := FormatValidationError(errors)
	response := apiResponse(ctx, message, "error", http.StatusUnprocessableEntity, nil, errorMessages, nil)
	ctx.JSON(response.Meta.Code, response)
}

func AbortResponse(ctx *gin.Context, message string, data any) {
	response := apiResponse(ctx, message, "error", http.StatusUnauthorized, data, nil, nil)
	ctx.AbortWithStatusJSON(response.Meta.Code, response)
}

func NotFoundResponse(ctx *gin.Context, message string, data any) {
	response := apiResponse(ctx, message, "error", http.StatusNotFound, data, nil, nil)
	ctx.JSON(response.Meta.Code, response)
}

func InternalServerResponse(ctx *gin.Context, message string, data any, errorMessage any) {
	if message == "" {
		message = "server could not process your requiest"
	}
	response := apiResponse(ctx, message, "error", http.StatusInternalServerError, data, errorMessage, nil)
	ctx.JSON(response.Meta.Code, response)
}
