package models

import "time"

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
}

type CampaignImage struct {
	ID         int `gorm:"primaryKey"`
	CampaignId int
	FileName   string
	IsPrimary  bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
