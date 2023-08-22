package helper

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Meta struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Response struct {
	Meta  Meta `json:"meta"`
	Data  any  `json:"data,omitempty"`
	Error any  `json:"error,omitempty"`
}

func apiResponse(ctx *gin.Context, message, status string, code int, data any, errorMessage any) Response {
	meta := Meta{
		Status:  status,
		Code:    code,
		Message: message,
	}

	response := Response{
		Meta:  meta,
		Data:  data,
		Error: errorMessage,
	}
	return response
}

func SuccessResponse(ctx *gin.Context, message string, data any) {
	response := apiResponse(ctx, message, "success", http.StatusOK, data, nil)
	ctx.JSON(response.Meta.Code, response)
}

func BadRequestResponse(ctx *gin.Context, message string, data any, errorMessage any) {
	response := apiResponse(ctx, message, "error", http.StatusBadRequest, data, errorMessage)
	ctx.JSON(response.Meta.Code, response)
}

func FailedValidationResponse(ctx *gin.Context, message string, errors validator.ValidationErrors) {
	errorMessages := FormatValidationError(errors)
	response := apiResponse(ctx, message, "error", http.StatusUnprocessableEntity, nil, errorMessages)
	ctx.JSON(response.Meta.Code, response)
}

func AbortResponse(ctx *gin.Context, message string, data any) {
	response := apiResponse(ctx, message, "error", http.StatusUnauthorized, data, nil)
	ctx.AbortWithStatusJSON(response.Meta.Code, response)
}

func NotFoundResponse(ctx *gin.Context, message string, data any) {
	response := apiResponse(ctx, message, "error", http.StatusNotFound, data, nil)
	ctx.JSON(response.Meta.Code, response)
}
