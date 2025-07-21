package controller

import (
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

func (ctrl *subscriptionController) Create(c *gin.Context) {}
