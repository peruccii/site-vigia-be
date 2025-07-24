package controller

import (
	"net/http"

	"peruccii/site-vigia-be/internal/dto"
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
	var input dto.SignInUserRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"invalid input": err.Error()})
		return
	}

	response, err := ctrl.services.SignInUser(c, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (ctrl *authController) RecoverPassword(c *gin.Context) {
	var input dto.RecoverPasswordRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"invalid input": err.Error()})
		return
	}
}
