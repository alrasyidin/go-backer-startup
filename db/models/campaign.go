package models

import (
	"time"
)

type Campaign struct {
	ID               int `gorm:"primaryKey"`
	Name             string
	UserId           int
	ShortDescription string
	Description      string
	GoalAmount       int
	CurrentAmount    int
	BackerCount      int
	Perks            string
	Slug             string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	CampaignImages   []CampaignImage
	User             *User
}

type CampaignImage struct {
	ID         int `gorm:"primaryKey"`
	CampaignID int
	FileName   string
	IsPrimary  bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
