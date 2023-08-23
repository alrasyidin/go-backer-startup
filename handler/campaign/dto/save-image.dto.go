package dto

type CreateCampaignImageRequest struct {
	IsPrimary  bool `form:"is_primary"`
	CampaignID int  `form:"campaign_id"`
}
