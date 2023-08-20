package models

import "time"

type User struct {
	ID                int       `json:"id" gorm:"primaryKey"`
	Name              string    `json:"name"`
	Occupation        string    `json:"occupation"`
	Email             string    `json:"email"`
	AvatarFileName    string    `json:"avatar_file_name"`
	PasswordHash      string    `json:"password_hash"`
	Token             string    `json:"token"`
	Role              string    `json:"role"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
}
