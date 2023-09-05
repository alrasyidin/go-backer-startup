package service

import (
	"github.com/alrasyidin/bwa-backer-startup/db/models"
	"github.com/alrasyidin/bwa-backer-startup/handler/transaction/dto"
	customerror "github.com/alrasyidin/bwa-backer-startup/pkg/error"
	"github.com/alrasyidin/bwa-backer-startup/repository"
)

type ITransactionService interface {
	GetTransactionsByCampaignID(input dto.GetTransactionCampaignRequest) ([]models.Transaction, error)
}

type TransactionService struct {
	repo         repository.ITransactionRepo
	campaignRepo repository.ICampaignRepo
}

func NewTransactionService(repo repository.ITransactionRepo, campaignRepo repository.ICampaignRepo) *TransactionService {
	return &TransactionService{
		repo,
		campaignRepo,
	}
}

func (service *TransactionService) GetTransactionsByCampaignID(input dto.GetTransactionCampaignRequest) ([]models.Transaction, error) {
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
