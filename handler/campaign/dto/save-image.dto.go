package dto

type SaveCampaignImageRequest struct {
	IsPrimary  bool `form:"is_primary" binding:"required"`
	CampaignID int  `form:"campaign_id" binding:"required,numeric"`
}
