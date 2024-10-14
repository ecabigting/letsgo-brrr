package routes

import (
	"github.com/ecabigting/letsgo-brrr/usermanager-api/controllers"
	"github.com/ecabigting/letsgo-brrr/usermanager-api/middlewares"
	"github.com/ecabigting/letsgo-brrr/usermanager-api/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRoutes(router *gin.Engine, db *mongo.Database) {
	userService := services.NewUserService(db)
	userController := controllers.NewUserController(userService)

	// public routes
	router.POST("/user", userController.CreateUser)
	router.PATCH("/user/verify/:userId", userController.VerifyUser)

	// protected routes
	protected := router.Group("/user").Use(middlewares.AuthMiddleware())
	{
		// get user by id with auth token
		protected.GET("/:userId", userController.GetUser)
		// update user profile
		protected.PATCH("", userController.UpdateUser)
		// delete user
		protected.DELETE("/:userId", userController.DeleteUser)
	}

	// admin routes
	admin := router.Group("/user").Use(middlewares.AuthMiddleware())
	admin.Use(middlewares.RateLimiter())
}
