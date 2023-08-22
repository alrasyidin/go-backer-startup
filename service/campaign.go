package service

import (
	"fmt"

	"github.com/alrasyidin/bwa-backer-startup/db/models"
	"github.com/alrasyidin/bwa-backer-startup/handler/campaign/dto"
	"github.com/alrasyidin/bwa-backer-startup/repository"
	"github.com/gosimple/slug"
)

type ICampaignService interface {
	GetCampaigns(UserID int) ([]models.Campaign, error)
	GetCampaign(input dto.GetCampaignDetailRequest) (models.Campaign, error)
	CreateCampaign(input dto.CreateCampaignRequest) (models.Campaign, error)
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

func (service *CampaignService) CreateCampaign(input dto.CreateCampaignRequest) (models.Campaign, error) {
	slugCandidate := fmt.Sprintf("%s-%d", input.Name, input.User.ID)

	campaign := models.Campaign{
		Name:             input.Name,
		UserId:           input.User.ID,
		ShortDescription: input.ShortDescription,
		Description:      input.Description,
		GoalAmount:       input.GoalAmount,
		Perks:            input.Perks,
		Slug:             slug.Make(slugCandidate),
	}

	campaign, err := service.repo.Save(campaign)
	if err != nil {
		return campaign, err
	}

	return campaign, err
}
