package models

import (
	"time"
)

type Transaction struct {
	ID         int
	CampaignId int
	UserId     int
	Amount     int
	Code       string
	Status     string
	User       User
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
