package dto

type TransactionNotificationRequest struct {
	TransactionStatus string `json:"transaction_status"`
	FraudStatus       string `json:"fraud_status"`
	PaymentType       string `json:"payment_type"`
	OrderId           string `json:"order_id"`
}
