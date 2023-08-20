package dto

import "github.com/alrasyidin/bwa-backer-startup/db/models"

type RegisterUserRequest struct {
	Name       string `json:"name" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Occupation string `json:"occupation" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type RegisterUserResponse struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Occupation     string `json:"occupation"`
	AvatarFileName string `json:"avatar_file_name"`
	Token          string `json:"token"`
}

func FormatUser(user models.User, token string) RegisterUserResponse {
	responseUser := RegisterUserResponse{
		ID:             user.ID,
		Name:           user.Name,
		Email:          user.Email,
		Occupation:     user.Occupation,
		AvatarFileName: user.AvatarFileName,
		Token:          token,
	}

	return responseUser
}
