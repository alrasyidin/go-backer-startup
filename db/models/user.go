package models

import "time"

type User struct {
	ID                int
	Name              string
	Occupation        string
	Email             string
	AvatarFileName    string
	PasswordHash      string
	Token             string
	Role              string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	PasswordChangedAt time.Time
}
