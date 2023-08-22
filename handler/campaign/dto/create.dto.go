package dto

import "github.com/alrasyidin/bwa-backer-startup/db/models"

type CreateCampaignRequest struct {
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
	Description      string `json:"description"`
	GoalAmount       int    `json:"goal_amount"`
	Perks            string `json:"perks"`
	User             models.User
}
