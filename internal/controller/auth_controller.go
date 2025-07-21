package controller

import (
	"peruccii/site-vigia-be/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthController interface {
	SignInUser(c *gin.Context)
}

type authController struct {
	services services.AuthService
}

func NewAuthController(services services.AuthService) AuthController {
	return &authController{services: services}
}

func (ctrl *authController) SignInUser(c *gin.Context) {
}
