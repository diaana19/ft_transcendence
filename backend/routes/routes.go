package routes

import (
	"github.com/Transcendence/controllers"
	"github.com/Transcendence/repositories"
	"github.com/Transcendence/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, DB *gorm.DB) {

	userRepo := repositories.NewUserRepository(DB)
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	// Routes
	api := router.Group("/api")
	{
		api.POST("/auth/register", func(c *gin.Context) {
			controllers.RegisterUser(c, DB)
		})
		api.GET("/users", userController.GetUsers)
		api.GET("/users/:id", userController.GetUser)
		api.PUT("/users/:id", userController.UpdateUser)
		api.DELETE("/users/:id", userController.DeleteUser)
	}
}
