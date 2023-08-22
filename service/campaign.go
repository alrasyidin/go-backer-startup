package service

import (
	"github.com/alrasyidin/bwa-backer-startup/db/models"
	"github.com/alrasyidin/bwa-backer-startup/handler/campaign/dto"
	"github.com/alrasyidin/bwa-backer-startup/repository"
)

type ICampaignService interface {
	GetCampaigns(UserID int) ([]models.Campaign, error)
	GetCampaign(input dto.GetCampaignDetailRequest) (models.Campaign, error)
}

type CampaignService struct {
	repo repository.ICampaignRepo
}

// Constructor for CampaignService
func NewCampaignService(repo repository.ICampaignRepo) *CampaignService {
	return &CampaignService{
		repo,
	}
}

func (service *CampaignService) GetCampaigns(UserID int) ([]models.Campaign, error) {
	if UserID != 0 {
		campaigns, err := service.repo.FindByUserID(UserID)

		if err != nil {
			return campaigns, err
		}

		return campaigns, nil
	}
	campaigns, err := service.repo.FindAll()
	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}

func (service *CampaignService) GetCampaign(input dto.GetCampaignDetailRequest) (models.Campaign, error) {
	campaign, err := service.repo.FindByID(input.ID)
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}
