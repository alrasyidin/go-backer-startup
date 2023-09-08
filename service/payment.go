package service

import (
	"strconv"

	"github.com/alrasyidin/bwa-backer-startup/db/models"
	"github.com/alrasyidin/bwa-backer-startup/handler/transaction/dto"
	"github.com/alrasyidin/bwa-backer-startup/pkg/constant"
	"github.com/alrasyidin/bwa-backer-startup/repository"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"gorm.io/gorm"
)

type IPaymentService interface {
	GetPaymentURL(transaction models.Transaction, user models.User) (string, error)
	ProcessPayment(input dto.TransactionNotificationRequest) error
	WithTrx(tx *gorm.DB) *PaymentService
}

type PaymentConfig struct {
	ServerKey string
	EnvType   midtrans.EnvironmentType
}

type PaymentService struct {
	Config          *PaymentConfig
	TransactionRepo repository.ITransactionRepo
	CampaignRepo    repository.ICampaignRepo
}

func NewPayment(Config *PaymentConfig, TransactionRepo repository.ITransactionRepo, CampaignRepo repository.ICampaignRepo) IPaymentService {
	return &PaymentService{Config, TransactionRepo, CampaignRepo}
}

func (service *PaymentService) WithTrx(trx *gorm.DB) *PaymentService {
	service.TransactionRepo = service.TransactionRepo.WithTrx(trx)
	return service
}

func (service *PaymentService) GetPaymentURL(transaction models.Transaction, user models.User) (string, error) {
	midtransClient := snap.Client{}
	midtransClient.New(service.Config.ServerKey, service.Config.EnvType)

	request := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(transaction.ID),
			GrossAmt: int64(transaction.Amount),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: user.Name,
			Email: user.Email,
		},
	}

	snapURL, err := midtransClient.CreateTransactionUrl(request)
	if err != nil {
		return "", err
	}

	return snapURL, nil
}

func (service *PaymentService) ProcessPayment(input dto.TransactionNotificationRequest) error {
	transactionID, err := strconv.Atoi(input.OrderId)
	if err != nil {
		return err
	}

	// transaction, err := service.TransactionRepo.FindByID(transactionID)
	transaction, err := service.TransactionRepo.FindByID(transactionID)
	if err != nil {
		return err
	}

	if input.PaymentType == "credit_card" && input.TransactionStatus == "capture" && input.FraudStatus == "accept" {
		transaction.Status = constant.PENDING
	} else if input.TransactionStatus == "settlement" {
		transaction.Status = constant.PAID
	} else if input.TransactionStatus == "cancel" {
		transaction.Status = constant.CANCELLED
	} else if input.TransactionStatus == "deny" || input.TransactionStatus == "failure" || input.TransactionStatus == "expire" {
		transaction.Status = constant.FAILED
	}

	updatedTransaction, err := service.TransactionRepo.Update(transaction)
	if err != nil {
		return err
	}

	campaign, err := service.CampaignRepo.FindByID(updatedTransaction.CampaignID)
	if err != nil {
		return err
	}

	if updatedTransaction.Status == constant.PAID {
		campaign.BackerCount = campaign.BackerCount + 1
		campaign.CurrentAmount = campaign.CurrentAmount + updatedTransaction.Amount

		_, err := service.CampaignRepo.Update(campaign)
		if err != nil {
			return err
		}
	}
	return nil
}
