package handler

import (
	"errors"

	customerror "github.com/alrasyidin/bwa-backer-startup/pkg/error"
	"github.com/alrasyidin/bwa-backer-startup/pkg/helper"
	"github.com/alrasyidin/bwa-backer-startup/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/alrasyidin/bwa-backer-startup/handler/transaction/dto"
)

type TransactionHandler struct {
	service        service.ITransactionService
	paymentService service.IPaymentService
}

func NewTransactionHandler(service service.ITransactionService, paymentService service.IPaymentService) *TransactionHandler {
	return &TransactionHandler{service, paymentService}
}

func (handler *TransactionHandler) GetCampaignTransactions(c *gin.Context) {
	var input dto.GetTransactionCampaignRequest

	err := c.ShouldBindUri(&input)
	if err != nil {
		helper.BadRequestResponse(c, "failed to get transaction campaign's", nil, err.Error())
		return
	}

	input.User = helper.GetCurrentUser(c, customerror.ErrNotOwnedCampaign.Error())

	transactions, err := handler.service.GetCampaignTransactions(input)
	if err != nil {
		helper.BadRequestResponse(c, "failed to get transaction campaign's", nil, err.Error())
		return
	}

	helper.SuccessResponse(c, "Success get transaction campaign's", dto.FormatListTransactionCampaignResponse(transactions))
}

func (handler *TransactionHandler) GetUserTransactions(c *gin.Context) {
	currentUser := helper.GetCurrentUser(c, customerror.ErrNotOwnedCampaign.Error())

	transactions, err := handler.service.GetUserTransactions(currentUser.ID)
	if err != nil {
		helper.BadRequestResponse(c, "failed to get transaction user's", nil, err.Error())
		return
	}

	helper.SuccessResponse(c, "Success get transaction user's", dto.FormatListTransactionUserResponse(transactions))
}

func (handler *TransactionHandler) CreateTransaction(c *gin.Context) {
	var input dto.CreateTransactionRequest

	err := c.ShouldBindJSON(&input)
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			helper.FailedValidationResponse(c, "failed to create transaction", ve)
			return
		}
		helper.BadRequestResponse(c, "failed to create transaction", nil, err)
		return
	}
	input.User = helper.GetCurrentUser(c, customerror.ErrNotOwnedCampaign.Error())

	transaction, err := handler.service.CreateTransaction(input)
	if err != nil {
		helper.BadRequestResponse(c, "failed to create transaction", nil, err.Error())
		return
	}

	helper.SuccessResponse(c, "Success to create transaction", dto.FormatTransactionResponse(transaction))
}

func (handler *TransactionHandler) ProcessPaymentNotification(c *gin.Context) {
	var input dto.TransactionNotificationRequest

	err := c.ShouldBindJSON(&input)
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			helper.FailedValidationResponse(c, "failed to process payment notification", ve)
			return
		}
		helper.BadRequestResponse(c, "failed to process payment notification", nil, err)
		return
	}

	err = handler.paymentService.ProcessPayment(input)
	if err != nil {
		helper.BadRequestResponse(c, "failed to process payment notification", nil, err)
		return
	}
	helper.SuccessResponse(c, "Success to process payment notification", nil)
}
