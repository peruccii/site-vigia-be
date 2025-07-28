package controller

import (
	"net/http"
	"peruccii/site-vigia-be/internal/dto"
	"peruccii/site-vigia-be/internal/services"

	"github.com/gin-gonic/gin"
)

type SubscriptionController interface {
	Create(c *gin.Context)
}

type subscriptionController struct {
	services services.SubscriptionService
}

func NewSubscriptionController(service services.SubscriptionService) SubscriptionController {
	return &subscriptionController{services: service}
}

func (ctrl *subscriptionController) Create(c *gin.Context) {
	var input dto.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"invalid input": err.Error()})
	}

	if err := ctrl.services.CreateSubscription(c, input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Subscription created successfully"})
}
