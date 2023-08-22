package repository

import (
	"github.com/alrasyidin/bwa-backer-startup/db/models"
	"gorm.io/gorm"
)

type ICampaignRepo interface {
	FindAll() ([]models.Campaign, error)
	FindByUserID(UserID int) ([]models.Campaign, error)
}

type CampaignRepo struct {
	DB *gorm.DB
}

// Constructor for CampaignRepo
func NewCampaignRepo(db *gorm.DB) *CampaignRepo {
	return &CampaignRepo{DB: db}
}

func (repo *CampaignRepo) FindAll() ([]models.Campaign, error) {
	var campaigns []models.Campaign

	err := repo.DB.Preload("CampaignImages", "campaign_images.is_primary = true").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (repo *CampaignRepo) FindByUserID(UserID int) ([]models.Campaign, error) {
	var campaigns []models.Campaign

	err := repo.DB.Where("user_id = ?", UserID).Preload("CampaignImages", "campaign_images.is_primary = true").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}
