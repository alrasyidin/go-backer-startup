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

func APIResponse(ctx *gin.Context, message, status string, code int, data any) {
	meta := Meta{
		Status:  status,
		Code:    code,
		Message: message,
	}

	response := Response{
		Meta: meta,
		Data: data,
	}

	ctx.JSON(code, response)
}
func ErrorResponse(ctx *gin.Context, message, status string, code int, data any, errorMessage any) {
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

	ctx.JSON(code, response)
}

func SuccessResponse(ctx *gin.Context, message string, data any) {
	APIResponse(ctx, message, "success", http.StatusOK, data)
}

func BadRequestResponse(ctx *gin.Context, message string, data any, errorMessage any) {
	ErrorResponse(ctx, message, "error", http.StatusBadRequest, data, errorMessage)
}

func FailedValidationResponse(ctx *gin.Context, message string, errors validator.ValidationErrors) {
	errorMessages := FormatValidationError(errors)
	ErrorResponse(ctx, message, "error", http.StatusUnprocessableEntity, nil, errorMessages)
}
