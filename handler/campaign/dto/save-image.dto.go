package dto

import "github.com/alrasyidin/bwa-backer-startup/db/models"

type SaveCampaignImageRequest struct {
	IsPrimary  bool        `form:"is_primary" binding:"required"`
	CampaignID int         `form:"campaign_id" binding:"required,numeric"`
	User       models.User `swaggerignore:"true"`
}
