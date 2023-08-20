package service

import (
	"errors"
	"fmt"

	"github.com/alrasyidin/bwa-backer-startup/db/models"
	"github.com/alrasyidin/bwa-backer-startup/handler/user/dto"
	"github.com/alrasyidin/bwa-backer-startup/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type IUserService interface {
	Register(input dto.RegisterUserRequest) (models.User, error)
	Login(input dto.LoginUserRequest) (models.User, error)
	IsEmailAvailable(email string) (bool, error)
	UploadAvatar(ID int, fileLocation string) (models.User, error)
}

type UserService struct {
	repo repository.IUserRepo
}

var (
	ErrUserNotFound      = errors.New("User not found")
	ErrInvalidPassword   = errors.New("password is invalid")
	ErrEmailAlreadyTaken = errors.New("email has been registered")
)

// Constructor for UserService
func NewUserService(repo repository.IUserRepo) *UserService {
	return &UserService{
		repo,
	}
}

func (service *UserService) Register(input dto.RegisterUserRequest) (models.User, error) {
	var user models.User

	isEmailAvailable, err := service.IsEmailAvailable(input.Email)
	if err != nil {
		return user, err
	}

	if !isEmailAvailable {
		return user, ErrEmailAlreadyTaken
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	user = models.User{
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

func (service *UserService) IsEmailAvailable(email string) (bool, error) {
	user, err := service.repo.FindByEmail(email)
	fmt.Printf("user: %+v", user)
	if err != nil {
		return false, err
	}

	if user.ID != 0 {
		return false, ErrEmailAlreadyTaken
	}

	return true, nil
}

func (service *UserService) UploadAvatar(ID int, fileLocation string) (models.User, error) {
	user, err := service.repo.FindByID(ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, ErrUserNotFound
		}

		return user, err
	}

	user.AvatarFileName = fileLocation

	user, err = service.repo.Update(user)
	if err != nil {
		return user, err
	}

	return user, nil
}
