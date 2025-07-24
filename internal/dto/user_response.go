package dto

import "time"

type UserResponse struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	Email           string            `json:"email"`
	EmailVerifiedAt *time.Time        `json:"email_verified_at"`
	CreatedAt       string            `json:"created_at"`
	UpdatedAt       string            `json:"updated_at"`
	Websites        []WebsiteResponse `json:"websites"`
}
