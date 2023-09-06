package dto

import "github.com/alrasyidin/bwa-backer-startup/db/models"

type CreateTransactionRequest struct {
	Amount     int `json:"amount" bind:"required"`
	CampaignID int `json:"campaign_id" bind:"required"`
	User       models.User
}
