package dto

import "github.com/alrasyidin/bwa-backer-startup/db/models"

type CreateTransactionRequest struct {
	Amount     int `json:"amount" bind:"required"`
	CampaignID int `json:"campaign_id" bind:"required"`
	User       models.User
}

type TransactionResponse struct {
	ID         int    `json:"id"`
	CampaignID int    `json:"campaign_id"`
	UserId     int    `json:"user_id"`
	Amount     int    `json:"amount"`
	Code       string `json:"code"`
	PaymentURL string `json:"payment_url"`
	Status     string `json:"status"`
}

func FormatTransactionResponse(transaction models.Transaction) *TransactionResponse {
	return &TransactionResponse{
		ID:         transaction.ID,
		CampaignID: transaction.CampaignID,
		UserId:     transaction.UserId,
		Amount:     transaction.Amount,
		Code:       transaction.Code,
		PaymentURL: transaction.PaymentURL,
		Status:     transaction.Status,
	}
}
