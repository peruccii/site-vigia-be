package controller

import (
	"net/http"

	"peruccii/site-vigia-be/internal/dto"
	"peruccii/site-vigia-be/internal/services"
	"peruccii/site-vigia-be/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthController interface {
	SignInUser(c *gin.Context)
	RecoverPassword(c *gin.Context)
	ResetPassword(c *gin.Context)
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

	response, err := ctrl.services.RecoverPassword(c, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (ctrl *authController) ResetPassword(c *gin.Context) {
	var input dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"invalid input": err.Error()})
		return
	}
	token := c.Request.URL.Query().Get("token")
	claims, err := utils.VerifyToken(token)
	if err != nil {
		http.Error(c.Writer, "Token inválido", http.StatusUnauthorized)
		return
	}
	_ = ctrl.services.ResetPassword(c, claims["email"].(string), input)

	c.JSON(http.StatusOK, gin.H{
		"message": "Senha atualizada com succeso",
	})
}
