package dto

type EmailCheckRequest struct {
	Email string `json:"email" binding:"required,email"`
}
