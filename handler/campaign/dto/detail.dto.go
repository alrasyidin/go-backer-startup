package dto

type GetCampaignDetailRequest struct {
	ID int `uri:"id" binding:"required"`
}