package dto

import "github.com/alrasyidin/bwa-backer-startup/db/models"

type CreateCampaignRequest struct {
	Name             string      `json:"name" binding:"required"`
	ShortDescription string      `json:"short_description" binding:"required"`
	Description      string      `json:"description" binding:"required"`
	GoalAmount       int         `json:"goal_amount" binding:"required,numeric"`
	Perks            string      `json:"perks" binding:"required"`
	User             models.User `swaggerignore:"true"`
}
