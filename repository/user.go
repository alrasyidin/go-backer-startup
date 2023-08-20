package repository

import (
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
	err := repo.DB.Create(&user).Error
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
func (repo *UserRepo) FindByEmail(email string) (models.User, error) {
	var user models.User
	err := repo.DB.Where("email = ?", email).Find(&user).Error
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (repo *UserRepo) FindByID(id int) (models.User, error) {
	var user models.User
	err := repo.DB.Where("id = ?", id).Find(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

func (repo *UserRepo) Update(user models.User) (models.User, error) {
	err := repo.DB.Save(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}
