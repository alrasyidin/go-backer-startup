package repository

import (
	"context"
	"time"

	"github.com/alrasyidin/bwa-backer-startup/db/models"
	"gorm.io/gorm"
)

type IUserRepo interface {
	Save(user models.User) (models.User, error)
	FindByEmail(email string) (models.User, error)
	FindByID(id int) (models.User, error)
	Update(user models.User) (models.User, error)
}

type UserRepo struct {
	DB *gorm.DB
}

// Constructor for UserRepo
func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (repo *UserRepo) Save(user models.User) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := repo.DB.WithContext(ctx).Create(&user).Error
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
func (repo *UserRepo) FindByEmail(email string) (models.User, error) {
	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := repo.DB.WithContext(ctx).Where("email = ?", email).Find(&user).Error
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (repo *UserRepo) FindByID(id int) (models.User, error) {
	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := repo.DB.WithContext(ctx).Where("id = ?", id).Find(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

func (repo *UserRepo) Update(user models.User) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := repo.DB.WithContext(ctx).Save(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}
