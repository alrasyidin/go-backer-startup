package handler

import (
	"errors"
	"fmt"

	"github.com/alrasyidin/bwa-backer-startup/pkg/common"
	customerror "github.com/alrasyidin/bwa-backer-startup/pkg/error"
	"github.com/alrasyidin/bwa-backer-startup/pkg/helper"
	"github.com/alrasyidin/bwa-backer-startup/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/alrasyidin/bwa-backer-startup/handler/campaign/dto"
)

type CampaignHandler struct {
	service service.ICampaignService
}

func NewCampaignHandler(service service.ICampaignService) *CampaignHandler {
	return &CampaignHandler{
		service,
	}
}

// GetCampaigns 	godoc
// @Summary      Get List of Campaigns
// @Description  Get List of Campaigns
// @Tags         Campaign
// @Accept       json
// @Produce      json
// @Success      200  {object}  helper.Response
// @Failure      400  {object}  helper.Response
// @Failure      404  {object}  helper.Response
// @Failure      500  {object}  helper.Response
// @Router       /campaigns [get]
func (handler *CampaignHandler) GetCampaigns(c *gin.Context) {
	var input dto.GetCampaignsRequest

	err := c.ShouldBind(&input)
	if err != nil {
		helper.BadRequestResponse(c, "failed to get campaigns data", nil, err.Error())
		return
	}

	input.PaginationRequest = common.NewPaginationRequet(input.Page, input.PerPage)

	campaigns, pagination, err := handler.service.GetCampaigns(input)

	if err != nil {
		helper.BadRequestResponse(c, "failed to get campaigns data", nil, err.Error())
		return
	}

	helper.SuccessResponseWithPagination(c, "List of Campaigns", dto.FormatListCampaignResponse(campaigns), pagination)
}

// GetCampaign 	godoc
// @Summary      Get Campaign
// @Description  Get Detail of Campaign
// @Tags         Campaign
// @Param        id   path      int  true  "Campaign ID"
// @Accept       json
// @Produce      json
// @Success      200  {object}  helper.Response
// @Failure      400  {object}  helper.Response
// @Failure      404  {object}  helper.Response
// @Failure      500  {object}  helper.Response
// @Router       /campaigns/{id} [get]
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

// CreateCampaign 	godoc
// @Summary      Create Campaign
// @Description  Create a Campaign
// @Tags         Campaign
// @Param        campaign   body     dto.CreateCampaignRequest  true  "Campaign Request Body"
// @Accept       json
// @Produce      json
// @Security		 ApiKeyAuth
// @Success      200  {object}  helper.Response
// @Failure      400  {object}  helper.Response
// @Failure      404  {object}  helper.Response
// @Failure      500  {object}  helper.Response
// @Router       /campaigns [post]
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

	input.User = helper.GetCurrentUser(c, customerror.ErrNotOwnedCampaign.Error())

	campaign, err := handler.service.CreateCampaign(input)
	if err != nil {
		helper.InternalServerResponse(c, "failed to create campaign", nil, err.Error())
		return
	}

	helper.SuccessResponse(c, "success to create campaign", dto.FormatCampaignResponse(campaign))
}

// UpdateCampaign 	godoc
// @Summary      		Update Campaign
// @Description  		Update a Campaign
// @Tags         		Campaign
// @Param        		id   path      int  true  "Campaign ID"
// @Param        		campaign   body     dto.CreateCampaignRequest  true  "Campaign Request Body"
// @Accept       		json
// @Produce      		json
// @Security		 		ApiKeyAuth
// @Success      		200  {object}  helper.Response
// @Failure      		400  {object}  helper.Response
// @Failure      		404  {object}  helper.Response
// @Failure      		500  {object}  helper.Response
// @Router       		/campaigns/{id} [put]
func (handler *CampaignHandler) UpdateCampaign(c *gin.Context) {
	var inputID dto.GetCampaignDetailRequest

	err := c.ShouldBindUri(&inputID)
	if err != nil {
		helper.BadRequestResponse(c, "failed to get campaign id", nil, err.Error())
		return
	}

	var inputData dto.CreateCampaignRequest

	err = c.ShouldBindJSON(&inputData)
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			helper.FailedValidationResponse(c, "failed to update campaign", ve)
			return
		}
		helper.BadRequestResponse(c, "failed to update campaign", nil, err)
		return
	}

	inputData.User = helper.GetCurrentUser(c, customerror.ErrNotOwnedCampaign.Error())

	campaign, err := handler.service.UpdateCampaign(inputID, inputData)
	if err != nil {
		helper.BadRequestResponse(c, "failed to update campaign", nil, err.Error())
		return
	}

	helper.SuccessResponse(c, "success to update campaign", dto.FormatCampaignResponse(campaign))
}

// UploadImage		 	godoc
// @Summary      		Upload Campaign Image
// @Description  		Upload a Campaign Image
// @Tags         		Campaign Image
// @Param        		campaign  formData   dto.SaveCampaignImageRequest  true  "Campaign Upload Image"
// @Param        		image     formData   file  true  "Image"
// @Accept       		mpfd
// @Produce      		json
// @Security		 		ApiKeyAuth
// @Success      		200  {object}  helper.Response
// @Failure      		400  {object}  helper.Response
// @Failure      		404  {object}  helper.Response
// @Failure      		500  {object}  helper.Response
// @Router       		/campaign-iamges [post]
func (handler *CampaignHandler) UploadImage(c *gin.Context) {
	var input dto.SaveCampaignImageRequest
	data := map[string]bool{
		"is_uploaded": false,
	}
	err := c.ShouldBind(&input)
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			helper.FailedValidationResponse(c, "failed to upload campaign image", ve)
			return
		}
		helper.BadRequestResponse(c, "failed to upload campaign image", data, err)
		return
	}

	fileImage, err := c.FormFile("image")
	if err != nil {
		helper.BadRequestResponse(c, "Failed to upload campaign image", data, nil)
		return
	}
	input.User = helper.GetCurrentUser(c, customerror.ErrNotOwnedCampaign.Error())

	path := fmt.Sprintf("images/%d-%s", input.User.ID, fileImage.Filename)
	err = c.SaveUploadedFile(fileImage, path)
	if err != nil {
		helper.BadRequestResponse(c, "Failed to upload campaign image", data, nil)
		return
	}

	_, err = handler.service.SaveCampaignImage(input, path)
	if err != nil {
		helper.BadRequestResponse(c, "Failed to upload campaign image", data, nil)
		return
	}

	data["is_uploaded"] = true
	helper.SuccessResponse(c, "Success to upload campaign image", data)

}
