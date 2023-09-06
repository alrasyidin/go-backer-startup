package service

import (
	"github.com/alrasyidin/bwa-backer-startup/db/models"
	"github.com/alrasyidin/bwa-backer-startup/handler/transaction/dto"
	customerror "github.com/alrasyidin/bwa-backer-startup/pkg/error"
	"github.com/alrasyidin/bwa-backer-startup/pkg/payment"
	"github.com/alrasyidin/bwa-backer-startup/repository"
)

type ITransactionService interface {
	GetCampaignTransactions(input dto.GetTransactionCampaignRequest) ([]models.Transaction, error)
	GetUserTransactions(userID int) ([]models.Transaction, error)
	CreateTransaction(input dto.CreateTransactionRequest) (models.Transaction, error)
}

type TransactionService struct {
	repo            repository.ITransactionRepo
	campaignRepo    repository.ICampaignRepo
	paymentMidtrans payment.IMidtrans
}

func NewTransactionService(repo repository.ITransactionRepo, campaignRepo repository.ICampaignRepo, paymentMidtrans payment.IMidtrans) *TransactionService {
	return &TransactionService{
		repo,
		campaignRepo,
		paymentMidtrans,
	}
}

func (service *TransactionService) GetCampaignTransactions(input dto.GetTransactionCampaignRequest) ([]models.Transaction, error) {
	campaign, err := service.campaignRepo.FindByID(input.ID)
	if err != nil {
		return []models.Transaction{}, err
	}

	if campaign.UserId != input.User.ID {
		return []models.Transaction{}, customerror.ErrNotOwnedCampaign
	}

	transactions, err := service.repo.FindByCampaignID(input.ID)

	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (service *TransactionService) GetUserTransactions(userID int) ([]models.Transaction, error) {
	transactions, err := service.repo.FindByUserID(userID)

	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (service *TransactionService) CreateTransaction(input dto.CreateTransactionRequest) (models.Transaction, error) {
	transaction := models.Transaction{
		CampaignID: input.CampaignID,
		UserId:     input.User.ID,
		Amount:     input.Amount,
		Code:       "",
		Status:     "PENDING",
	}

	newTransaction, err := service.repo.Save(transaction)
	if err != nil {
		return newTransaction, err
	}

	paymentURl, err := service.paymentMidtrans.GetPaymentURL(transaction, input.User)

	if err != nil {
		return newTransaction, err
	}
	newTransaction.PaymentURL = paymentURl
	newTransaction, err = service.repo.Update(newTransaction)
	if err != nil {
		return newTransaction, err
	}

	return newTransaction, nil
}
