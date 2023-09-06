package repository

import (
	"log"

	"github.com/alrasyidin/bwa-backer-startup/db/models"
	"gorm.io/gorm"
)

type ITransactionRepo interface {
	FindByCampaignID(ID int) ([]models.Transaction, error)
	FindByUserID(ID int) ([]models.Transaction, error)
}

type TransactionRepo struct {
	DB *gorm.DB
}

// Constructor for TransactionRepo
func NewTransactionRepo(db *gorm.DB) *TransactionRepo {
	return &TransactionRepo{DB: db}
}

func (repo *TransactionRepo) FindByCampaignID(ID int) ([]models.Transaction, error) {
	var transactions []models.Transaction

	err := repo.DB.Preload("User").Where("campaign_id = ?", ID).Order("id DESC").Find(&transactions).Error
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (repo *TransactionRepo) FindByUserID(ID int) ([]models.Transaction, error) {
	var transactions []models.Transaction

	err := repo.DB.Preload("Campaign.CampaignImages", &models.CampaignImage{
		IsPrimary: true,
	}).Where(&models.Transaction{
		UserId: ID,
	}).Find(&transactions).Error

	log.Printf("transactions %+v\n", transactions)
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}
