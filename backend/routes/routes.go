package routes

import (
	"github.com/Transcendence/controllers"
	"github.com/Transcendence/middleware"
	"github.com/Transcendence/repositories"
	"github.com/Transcendence/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, DB *gorm.DB) {

	userRepo := repositories.NewUserRepository(DB)
	authService := services.NewAuthService(userRepo)
	authController := controllers.NewAuthController(authService)

	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	messageRepo := repositories.NewMessageRepository(DB)
	chatService := services.NewChatService(messageRepo, userRepo)
	chatController := controllers.NewChatController(chatService)

	api := router.Group("/api")
	{
		api.POST("/auth/register", authController.RegisterUser)
		api.POST("/auth/login", authController.LoginUser)
		api.POST("/auth/refresh", authController.RefreshToken)
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			api.GET("/users", userController.GetUsers)
			api.GET("/users/:id", userController.GetUser)
			api.PUT("/users/:id", userController.UpdateUser)
			api.DELETE("/users/:id", userController.DeleteUser)

			protected.POST("/chat/messages", chatController.SendMessage)
			protected.GET("/chat/messages", chatController.ListConversation)
			protected.GET("/chat/poll", chatController.Poll)
		}
	}
}
