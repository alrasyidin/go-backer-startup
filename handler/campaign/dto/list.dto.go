package dto

import (
	"github.com/alrasyidin/bwa-backer-startup/db/models"
	"github.com/alrasyidin/bwa-backer-startup/pkg/common"
)

type GetCampaignsRequest struct {
	common.PaginationRequest
	UserID int `form:"user_id"`
}

type CampaignResponse struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	UserId           int    `json:"user_id"`
	ShortDescription string `json:"short_description"`
	Description      string `json:"description"`
	GoalAmount       int    `json:"goal_amount"`
	CurrentAmount    int    `json:"current_amount"`
	Slug             string `json:"slug"`
	ImageURL         string `json:"image_url"`
}

func FormatCampaignResponse(campaign models.Campaign) CampaignResponse {
	campaignResponse := CampaignResponse{
		ID:               campaign.ID,
		Name:             campaign.Name,
		UserId:           campaign.UserId,
		ShortDescription: campaign.ShortDescription,
		Description:      campaign.Description,
		GoalAmount:       campaign.GoalAmount,
		CurrentAmount:    campaign.CurrentAmount,
		Slug:             campaign.Slug,
		ImageURL:         "",
	}

	if len(campaign.CampaignImages) > 0 {
		campaignResponse.ImageURL = campaign.CampaignImages[0].FileName
	}

	return campaignResponse
}

func FormatListCampaignResponse(campaigns []models.Campaign) []CampaignResponse {
	listCampaignResponse := []CampaignResponse{}

	for _, campaign := range campaigns {
		listCampaignResponse = append(listCampaignResponse, FormatCampaignResponse(campaign))
	}

	return listCampaignResponse
}
