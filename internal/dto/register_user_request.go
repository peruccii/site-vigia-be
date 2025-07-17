package dto

type RegisterUserRequest struct {
	Name     string `json:"name" validate:"required,min=1,max=255"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=128"`
}
