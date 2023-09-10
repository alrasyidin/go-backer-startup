package handler

import (
	"errors"

	customerror "github.com/alrasyidin/bwa-backer-startup/pkg/error"
	"github.com/alrasyidin/bwa-backer-startup/pkg/helper"
	"github.com/alrasyidin/bwa-backer-startup/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	"github.com/alrasyidin/bwa-backer-startup/handler/transaction/dto"
)

type TransactionHandler struct {
	service        service.ITransactionService
	paymentService service.IPaymentService
}

func NewTransactionHandler(service service.ITransactionService, paymentService service.IPaymentService) *TransactionHandler {
	return &TransactionHandler{service, paymentService}
}

// GetCampaignTransaction			godoc
// @Summary      							Get campaign transaction
// @Description  							Get campaign transaction
// @Tags         							Transaction
// @Accept       							json
// @Produce      							json
// @Security      						ApiKeyAuth
// @Success      							200  {object}  helper.Response
// @Failure      							400  {object}  helper.Response
// @Failure      							404  {object}  helper.Response
// @Failure      							500  {object}  helper.Response
// @Router       							/campaigns/:id/transactions [get]
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

// GetUserTransaction		 godoc
// @Summary      				 Get user transaction
// @Description  				 Get user transaction
// @Tags         				 Transaction
// @Accept       				 json
// @Produce      				 json
// @Security      			 ApiKeyAuth
// @Success      				 200  {object}  helper.Response
// @Failure      				 400  {object}  helper.Response
// @Failure      				 404  {object}  helper.Response
// @Failure      				 500  {object}  helper.Response
// @Router       				 /transactions/ [get]
func (handler *TransactionHandler) GetUserTransactions(c *gin.Context) {
	currentUser := helper.GetCurrentUser(c, customerror.ErrNotOwnedCampaign.Error())

	transactions, err := handler.service.GetUserTransactions(currentUser.ID)
	if err != nil {
		helper.BadRequestResponse(c, "failed to get transaction user's", nil, err.Error())
		return
	}

	helper.SuccessResponse(c, "Success get transaction user's", dto.FormatListTransactionUserResponse(transactions))
}

// CreateTransaction		 godoc
// @Summary      				 Create transaction
// @Description  				 Create a transaction
// @Tags         				 Transaction
// @Accept       				 json
// @Produce      				 json
// @Security      			 ApiKeyAuth
// @Success      				 200  {object}  helper.Response
// @Failure      				 400  {object}  helper.Response
// @Failure      				 404  {object}  helper.Response
// @Failure      				 500  {object}  helper.Response
// @Router       				 /transactions/ [post]
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

// ProcessTransaction		 godoc
// @Summary      				 Process payment notification transaction
// @Description  				 Process payment notification transaction
// @Tags         				 Transaction
// @Accept       				 json
// @Produce      				 json
// @Success      				 200  {object}  helper.Response
// @Failure      				 400  {object}  helper.Response
// @Failure      				 404  {object}  helper.Response
// @Failure      				 500  {object}  helper.Response
// @Router       				 /transactions/notification [post]
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

	tx := c.MustGet("db_trx").(*gorm.DB)

	err = handler.paymentService.WithTrx(tx).ProcessPayment(input)
	if err != nil {
		helper.BadRequestResponse(c, "failed to process payment notification", nil, err)
		return
	}
	helper.SuccessResponse(c, "Success to process payment notification", nil)
}
