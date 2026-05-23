package routes

import (
	"github.com/Transcendence/config"
	"github.com/Transcendence/controllers"
	"github.com/Transcendence/middleware"
	"github.com/Transcendence/repositories"
	"github.com/Transcendence/services"
	"github.com/Transcendence/socket"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func createPostRoutes(api *gin.RouterGroup, rdb *redis.Client, postController *controllers.PostController) {
	posts := api.Group("/posts")
	{
		posts.GET("", middleware.OptionalAuthMiddleware(), postController.GetPosts)
		posts.GET("/user/:userId", middleware.OptionalAuthMiddleware(), postController.GetPostsByUser)
		posts.GET("/:id", middleware.OptionalAuthMiddleware(), postController.GetPost)
		posts.GET("/:id/comments", postController.GetComments)

		protected := posts.Group("")
		protected.Use(middleware.AuthMiddleware(rdb))
		{
			protected.POST("", postController.CreatePost)
			protected.PUT("/:id", postController.UpdatePost)
			protected.DELETE("/:id", postController.DeletePost)

			protected.POST("/:id/like", postController.ToggleLike)

			protected.POST("/:id/comments", postController.CreateComment)
			protected.PUT("/:id/comments/:commentId", postController.UpdateComment)
			protected.DELETE("/:id/comments/:commentId", postController.DeleteComment)
		}
	}
}

func SetupRoutes(router *gin.Engine, DB *gorm.DB, rdb *redis.Client, cfg *config.Config) {
	notifRepo := repositories.NewNotificationRepositories(DB)
	notifPubSub := repositories.NewNotificationPubSub(rdb)
	notifService := services.NewNotificationService(notifRepo, notifPubSub)
	notifController := controllers.NewNotificationController(notifService)

	msgRepo := repositories.NewMessageRepository(DB)
	msgController := controllers.NewMsgController(msgRepo)

	userRepo := repositories.NewUserRepository(DB)
	authService := services.NewAuthService(userRepo)
	twoFAService := services.NewTwoFAService(userRepo)
	authController := controllers.NewAuthController(authService, twoFAService, rdb)
	twoFAController := controllers.NewTwoFAController(twoFAService)

	postRepo := repositories.NewPostRepository(DB)
	postService := services.NewPostService(postRepo)

	userService := services.NewUserService(userRepo)
	friendService := &services.FriendService{DB: DB}
	friendController := &controllers.FriendController{
		Service:             friendService,
		NotificationService: notifService,
	}
	userController := controllers.NewUserController(userService, friendService)

	fileRepo := repositories.NewFileRepository(DB)
	uploadService := services.NewUploadService(fileRepo)
	uploadController := &controllers.UploadController{
		Service:       uploadService,
		FriendService: friendService,
	}

	postController := controllers.NewPostController(postService, notifService, uploadService)
	gdprService := services.NewGDPRService(DB)
	gdprController := controllers.NewGDPRController(gdprService)

	wsManager := socket.NewWSManager()
	chatHandler := socket.NewChatHandler(wsManager, rdb, notifService, msgRepo, fileRepo)

	chatService := services.NewChatService(msgRepo, userRepo)
	chatController := controllers.NewChatController(chatService)

	oauthService := services.NewOAuthService(userRepo, rdb, cfg)
	oauthController := controllers.NewOAuthController(oauthService, cfg)

	api := router.Group("/api")
	{
		api.POST("/auth/2fa/verify", middleware.RateLimitMiddleware(rdb), authController.Verify2FA)
		api.POST("/auth/register", middleware.RateLimitMiddleware(rdb), authController.RegisterUser)
		api.POST("/auth/login", middleware.RateLimitMiddleware(rdb), authController.LoginUser)
		api.POST("/auth/refresh", middleware.RateLimitMiddleware(rdb), authController.RefreshToken)
		api.GET("/ws/chat", middleware.WSAuthMiddleware(), chatHandler.HandleWS)

		api.GET("/auth/oauth/github/login", oauthController.OAuthLogin)
		api.GET("/auth/oauth/github/callback", oauthController.OAuthCallback)

		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(rdb))
		{
			protected.POST("/auth/logout", authController.LogoutUser)
			protected.GET("users", userController.GetUsers)
			protected.GET("users/:id", userController.GetUser)
			protected.PUT("users/:id", userController.UpdateUser)
			protected.DELETE("users/:id", userController.DeleteUser)

			protected.POST("friends/request/:id", friendController.SendFriendRequest)
			protected.POST("friends/accept/:id", friendController.AcceptFriend)
			protected.POST("friends/reject/:id", friendController.RejectFriendRequest)
			protected.DELETE("friends/:id", friendController.RemoveFriend)

			// Room message history — served by the WebSocket primary path
			protected.GET("rooms/:roomId/messages", msgController.GetRoomMsg)
			protected.GET("messages/:messageId/replies", msgController.GetReplies)

			// Poll/DM fallback — used when WebSocket is unavailable
			protected.POST("chat/messages", chatController.SendMessage)
			protected.GET("chat/messages", chatController.ListConversation)
			protected.GET("chat/poll", chatController.Poll)
			protected.POST("friends/follow/:id", friendController.FollowUser)
			protected.DELETE("friends/follow/:id", friendController.UnfollowUser)

			protected.GET("users/:id/followers", friendController.GetFollowers)
			protected.GET("users/:id/following", friendController.GetFollowing)
			protected.GET("users/:id/friends", friendController.GetFriends)

			protected.GET("notification", notifController.GetUnread)
			protected.PATCH("notification/read", notifController.MarkAllRead)

			protected.POST("upload", uploadController.UploadFile)
			protected.GET("/files/:id", uploadController.ServeFile)
			protected.GET("gdpr/export", gdprController.ExportUserData)
			protected.DELETE("gdpr/delete", gdprController.DeleteUserData)

			protected.POST("/2fa/setup", twoFAController.Setup)
			protected.POST("/2fa/enable", twoFAController.Enable)
			protected.POST("/2fa/disable", twoFAController.Disable)
		}

		createPostRoutes(api, rdb, postController)
	}
}
