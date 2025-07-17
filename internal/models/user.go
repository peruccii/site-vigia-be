package models

import (
"github.com/google/uuid"
	"time"
)

type User struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	Name            string     `json:"name" db:"name"`
	Email           string     `json:"email" db:"email"`
	PasswordHash    string     `json:"-" db:"password_hash"`
	EmailVerifiedAt *time.Time `json:"email_verified_at" db:"email_verified_at"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}
