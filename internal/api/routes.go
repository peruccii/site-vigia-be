package api

import (
	"database/sql"

	"peruccii/site-vigia-be/db"
	"peruccii/site-vigia-be/internal/api/middleware"
	"peruccii/site-vigia-be/internal/controller"
	"peruccii/site-vigia-be/internal/repository"
	"peruccii/site-vigia-be/internal/services"

	"github.com/gin-gonic/gin"
)

func SetupRouter(database *sql.DB) *gin.Engine {
	r := gin.Default()

	queries := db.New(database)
	userRepo := repository.NewUserRepository(queries)
	authRepo := repository.NewAuthRepository(queries)
	planRepo := repository.NewPlanRepository(queries)
	paymentRepo := repository.NewPaymentRepository(queries)
	subscriptionRepo := repository.NewSubscriptionRepository(queries)

	authService := services.NewAuthService(authRepo, userRepo)
	userService := services.NewUserService(userRepo)
	planService := services.NewPlanService(planRepo)
	// subscriptionService := services.NewSubscriptionService(subscriptionRepo)
	stripeService := services.NewStripeService(paymentRepo, *subscriptionRepo, *planRepo)

	authController := controller.NewAuthController(authService)
	userController := controller.NewUserController(userService)
	planController := controller.NewPlanController(planService)
	subscriptionController := controller.NewSubscriptionHandler(stripeService)

	api := r.Group("/api")

	auth := api.Group("/auth")
	{
		auth.POST("/register", userController.Register)
		auth.POST("/login", authController.SignInUser)
		auth.PUT("/recover-password", authController.RecoverPassword).Use(middleware.AuthMiddleware(authService))
		auth.PATCH("/reset-password", authController.ResetPassword).Use(middleware.AuthMiddleware(authService))
	}

	protected := api.Group("")
	{
		users := protected.Group("/users").Use(middleware.AuthMiddleware(authService))
		{
			users.GET("/:email", userController.FindByEmail)
		}

		websites := protected.Group("/websites").Use(middleware.AuthMiddleware(authService))
		{
			websites.POST("")
		}

		plans := protected.Group("/plan").Use(middleware.AuthMiddleware(authService))
		{
			plans.POST("", planController.Create)
		}

		subscriptions := protected.Group("/subscriptions").Use(middleware.AuthMiddleware(authService))
		{
			subscriptions.POST("", subscriptionController.CreateCheckout)
		}
	}

	return r
}
