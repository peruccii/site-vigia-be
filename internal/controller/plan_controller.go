package controller

import (
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
}
