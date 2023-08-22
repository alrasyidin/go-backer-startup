package repository

import (
	"errors"
	"fmt"

	"github.com/alrasyidin/bwa-backer-startup/db/models"
	customerror "github.com/alrasyidin/bwa-backer-startup/pkg/error"
	"gorm.io/gorm"
)

type ICampaignRepo interface {
	FindAll() ([]models.Campaign, error)
	FindByUserID(UserID int) ([]models.Campaign, error)
	FindByID(ID int) (models.Campaign, error)
	Save(campaign models.Campaign) (models.Campaign, error)
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

func (repo *CampaignRepo) FindByID(ID int) (models.Campaign, error) {
	var campaign models.Campaign

	err := repo.DB.Where("id = ?", ID).Preload("User").Preload("CampaignImages").First(&campaign).Error
	fmt.Println("err", err)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return campaign, customerror.ErrNotFound
		}

		return campaign, err
	}

	return campaign, nil
}

func (repo *CampaignRepo) Save(campaign models.Campaign) (models.Campaign, error) {
	err := repo.DB.Create(campaign).Error

	if err != nil {
		return campaign, err
	}

	return campaign, err
}
