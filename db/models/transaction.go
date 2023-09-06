package models

import (
	"time"
)

type Transaction struct {
	ID         int `gorm:"primaryKey"`
	CampaignId int
	UserId     int
	Amount     int
	Code       string
	Status     string
	User       *User
	Campaign   *Campaign
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
