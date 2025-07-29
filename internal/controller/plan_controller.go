package controller

import (
	"net/http"

	"peruccii/site-vigia-be/internal/dto"
	"peruccii/site-vigia-be/internal/services"

	"github.com/gin-gonic/gin"
)

type PlanController interface {
	Create(c *gin.Context)
}

type planController struct {
	services services.PlanService
}

func NewPlanController(service services.PlanService) PlanController {
	return &planController{services: service}
}

func (ctrl *planController) Create(c *gin.Context) {
	var input dto.CreatePlanRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"invalid input": err.Error()})
	}

	if err := ctrl.services.CreatePlan(c, input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Plan created successfully"})
}
