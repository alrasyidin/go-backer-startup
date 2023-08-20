package service

import (
	"errors"

	"github.com/alrasyidin/bwa-backer-startup/db/models"
	"github.com/alrasyidin/bwa-backer-startup/handler/user/dto"
	"github.com/alrasyidin/bwa-backer-startup/repository"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	Register(input dto.RegisterUserRequest) (models.User, error)
	Login(input dto.RegisterUserRequest) (models.User, error)
}

type UserService struct {
	repo repository.IUserRepo
}

var (
	ErrUserNotFound    = errors.New("User not found")
	ErrInvalidPassword = errors.New("password is invalid")
)

// Constructor for UserService
func NewUserService(repo repository.IUserRepo) *UserService {
	return &UserService{
		repo,
	}
}

func (service *UserService) Register(input dto.RegisterUserRequest) (models.User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	user := models.User{
		Name:         input.Name,
		Occupation:   input.Occupation,
		Email:        input.Email,
		PasswordHash: string(passwordHash),
		Token:        "user",
	}

	newUser, err := service.repo.Save(user)
	if err != nil {
		return models.User{}, err
	}

	return newUser, nil
}

func (service *UserService) Login(input dto.LoginUserRequest) (models.User, error) {
	user, err := service.repo.FindByEmail(input.Email)
	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		return user, ErrInvalidPassword
	}

	return user, nil
}
