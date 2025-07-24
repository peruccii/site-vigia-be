package dto

type RecoverPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}
