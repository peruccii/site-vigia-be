package controller

import (
	"net/http"

	"peruccii/site-vigia-be/db"
	"peruccii/site-vigia-be/internal/dto"
	"peruccii/site-vigia-be/internal/services"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	Register(c *gin.Context)
	FindByEmail(c *gin.Context)
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

	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}
	user, err := ctrl.service.FindByEmail(c, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == (db.User{}) {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	response := dto.UserResponse{
		ID:              user.ID.String(),
		Name:            user.Name,
		Email:           user.Email,
		EmailVerifiedAt: user.EmailVerifiedAt,
		CreatedAt:       user.CreatedAt.Local().String(),
		UpdatedAt:       user.UpdatedAt.Local().String(),
	}

	c.JSON(http.StatusOK, response)
}
