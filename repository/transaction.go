package repository

import (
	"github.com/alrasyidin/bwa-backer-startup/db/models"
	"gorm.io/gorm"
)

type ITransactionRepo interface {
	FindByCampaignID(ID int) ([]models.Transaction, error)
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

	err := repo.DB.Model(&models.Transaction{}).Preload("User").Where("campaign_id = ?", ID).Order("id DESC").Find(&transactions).Error
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}
