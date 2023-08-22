package handler

import (
	"fmt"
	"strconv"

	"github.com/alrasyidin/bwa-backer-startup/handler/campaign/dto"
	"github.com/alrasyidin/bwa-backer-startup/pkg/helper"
	"github.com/alrasyidin/bwa-backer-startup/service"
	"github.com/gin-gonic/gin"
)

type CampaignHandler struct {
	service service.ICampaignService
}

// Constructor for CampaignHandler
func NewCampaignHandler(service service.ICampaignService) *CampaignHandler {
	return &CampaignHandler{
		service,
	}
}

func (handler *CampaignHandler) GetCampaigns(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))
	campaigns, err := handler.service.GetCampaigns(userID)

	if err != nil {
		helper.BadRequestResponse(c, "Error get campaigns data", nil, err.Error())
		return
	}

	helper.SuccessResponse(c, "List of Campaigns", dto.FormatListCampaignResponse(campaigns))
}

func (handler *CampaignHandler) GetCampaign(c *gin.Context) {
	var input dto.GetCampaignDetailRequest

	err := c.ShouldBindUri(&input)
	if err != nil {
		helper.BadRequestResponse(c, "failed to get campaign detail", nil, err.Error())
		return
	}

	fmt.Println("input", input)

	campaign, err := handler.service.GetCampaign(input)

	if err != nil {
		helper.NotFoundResponse(c, err.Error(), nil)
		return
	}

	helper.SuccessResponse(c, "Campaign detail", campaign)
}
