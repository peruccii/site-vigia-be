package controller

import (
	"net/http"
	"strconv"

	"peruccii/site-vigia-be/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SubscriptionHandler struct {
	stripeService services.StripeService
}

func NewSubscriptionHandler(stripeService services.StripeService) *SubscriptionHandler {
	return &SubscriptionHandler{
		stripeService: stripeService,
	}
}

// POST /subscriptions/checkout
func (h *SubscriptionHandler) CreateCheckout(c *gin.Context) {
	var req struct {
		PlanID string `json:"plan_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Pegar dados do usuário autenticado
	userID := c.MustGet("user_id").(uuid.UUID)
	userName := c.MustGet("user_name").(string)
	userEmail := c.MustGet("user_email").(string)

	planID, err := strconv.Atoi(req.PlanID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan_id"})
		return
	}

	checkoutReq := &services.CreateSubscriptionRequest{
		UserID:     userID,
		PlanID:     int32(planID),
		UserName:   userName,
		UserEmail:  userEmail,
		SuccessURL: "https://app.sentinelsimples.com/success",
		CancelURL:  "https://app.sentinelsimples.com/pricing",
	}

	session, err := h.stripeService.CreateSubscriptionCheckout(c.Request.Context(), checkoutReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"checkout_url": session.URL,
		"session_id":   session.ID,
	})
}

/*
// POST /subscriptions/cancel
func (h *SubscriptionHandler) CancelSubscription(c *gin.Context) {
	var req struct {
		SubscriptionID string `json:"subscription_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.stripeService.CancelSubscription(c.Request.Context(), req.SubscriptionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Subscription will be canceled at the end of current period",
	})
}*/

// POST /webhook/stripe
func (h *SubscriptionHandler) HandleWebhook(c *gin.Context) {
	payload, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	signature := c.GetHeader("Stripe-Signature")

	err = h.stripeService.ProcessWebhook(c.Request.Context(), payload, signature)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"received": true})
}
