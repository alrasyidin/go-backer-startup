package dto

import (
	"strings"

	"github.com/alrasyidin/bwa-backer-startup/db/models"
)

type GetCampaignDetailRequest struct {
	ID int `uri:"id" binding:"required"`
}

type CampaignDetailResponse struct {
	ID               int                     `json:"id"`
	Name             string                  `json:"name"`
	UserId           int                     `json:"user_id"`
	ImageURL         string                  `json:"image_url"`
	ShortDescription string                  `json:"short_description"`
	Description      string                  `json:"description"`
	GoalAmount       int                     `json:"goal_amount"`
	CurrentAmount    int                     `json:"current_amount"`
	BackerCount      int                     `json:"backer_count"`
	Perks            []string                `json:"perks"`
	Slug             string                  `json:"slug"`
	User             CampaignUserResponse    `json:"user"`
	Images           []CampaignImageResponse `json:"images"`
}

type CampaignUserResponse struct {
	ImageURL string `json:"image_url"`
	Name     string `json:"name"`
}

type CampaignImageResponse struct {
	ImageURL  string `json:"image_url"`
	IsPrimary bool   `json:"is_primary"`
}

func FormatCampaignDetailResponse(campaign models.Campaign) CampaignDetailResponse {
	campaignDetailResponse := CampaignDetailResponse{
		ID:               campaign.ID,
		Name:             campaign.Name,
		UserId:           campaign.UserId,
		ShortDescription: campaign.ShortDescription,
		Description:      campaign.Description,
		GoalAmount:       campaign.GoalAmount,
		CurrentAmount:    campaign.CurrentAmount,
		Slug:             campaign.Slug,
		ImageURL:         "",
		Perks:            []string{},
		User:             CampaignUserResponse{},
		Images:           []CampaignImageResponse{},
	}

	if campaign.Perks != "" {
		perks := strings.Split(campaign.Perks, ",")
		for _, perk := range perks {
			campaignDetailResponse.Perks = append(campaignDetailResponse.Perks, strings.TrimSpace(perk))
		}
	}

	if len(campaign.CampaignImages) > 0 {
		campaignDetailResponse.ImageURL = campaign.CampaignImages[0].FileName
	}

	campaignUserResponse := CampaignUserResponse{
		ImageURL: campaign.User.AvatarFileName,
		Name:     campaign.User.Name,
	}
	campaignDetailResponse.User = campaignUserResponse

	if len(campaign.CampaignImages) > 0 {
		for _, campaignImage := range campaign.CampaignImages {
			image := CampaignImageResponse{
				ImageURL:  campaignImage.FileName,
				IsPrimary: campaignImage.IsPrimary,
			}
			campaignDetailResponse.Images = append(campaignDetailResponse.Images, image)
		}
	}

	return campaignDetailResponse
}
