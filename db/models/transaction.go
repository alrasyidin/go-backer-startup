package models

import (
	"time"
)

type Transaction struct {
	ID         int `gorm:"primaryKey"`
	CampaignID int
	UserId     int
	Amount     int
	Code       string
	PaymentURL string
	Status     string
	User       *User
	Campaign   *Campaign
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
