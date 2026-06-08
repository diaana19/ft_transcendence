package routes

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"ft_transcendence/backend/internal/config"
	"ft_transcendence/backend/internal/controllers"
	"ft_transcendence/backend/internal/repositories"
	"ft_transcendence/backend/internal/services"
	"ft_transcendence/backend/internal/socket"
)

type Controllers struct {
	Auth         *controllers.AuthController
	TwoFA        *controllers.TwoFAController
	User         *controllers.UserController
	Friend       *controllers.FriendController
	Post         *controllers.PostController
	Notification *controllers.NotificationController
	Upload       *controllers.UploadController
	GDPR         *controllers.GDPRController
	OAuth        *controllers.OAuthController
	Gamification *controllers.GamificationController
	ChatWS       *socket.ChatHandler
	WSManager    *socket.WSManager
}

func Wire(pdb *gorm.DB, rdb *redis.Client, cfg *config.Config) *Controllers {
	userRepo := repositories.NewUserRepository(pdb)
	msgRepo := repositories.NewMessageRepository(pdb)
	postRepo := repositories.NewPostRepository(pdb)
	fileRepo := repositories.NewFileRepository(pdb)
	notifRepo := repositories.NewNotificationRepositories(pdb)
	notifPubSub := repositories.NewNotificationPubSub(rdb)

	authService := services.NewAuthService(userRepo)
	twoFAService := services.NewTwoFAService(userRepo)
	mailService := services.NewMailService()
	userService := services.NewUserService(userRepo)
	friendService := &services.FriendService{DB: pdb}
	postService := services.NewPostService(postRepo)
	notifService := services.NewNotificationService(notifRepo, notifPubSub)
	uploadService := services.NewUploadService(fileRepo)
	gdprService := services.NewGDPRService(pdb)
	oauthService := services.NewOAuthService(userRepo, rdb, cfg)
	gamificationService := services.NewGamificationService(pdb, friendService)

	wsManager := socket.NewWSManager()
	chatWS := socket.NewChatHandler(wsManager, rdb, notifService, msgRepo, userRepo, fileRepo, cfg.FrontendURL)

	return &Controllers{
		Auth:         controllers.NewAuthController(authService, twoFAService, mailService, rdb),
		TwoFA:        controllers.NewTwoFAController(twoFAService),
		User:         controllers.NewUserController(userService, friendService, mailService),
		Friend:       &controllers.FriendController{Service: friendService, NotificationService: notifService},
		Post:         controllers.NewPostController(postService, notifService, uploadService),
		Notification: controllers.NewNotificationController(notifService),
		Upload:       &controllers.UploadController{Service: uploadService, FriendService: friendService},
		GDPR:         controllers.NewGDPRController(gdprService, mailService),
		OAuth:        controllers.NewOAuthController(oauthService, cfg),
		Gamification: controllers.NewGamificationController(gamificationService),
		ChatWS:       chatWS,
		WSManager:    wsManager,
	}
}
