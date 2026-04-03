package routes

import (
	"github.com/Transcendence/controllers"
	"github.com/Transcendence/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, DB *gorm.DB) {
	
	userService := services.NewUserService(DB)
	userController := controllers.NewUserController(userService)
	
	// Routes
	api := router.Group("/api")
	{
		api.POST("/auth/register", controllers.RegisterUser)
		api.GET("/users",	userController.GetUsers)
		api.POST("/users",	userController.CreateUser)
		api.GET("/users/:id",	userController.GetUser)
		api.PUT("/users/:id",	userController.UpdateUser)
		api.DELETE("/users/:id",	userController.DeleteUser)
	}
}