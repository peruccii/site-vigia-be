package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"peruccii/site-vigia-be/db"
	"peruccii/site-vigia-be/internal/api/middleware"
	"peruccii/site-vigia-be/internal/controller"
	"peruccii/site-vigia-be/internal/repository"
	"peruccii/site-vigia-be/internal/services"
)

func SetupRouter(database *sql.DB) *gin.Engine {
	r := gin.Default()

	queries := db.New(database)
	userRepo := repository.NewUserRepository(queries)
	authRepo := repository.NewAuthRepository(queries)
	planRepo := repository.NewPlanRepository(queries)
	subscriptionRepo := repository.NewSubscriptionRepository(queries)

	authService := services.NewAuthService(authRepo, userRepo)
	userService := services.NewUserService(userRepo)
	planService := services.NewPlanService(planRepo)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo)

	authController := controller.NewAuthController(authService)
	userController := controller.NewUserController(userService)
	planController := controller.NewPlanController(planService)
	subscriptionController := controller.NewSubscriptionController(subscriptionService)

	api := r.Group("/api")

	auth := api.Group("/auth")
	{
		auth.POST("/register", userController.Register)
		auth.POST("/login", authController.SignInUser)
	}

	protected := api.Group("")
	{
		plans := protected.Group("/plan").Use(middleware.AuthMiddleware(authService))
		{
			plans.POST("", planController.Create)
		}

		users := protected.Group("/users").Use(middleware.AuthMiddleware(authService))
		{
			//	 users.GET("", userController.)
		}

		subscriptions := protected.Group("/subscriptions").Use(middleware.AuthMiddleware(authService))
		{
			subscriptions.POST("", subscriptionController.Create)
		}
	}

	return r
}

