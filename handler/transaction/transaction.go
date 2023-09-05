package handler

import (
	customerror "github.com/alrasyidin/bwa-backer-startup/pkg/error"
	"github.com/alrasyidin/bwa-backer-startup/pkg/helper"
	"github.com/alrasyidin/bwa-backer-startup/service"
	"github.com/gin-gonic/gin"

	"github.com/alrasyidin/bwa-backer-startup/handler/transaction/dto"
)

type TransactionHandler struct {
	service service.ITransactionService
}

func NewTransactionHandler(service service.ITransactionService) *TransactionHandler {
	return &TransactionHandler{service}
}

func (handler *TransactionHandler) GetCampaignTransactions(c *gin.Context) {
	var input dto.GetTransactionCampaignRequest

	err := c.ShouldBindUri(&input)
	if err != nil {
		helper.BadRequestResponse(c, "failed to get transaction campaign's", nil, err.Error())
		return
	}

	input.User = helper.GetCurrentUser(c, customerror.ErrNotOwnedCampaign.Error())

	transactions, err := handler.service.GetTransactionsByCampaignID(input)
	if err != nil {
		helper.BadRequestResponse(c, "failed to get transaction campaign's", nil, err.Error())
		return
	}

	helper.SuccessResponse(c, "Success get transaction campaign's", dto.FormatListTransactionCampaignResponse(transactions))
}
