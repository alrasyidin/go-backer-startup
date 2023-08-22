package dto

import "github.com/alrasyidin/bwa-backer-startup/db/models"

type CampaginRequest struct {
}

type CampaginResponse struct {
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

func FormatCampaignResponse(campaign models.Campaign) CampaginResponse {
	campaginResponse := CampaginResponse{
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
		campaginResponse.ImageURL = campaign.CampaignImages[0].FileName
	}

	return campaginResponse
}

func FormatListCampaignResponse(campaigns []models.Campaign) []CampaginResponse {
	listCampaginResponse := []CampaginResponse{}

	for _, campaign := range campaigns {
		listCampaginResponse = append(listCampaginResponse, FormatCampaignResponse(campaign))
	}

	return listCampaginResponse
}
