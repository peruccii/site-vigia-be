package dto

type ResetPasswordRequest struct {
	NewPassword string `json:"new_password" validate:"required,min=8,max=100"`
}
