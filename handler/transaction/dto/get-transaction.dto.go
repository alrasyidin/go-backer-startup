package dto

import (
	"time"

	"github.com/alrasyidin/bwa-backer-startup/db/models"
)

type GetTransactionCampaignRequest struct {
	ID   int `uri:"id" binding:"required"`
	User models.User
}

type GetTransactionCampaignResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

func FormatTransactionCampaignResponse(transaction models.Transaction) GetTransactionCampaignResponse {
	return GetTransactionCampaignResponse{
		ID:        transaction.ID,
		Name:      transaction.User.Name,
		Amount:    transaction.Amount,
		CreatedAt: transaction.CreatedAt,
	}
}

func FormatListTransactionCampaignResponse(transactions []models.Transaction) []GetTransactionCampaignResponse {
	listTransactionResponse := []GetTransactionCampaignResponse{}

	if len(transactions) == 0 {
		return listTransactionResponse
	}
	for _, transaction := range transactions {
		listTransactionResponse = append(listTransactionResponse, FormatTransactionCampaignResponse(transaction))
	}

	return listTransactionResponse
}

type GetTransactionUserResponse struct {
	ID        int                  `json:"id"`
	Status    string               `json:"status"`
	Amount    int                  `json:"amount"`
	CreatedAt time.Time            `json:"created_at"`
	Campaign  CampaignInfoResponse `json:"campaign"`
}

type CampaignInfoResponse struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

func FormatTransactionUserResponse(transaction models.Transaction) GetTransactionUserResponse {
	response := GetTransactionUserResponse{
		ID:        transaction.ID,
		Status:    transaction.Status,
		Amount:    transaction.Amount,
		CreatedAt: transaction.CreatedAt,
		Campaign: CampaignInfoResponse{
			Name: transaction.Campaign.Name,
		},
	}

	if len(transaction.Campaign.CampaignImages) > 0 {
		response.Campaign.ImageURL = transaction.Campaign.CampaignImages[0].FileName
	}

	return response
}

func FormatListTransactionUserResponse(transactions []models.Transaction) []GetTransactionUserResponse {
	listTransactionResponse := []GetTransactionUserResponse{}

	if len(transactions) == 0 {
		return listTransactionResponse
	}

	for _, transaction := range transactions {
		listTransactionResponse = append(listTransactionResponse, FormatTransactionUserResponse(transaction))
	}

	return listTransactionResponse
}
