package routes

import (
	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
	"github.com/Transcendence/controllers"
)

func SetupRoutes(router *gin.Engine, DB *gorm.DB) {

	// Routes
	api := router.Group("/api")
	{
		api.POST("/auth/register", controllers.RegisterUser)
	}
}