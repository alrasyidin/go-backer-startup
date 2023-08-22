package handler

import (
	"errors"
	"strconv"

	"github.com/alrasyidin/bwa-backer-startup/db/models"
	"github.com/alrasyidin/bwa-backer-startup/handler/campaign/dto"
	"github.com/alrasyidin/bwa-backer-startup/middleware"
	"github.com/alrasyidin/bwa-backer-startup/pkg/helper"
	"github.com/alrasyidin/bwa-backer-startup/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

	campaign, err := handler.service.GetCampaign(input)

	if err != nil {
		helper.NotFoundResponse(c, err.Error(), nil)
		return
	}

	helper.SuccessResponse(c, "Campaign detail", dto.FormatCampaignDetailResponse(campaign))
}

func (handler *CampaignHandler) CreateCampaign(c *gin.Context) {
	var input dto.CreateCampaignRequest

	err := c.ShouldBindJSON(&input)
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			helper.FailedValidationResponse(c, "failed to create campaign", ve)
			return
		}
		helper.BadRequestResponse(c, "failed to create campaign", nil, err)
		return
	}

	currentUser, ok := c.MustGet(middleware.AuthorizationUserKey).(models.User)

	if !ok {
		helper.InternalServerResponse(c, "failed to create campaign, user not valid", nil, err.Error())
		return
	}

	input.User = currentUser

	campaign, err := handler.service.CreateCampaign(input)
	if err != nil {
		helper.InternalServerResponse(c, "failed to create campaign", nil, err.Error())
		return
	}

	helper.SuccessResponse(c, "success to create campaign", dto.FormatCampaignResponse(campaign))
}
