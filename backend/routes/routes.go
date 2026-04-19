package routes

import (
	"github.com/Transcendence/controllers"
	"github.com/Transcendence/middleware"
	"github.com/Transcendence/repositories"
	"github.com/Transcendence/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func create_post_routes(api *gin.RouterGroup, DB *gorm.DB) {
	postRepo := repositories.NewPostRepository(DB)
	postService := services.NewPostService(postRepo)
	postController := controllers.NewPostController(postService)

	posts := api.Group("/posts")
	{
		// Public – read (liked field populated only when token is present)
		posts.GET("", middleware.OptionalAuthMiddleware(), postController.GetPosts)
		posts.GET("/:id", middleware.OptionalAuthMiddleware(), postController.GetPost)
		posts.GET("/:id/comments", postController.GetComments)

		// Protected – require authentication
		protected := posts.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// Post CRUD
			protected.POST("", postController.CreatePost)
			protected.PUT("/:id", postController.UpdatePost)
			protected.DELETE("/:id", postController.DeletePost)

			// Likes
			protected.POST("/:id/like", postController.ToggleLike)

			// Comments
			protected.POST("/:id/comments", postController.CreateComment)
			protected.PUT("/:id/comments/:commentId", postController.UpdateComment)
			protected.DELETE("/:id/comments/:commentId", postController.DeleteComment)
		}
	}
}

func SetupRoutes(router *gin.Engine, DB *gorm.DB) {

	userRepo := repositories.NewUserRepository(DB)
	authService := services.NewAuthService(userRepo)
	authController := controllers.NewAuthController(authService)

	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	api := router.Group("/api")
	{
		api.POST("/auth/register", authController.RegisterUser)
		api.POST("/auth/login", authController.LoginUser)
		api.POST("/auth/refresh", authController.RefreshToken)

		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/users", userController.GetUsers)
			protected.GET("/users/:id", userController.GetUser)
			protected.PUT("/users/:id", userController.UpdateUser)
			protected.DELETE("/users/:id", userController.DeleteUser)
		}

		create_post_routes(api, DB)
	}
}
