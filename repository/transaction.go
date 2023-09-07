package repository

import (
	"context"
	"errors"
	"time"

	"github.com/alrasyidin/bwa-backer-startup/db/models"
	customerror "github.com/alrasyidin/bwa-backer-startup/pkg/error"
	"gorm.io/gorm"
)

type ITransactionRepo interface {
	FindByCampaignID(ID int) ([]models.Transaction, error)
	FindByUserID(ID int) ([]models.Transaction, error)
	FindByID(ID int) (models.Transaction, error)
	Save(transaction models.Transaction) (models.Transaction, error)
	Update(transaction models.Transaction) (models.Transaction, error)
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := repo.DB.WithContext(ctx).Preload("User").Where("campaign_id = ?", ID).Order("id DESC").Find(&transactions).Error
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (repo *TransactionRepo) FindByUserID(ID int) ([]models.Transaction, error) {
	var transactions []models.Transaction

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := repo.DB.WithContext(ctx).Preload("Campaign.CampaignImages", &models.CampaignImage{
		IsPrimary: true,
	}).Where(&models.Transaction{
		UserId: ID,
	}).Find(&transactions).Error

	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (repo *TransactionRepo) FindByID(ID int) (models.Transaction, error) {
	var transaction models.Transaction
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := repo.DB.WithContext(ctx).Where("id = ?", ID).First(&transaction).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transaction, customerror.ErrNotFound
		}
		return transaction, err
	}

	return transaction, nil
}

func (repo *TransactionRepo) Save(transaction models.Transaction) (models.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := repo.DB.WithContext(ctx).Create(&transaction).Error
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (repo *TransactionRepo) Update(transaction models.Transaction) (models.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := repo.DB.WithContext(ctx).Save(&transaction).Error
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}
