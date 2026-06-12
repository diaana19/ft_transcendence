package test

import (
	"testing"

	"ft_transcendence/backend/internal/models"
	"ft_transcendence/backend/internal/repositories"
	"ft_transcendence/backend/internal/services"
	"ft_transcendence/backend/internal/utils"
)

func TestFriendService_GetFollowers(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	user1ID := utils.NewID()
	user2ID := utils.NewID()
	targetID := utils.NewID()
	user1 := models.User{ID: user1ID, Username: "follower1", Email: "f1@test.com"}
	user2 := models.User{ID: user2ID, Username: "follower2", Email: "f2@test.com"}
	target := models.User{ID: targetID, Username: "target", Email: "target@test.com"}
	db.Create(&user1)
	db.Create(&user2)
	db.Create(&target)

	service.Follow(user1ID, targetID)
	service.Follow(user2ID, targetID)

	followers, err := service.GetFollowers(targetID)
	if err != nil {
		t.Fatalf("GetFollowers: %v", err)
	}
	if len(followers) != 2 {
		t.Fatalf("expected 2 followers, got %d", len(followers))
	}
}

func TestFriendService_GetFollowing(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	userID := utils.NewID()
	user1ID := utils.NewID()
	user2ID := utils.NewID()
	user := models.User{ID: userID, Username: "followinguser", Email: "fu@test.com"}
	user1 := models.User{ID: user1ID, Username: "following1", Email: "f1@test.com"}
	user2 := models.User{ID: user2ID, Username: "following2", Email: "f2@test.com"}
	db.Create(&user)
	db.Create(&user1)
	db.Create(&user2)

	service.Follow(userID, user1ID)
	service.Follow(userID, user2ID)

	following, err := service.GetFollowing(userID)
	if err != nil {
		t.Fatalf("GetFollowing: %v", err)
	}
	if len(following) != 2 {
		t.Fatalf("expected 2 following, got %d", len(following))
	}
}

func TestFriendService_GetFriends(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	user1ID := utils.NewID()
	user2ID := utils.NewID()
	user1 := models.User{ID: user1ID, Username: "friend1", Email: "fr1@test.com"}
	user2 := models.User{ID: user2ID, Username: "friend2", Email: "fr2@test.com"}
	db.Create(&user1)
	db.Create(&user2)

	service.SendRequest(user1ID, user2ID)
	service.AcceptRequest(user2ID, user1ID)

	friends, err := service.GetFriends(user1ID)
	if err != nil {
		t.Fatalf("GetFriends: %v", err)
	}
	if len(friends) != 1 {
		t.Fatalf("expected 1 friend, got %d", len(friends))
	}
}

func TestFriendService_RejectRequest(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	senderID := utils.NewID()
	receiverID := utils.NewID()
	sender := models.User{ID: senderID, Username: "sender", Email: "sender@test.com"}
	receiver := models.User{ID: receiverID, Username: "receiver", Email: "receiver@test.com"}
	db.Create(&sender)
	db.Create(&receiver)

	_, err := service.SendRequest(senderID, receiverID)
	if err != nil {
		t.Fatalf("SendRequest: %v", err)
	}

	err = service.RejectRequest(receiverID, senderID)
	if err != nil {
		t.Fatalf("RejectRequest: %v", err)
	}

	areFriends, _ := service.AreFriends(senderID, receiverID)
	if areFriends {
		t.Fatal("expected not friends after reject")
	}
}

func TestFriendService_RemoveFriend(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	user1ID := utils.NewID()
	user2ID := utils.NewID()
	user1 := models.User{ID: user1ID, Username: "rmfriend1", Email: "rm1@test.com"}
	user2 := models.User{ID: user2ID, Username: "rmfriend2", Email: "rm2@test.com"}
	db.Create(&user1)
	db.Create(&user2)

	service.SendRequest(user1ID, user2ID)
	service.AcceptRequest(user2ID, user1ID)

	err := service.RemoveFriend(user1ID, user2ID)
	if err != nil {
		t.Fatalf("RemoveFriend: %v", err)
	}

	areFriends, _ := service.AreFriends(user1ID, user2ID)
	if areFriends {
		t.Fatal("expected not friends after remove")
	}
}

func TestFriendService_Unfollow(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	followerID := utils.NewID()
	targetID := utils.NewID()
	follower := models.User{ID: followerID, Username: "unfollower", Email: "unf@test.com"}
	target := models.User{ID: targetID, Username: "untarget", Email: "unt@test.com"}
	db.Create(&follower)
	db.Create(&target)

	service.Follow(followerID, targetID)

	err := service.Unfollow(followerID, targetID)
	if err != nil {
		t.Fatalf("Unfollow: %v", err)
	}

	following, _ := service.GetFollowing(followerID)
	if len(following) != 0 {
		t.Fatalf("expected 0 following, got %d", len(following))
	}
}

func TestAuthService_GetUserByID(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewAuthService(repo)

	userID := utils.NewID()
	user := models.User{ID: userID, Username: "getbyid", Email: "getbyid@test.com"}
	db.Create(&user)

	got, err := service.GetUserByID(userID)
	if err != nil {
		t.Fatalf("GetUserByID: %v", err)
	}
	if got.ID != userID {
		t.Fatalf("expected ID %s, got %s", userID, got.ID)
	}
}

func TestAuthService_GetUserByIDNotFound(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewAuthService(repo)

	_, err := service.GetUserByID(utils.NewID())
	if err == nil {
		t.Fatal("expected error for non-existent user")
	}
}

func TestUserService_GetUser(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewUserService(repo)

	userID := utils.NewID()
	user := models.User{ID: userID, Username: "getuser", Email: "getuser@test.com"}
	db.Create(&user)

	got, err := service.GetUser(userID)
	if err != nil {
		t.Fatalf("GetUser: %v", err)
	}
	if got.ID != userID {
		t.Fatalf("expected ID %s, got %s", userID, got.ID)
	}
}

func TestUserService_GetUserNotFound(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewUserService(repo)

	_, err := service.GetUser(utils.NewID())
	if err == nil {
		t.Fatal("expected error for non-existent user")
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewUserService(repo)

	userID := utils.NewID()
	user := models.User{ID: userID, Username: "upduser", Email: "upduser@test.com"}
	db.Create(&user)

	updated, err := service.UpdateUser(userID, models.UpdateUserInput{Bio: "new bio"})
	if err != nil {
		t.Fatalf("UpdateUser: %v", err)
	}
	if updated.Bio != "new bio" {
		t.Fatalf("expected bio %q, got %q", "new bio", updated.Bio)
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewUserService(repo)

	userID := utils.NewID()
	user := models.User{ID: userID, Username: "deluser", Email: "deluser@test.com"}
	db.Create(&user)

	err := service.DeleteUser(userID)
	if err != nil {
		t.Fatalf("DeleteUser: %v", err)
	}

	_, err = service.GetUser(userID)
	if err == nil {
		t.Fatal("expected error for deleted user")
	}
}

func TestNotificationService_GetUnread(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewNotificationRepositories(db)
	pubsub := repositories.NewNotificationPubSub(sharedRDB)
	service := services.NewNotificationService(repo, pubsub)

	userID := utils.NewID()
	user := models.User{ID: userID, Username: "notifuser", Email: "notif@test.com"}
	db.Create(&user)

	service.SendNotification(userID, "", utils.NewID(), "actor", "test", "test content", "")

	notifs, err := service.GetUnread(userID)
	if err != nil {
		t.Fatalf("GetUnread: %v", err)
	}
	if len(notifs) != 1 {
		t.Fatalf("expected 1 notification, got %d", len(notifs))
	}
}

func TestNotificationService_MarkAllRead(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewNotificationRepositories(db)
	pubsub := repositories.NewNotificationPubSub(sharedRDB)
	service := services.NewNotificationService(repo, pubsub)

	userID := utils.NewID()
	user := models.User{ID: userID, Username: "markread", Email: "markread@test.com"}
	db.Create(&user)

	service.SendNotification(userID, "", utils.NewID(), "actor1", "test", "content1", "")
	service.SendNotification(userID, "", utils.NewID(), "actor2", "test", "content2", "")

	err := service.MarkAllRead(userID)
	if err != nil {
		t.Fatalf("MarkAllRead: %v", err)
	}

	notifs, _ := service.GetUnread(userID)
	if len(notifs) != 0 {
		t.Fatalf("expected 0 unread, got %d", len(notifs))
	}
}

func TestNotificationService_MarkRead(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewNotificationRepositories(db)
	pubsub := repositories.NewNotificationPubSub(sharedRDB)
	service := services.NewNotificationService(repo, pubsub)

	userID := utils.NewID()
	user := models.User{ID: userID, Username: "markreadsingle", Email: "markreadsingle@test.com"}
	db.Create(&user)

	service.SendNotification(userID, "", utils.NewID(), "actor", "test", "content", "")

	notifs, _ := service.GetUnread(userID)
	if len(notifs) != 1 {
		t.Fatalf("expected 1 notification")
	}

	err := service.MarkRead(userID, notifs[0].ID)
	if err != nil {
		t.Fatalf("MarkRead: %v", err)
	}

	notifs, _ = service.GetUnread(userID)
	if len(notifs) != 0 {
		t.Fatalf("expected 0 unread, got %d", len(notifs))
	}
}
