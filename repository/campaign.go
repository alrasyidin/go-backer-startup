package repository

import (
	"context"
	"errors"
	"time"

	"github.com/alrasyidin/bwa-backer-startup/db/models"
	customerror "github.com/alrasyidin/bwa-backer-startup/pkg/error"
	"gorm.io/gorm"
)

type ICampaignRepo interface {
	FindAll() ([]models.Campaign, error)
	FindByUserID(UserID int) ([]models.Campaign, error)
	FindByID(ID int) (models.Campaign, error)
	Save(campaign models.Campaign) (models.Campaign, error)
	Update(campaign models.Campaign) (models.Campaign, error)
	SaveImage(campaignImage models.CampaignImage) (models.CampaignImage, error)
	MarkAllImageAsNonPrimary(campaignID int) (bool, error)
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := repo.DB.WithContext(ctx).Preload("CampaignImages", "campaign_images.is_primary = true").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (repo *CampaignRepo) FindByUserID(UserID int) ([]models.Campaign, error) {
	var campaigns []models.Campaign

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := repo.DB.WithContext(ctx).Where("user_id = ?", UserID).Preload("CampaignImages", "campaign_images.is_primary = true").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (repo *CampaignRepo) FindByID(ID int) (models.Campaign, error) {
	var campaign models.Campaign

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := repo.DB.WithContext(ctx).Where("id = ?", ID).Preload("User").Preload("CampaignImages").First(&campaign).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return campaign, customerror.ErrNotFound
		}

		return campaign, err
	}

	return campaign, nil
}

func (repo *CampaignRepo) Save(campaign models.Campaign) (models.Campaign, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := repo.DB.WithContext(ctx).Create(&campaign).Error

	if err != nil {
		return campaign, err
	}

	return campaign, err
}

func (repo *CampaignRepo) Update(campaign models.Campaign) (models.Campaign, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := repo.DB.WithContext(ctx).Save(&campaign).Error

	if err != nil {
		return campaign, err
	}

	return campaign, err
}

func (repo *CampaignRepo) SaveImage(campaignImage models.CampaignImage) (models.CampaignImage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := repo.DB.WithContext(ctx).Create(&campaignImage).Error

	if err != nil {
		return campaignImage, err
	}

	return campaignImage, err
}

func (repo *CampaignRepo) MarkAllImageAsNonPrimary(campaignID int) (bool, error) {
	// UPDATE campaign_images SET is_primary = true WHERE campaign_id = 1
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := repo.DB.WithContext(ctx).Model(&models.CampaignImage{}).Where("campaign_id = ?", campaignID).Update("is_primary", false).Error

	if err != nil {
		return false, err
	}

	return true, err
}
