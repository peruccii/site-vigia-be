package controller

import (
	"net/http"

	"peruccii/site-vigia-be/internal/dto"
	"peruccii/site-vigia-be/internal/services"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	Register(c *gin.Context)
}

type userController struct {
	service services.UserService
}

func NewUserController(service services.UserService) UserController {
	return &userController{service: service}
}

func (ctrl *userController) Register(c *gin.Context) {
	var input dto.RegisterUserRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"invalid input": err.Error()})
		return
	}

	if err := ctrl.service.RegisterUser(c, input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, "CREATED")
}

func (ctrl *userController) FindByEmail(c *gin.Context) {
    email := c.Param("email")

    user, err := ctrl.service.
}
