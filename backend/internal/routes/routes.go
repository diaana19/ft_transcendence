package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"ft_transcendence/backend/internal/controllers"
	"ft_transcendence/backend/internal/middleware"
	"ft_transcendence/backend/internal/socket"
)

// SetupRoutes registers every HTTP and WebSocket endpoint under /api. The
// route tree is broken up by feature so this function reads as a manifest;
// the per-feature register* helpers below own the path details. rdb is passed
// for the middlewares that need direct Redis access (rate limit, auth session
// lookups).
func SetupRoutes(router *gin.Engine, c *Controllers, rdb *redis.Client) {
	api := router.Group("/api")

	registerPublicAuthRoutes(api, c.Auth, rdb)
	registerOAuthRoutes(api, c.OAuth)
	registerWebSocketRoutes(api, c.ChatWS)
	registerPostRoutes(api, rdb, c.Post)

	protected := api.Group("", middleware.AuthMiddleware(rdb))
	registerProtectedAuthRoutes(protected, c.Auth)
	registerUserRoutes(protected, c.User)
	registerFriendRoutes(protected, c.Friend)
	registerNotificationRoutes(protected, c.Notification)
	registerUploadRoutes(api, rdb, c.Upload)
	registerGDPRRoutes(protected, c.GDPR)
	registerTwoFARoutes(protected, c.TwoFA)
	registerGamificationRoutes(protected, c.Gamification)
}

// /auth/{register,login,refresh,2fa/verify} — public, rate-limited per IP.
func registerPublicAuthRoutes(api *gin.RouterGroup, c *controllers.AuthController, rdb *redis.Client) {
	auth := api.Group("/auth", middleware.RateLimitMiddleware(rdb))
	auth.POST("/register", c.RegisterUser)
	auth.POST("/login", c.LoginUser)
	auth.POST("/refresh", c.RefreshToken)
	auth.POST("/2fa/verify", c.Verify2FA)
}

// GitHub OAuth — public login surface, currently not rate-limited.
func registerOAuthRoutes(api *gin.RouterGroup, c *controllers.OAuthController) {
	oauth := api.Group("/auth/oauth/github")
	oauth.GET("/login", c.OAuthLogin)
	oauth.GET("/callback", c.OAuthCallback)
}

// WebSocket chat — auth happens inside WSAuthMiddleware (reads token from a
// query param because browsers can't set Authorization on a WS handshake).
func registerWebSocketRoutes(api *gin.RouterGroup, c *socket.ChatHandler) {
	api.GET("/ws/chat", middleware.WSAuthMiddleware(), c.HandleWS)
}

func registerProtectedAuthRoutes(protected *gin.RouterGroup, c *controllers.AuthController) {
	protected.POST("/auth/logout", c.LogoutUser)
	protected.GET("/auth/me", c.Me)
}

func registerUserRoutes(protected *gin.RouterGroup, c *controllers.UserController) {
	protected.GET("/users", c.GetUsers)
	protected.GET("/users/:id", c.GetUser)
	protected.PUT("/users/:id", c.UpdateUser)
	protected.DELETE("/users/:id", c.DeleteUser)
}

func registerFriendRoutes(protected *gin.RouterGroup, c *controllers.FriendController) {
	protected.POST("/friends/request/:id", c.SendFriendRequest)
	protected.POST("/friends/accept/:id", c.AcceptFriend)
	protected.POST("/friends/reject/:id", c.RejectFriendRequest)
	protected.DELETE("/friends/:id", c.RemoveFriend)
	protected.POST("/friends/follow/:id", c.FollowUser)
	protected.DELETE("/friends/follow/:id", c.UnfollowUser)
	protected.GET("/users/:id/followers", c.GetFollowers)
	protected.GET("/users/:id/following", c.GetFollowing)
	protected.GET("/users/:id/friends", c.GetFriends)
}

func registerNotificationRoutes(protected *gin.RouterGroup, c *controllers.NotificationController) {
	protected.GET("/notification", c.GetUnread)
	protected.PATCH("/notification/read", c.MarkAllRead)
	protected.PATCH("/notification/:id/read", c.MarkRead)
}

// Uploads have mixed visibility. POST /upload requires auth (only logged-in
// users can write). GET /files/:id is mounted on the public api group with
// OptionalAuthMiddleware so the controller's visibility check runs first:
// public files stream without a token, friends/private files still require
// an authenticated caller.
func registerUploadRoutes(api *gin.RouterGroup, rdb *redis.Client, c *controllers.UploadController) {
	api.GET("/files/:id", middleware.OptionalAuthMiddleware(rdb), c.ServeFile)

	protected := api.Group("", middleware.AuthMiddleware(rdb))
	protected.POST("/upload", c.UploadFile)
}

func registerGDPRRoutes(protected *gin.RouterGroup, c *controllers.GDPRController) {
	protected.GET("/gdpr/export", c.ExportUserData)
	protected.DELETE("/gdpr/delete", c.DeleteUserData)
}

func registerTwoFARoutes(protected *gin.RouterGroup, c *controllers.TwoFAController) {
	protected.POST("/2fa/setup", c.Setup)
	protected.POST("/2fa/enable", c.Enable)
	protected.POST("/2fa/disable", c.Disable)
}

func registerGamificationRoutes(protected *gin.RouterGroup, c *controllers.GamificationController) {
	protected.GET("/users/:id/gamification", c.Get)
	protected.GET("/leaderboard", c.GetLeaderboard)
}

// Posts have mixed visibility: list / single / comments are public with optional
// auth, mutations require it. Kept self-contained because of that quirk.
func registerPostRoutes(api *gin.RouterGroup, rdb *redis.Client, c *controllers.PostController) {
	posts := api.Group("/posts")
	posts.GET("", middleware.OptionalAuthMiddleware(rdb), c.GetPosts)
	posts.GET("/user/:userId", middleware.OptionalAuthMiddleware(rdb), c.GetPostsByUser)
	posts.GET("/:id", middleware.OptionalAuthMiddleware(rdb), c.GetPost)
	posts.GET("/:id/comments", c.GetComments)

	protected := posts.Group("", middleware.AuthMiddleware(rdb))
	protected.POST("", c.CreatePost)
	protected.PUT("/:id", c.UpdatePost)
	protected.DELETE("/:id", c.DeletePost)
	protected.POST("/:id/like", c.ToggleLike)
	protected.POST("/:id/comments", c.CreateComment)
	protected.PUT("/:id/comments/:commentId", c.UpdateComment)
	protected.DELETE("/:id/comments/:commentId", c.DeleteComment)
}
