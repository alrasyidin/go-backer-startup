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
	for _, transaction := range transactions {
		listTransactionResponse = append(listTransactionResponse, FormatTransactionCampaignResponse(transaction))
	}

	return listTransactionResponse
}
