package service

import (
	"fmt"

	"github.com/alrasyidin/bwa-backer-startup/db/models"
	"github.com/alrasyidin/bwa-backer-startup/handler/campaign/dto"
	"github.com/alrasyidin/bwa-backer-startup/pkg/common"
	customerror "github.com/alrasyidin/bwa-backer-startup/pkg/error"
	"github.com/alrasyidin/bwa-backer-startup/repository"
	"github.com/gosimple/slug"
)

type ICampaignService interface {
	GetCampaigns(input dto.GetCampaignsRequest) ([]models.Campaign, *common.Pagination, error)
	GetCampaign(input dto.GetCampaignDetailRequest) (models.Campaign, error)
	CreateCampaign(input dto.CreateCampaignRequest) (models.Campaign, error)
	UpdateCampaign(inputID dto.GetCampaignDetailRequest, inputData dto.CreateCampaignRequest) (models.Campaign, error)
	SaveCampaignImage(input dto.SaveCampaignImageRequest, fileLocation string) (models.CampaignImage, error)
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

func (service *CampaignService) GetCampaigns(input dto.GetCampaignsRequest) ([]models.Campaign, *common.Pagination, error) {
	if input.UserID != 0 {
		campaigns, totalItems, err := service.repo.FindByUserIDWithCount(input.UserID, input.Page, input.PerPage)

		if err != nil {
			return campaigns, nil, err
		}

		pagination := common.NewPagination(input.Page, input.PerPage, totalItems)
		return campaigns, pagination, nil
	}
	campaigns, totalItems, err := service.repo.FindAllWithCount(input.Page, input.PerPage)
	pagination := common.NewPagination(input.Page, input.PerPage, totalItems)

	if err != nil {
		return campaigns, nil, err
	}
	return campaigns, pagination, nil
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

	return campaign, nil
}

func (service *CampaignService) UpdateCampaign(inputID dto.GetCampaignDetailRequest, inputData dto.CreateCampaignRequest) (models.Campaign, error) {
	campaign, err := service.repo.FindByID(inputID.ID)
	if err != nil {
		return campaign, err
	}

	if campaign.UserId != inputData.User.ID {
		return campaign, customerror.ErrNotOwnedCampaign
	}

	campaign.Name = inputData.Name
	campaign.ShortDescription = inputData.ShortDescription
	campaign.Description = inputData.Description
	campaign.GoalAmount = inputData.GoalAmount
	campaign.Perks = inputData.Perks

	updatedCampaign, err := service.repo.Update(campaign)
	if err != nil {
		return updatedCampaign, err
	}

	return updatedCampaign, nil
}

func (service *CampaignService) SaveCampaignImage(input dto.SaveCampaignImageRequest, fileLocation string) (models.CampaignImage, error) {
	campaign, err := service.repo.FindByID(input.CampaignID)
	if err != nil {
		return models.CampaignImage{}, err
	}

	if campaign.UserId != input.User.ID {
		return models.CampaignImage{}, customerror.ErrNotOwnedCampaign
	}

	if input.IsPrimary {
		_, err := service.repo.MarkAllImageAsNonPrimary(input.CampaignID)
		if err != nil {
			return models.CampaignImage{}, err
		}
	}
	campaignImage := models.CampaignImage{
		CampaignID: input.CampaignID,
		FileName:   fileLocation,
		IsPrimary:  input.IsPrimary,
	}
	newCampaignImage, err := service.repo.SaveImage(campaignImage)
	if err != nil {
		return newCampaignImage, err
	}

	return newCampaignImage, nil
}
