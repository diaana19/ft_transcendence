package test

import (
	"testing"

	"ft_transcendence/backend/internal/models"
	"ft_transcendence/backend/internal/repositories"
	"ft_transcendence/backend/internal/services"
	"ft_transcendence/backend/internal/utils"
)

// --- Friend Service: GetFollowers ---

func TestFriendService_GetFollowers(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	// Create test users
	user1ID := utils.NewID()
	user2ID := utils.NewID()
	targetID := utils.NewID()
	user1 := models.User{ID: user1ID, Username: "follower1", Email: "f1@test.com"}
	user2 := models.User{ID: user2ID, Username: "follower2", Email: "f2@test.com"}
	target := models.User{ID: targetID, Username: "target", Email: "target@test.com"}
	db.Create(&user1)
	db.Create(&user2)
	db.Create(&target)

	// Create follows
	service.Follow(user1ID, targetID)
	service.Follow(user2ID, targetID)

	// Get followers
	followers, err := service.GetFollowers(targetID)
	if err != nil {
		t.Fatalf("GetFollowers: %v", err)
	}
	if len(followers) != 2 {
		t.Fatalf("expected 2 followers, got %d", len(followers))
	}
}

// --- Friend Service: GetFollowing ---

func TestFriendService_GetFollowing(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	// Create test users
	userID := utils.NewID()
	user1ID := utils.NewID()
	user2ID := utils.NewID()
	user := models.User{ID: userID, Username: "followinguser", Email: "fu@test.com"}
	user1 := models.User{ID: user1ID, Username: "following1", Email: "f1@test.com"}
	user2 := models.User{ID: user2ID, Username: "following2", Email: "f2@test.com"}
	db.Create(&user)
	db.Create(&user1)
	db.Create(&user2)

	// Create follows
	service.Follow(userID, user1ID)
	service.Follow(userID, user2ID)

	// Get following
	following, err := service.GetFollowing(userID)
	if err != nil {
		t.Fatalf("GetFollowing: %v", err)
	}
	if len(following) != 2 {
		t.Fatalf("expected 2 following, got %d", len(following))
	}
}

// --- Friend Service: GetFriends ---

func TestFriendService_GetFriends(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	// Create test users
	user1ID := utils.NewID()
	user2ID := utils.NewID()
	user1 := models.User{ID: user1ID, Username: "friend1", Email: "fr1@test.com"}
	user2 := models.User{ID: user2ID, Username: "friend2", Email: "fr2@test.com"}
	db.Create(&user1)
	db.Create(&user2)

	// Make friends
	service.SendRequest(user1ID, user2ID)
	service.AcceptRequest(user2ID, user1ID)

	// Get friends
	friends, err := service.GetFriends(user1ID)
	if err != nil {
		t.Fatalf("GetFriends: %v", err)
	}
	if len(friends) != 1 {
		t.Fatalf("expected 1 friend, got %d", len(friends))
	}
}

// --- Friend Service: RejectRequest ---

func TestFriendService_RejectRequest(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	// Create test users
	senderID := utils.NewID()
	receiverID := utils.NewID()
	sender := models.User{ID: senderID, Username: "sender", Email: "sender@test.com"}
	receiver := models.User{ID: receiverID, Username: "receiver", Email: "receiver@test.com"}
	db.Create(&sender)
	db.Create(&receiver)

	// Send request
	err := service.SendRequest(senderID, receiverID)
	if err != nil {
		t.Fatalf("SendRequest: %v", err)
	}

	// Reject request
	err = service.RejectRequest(receiverID, senderID)
	if err != nil {
		t.Fatalf("RejectRequest: %v", err)
	}

	// Verify not friends
	areFriends, _ := service.AreFriends(senderID, receiverID)
	if areFriends {
		t.Fatal("expected not friends after reject")
	}
}

// --- Friend Service: RemoveFriend ---

func TestFriendService_RemoveFriend(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	// Create test users
	user1ID := utils.NewID()
	user2ID := utils.NewID()
	user1 := models.User{ID: user1ID, Username: "rmfriend1", Email: "rm1@test.com"}
	user2 := models.User{ID: user2ID, Username: "rmfriend2", Email: "rm2@test.com"}
	db.Create(&user1)
	db.Create(&user2)

	// Make friends
	service.SendRequest(user1ID, user2ID)
	service.AcceptRequest(user2ID, user1ID)

	// Remove friend
	err := service.RemoveFriend(user1ID, user2ID)
	if err != nil {
		t.Fatalf("RemoveFriend: %v", err)
	}

	// Verify not friends
	areFriends, _ := service.AreFriends(user1ID, user2ID)
	if areFriends {
		t.Fatal("expected not friends after remove")
	}
}

// --- Friend Service: Unfollow ---

func TestFriendService_Unfollow(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	// Create test users
	followerID := utils.NewID()
	targetID := utils.NewID()
	follower := models.User{ID: followerID, Username: "unfollower", Email: "unf@test.com"}
	target := models.User{ID: targetID, Username: "untarget", Email: "unt@test.com"}
	db.Create(&follower)
	db.Create(&target)

	// Follow
	service.Follow(followerID, targetID)

	// Unfollow
	err := service.Unfollow(followerID, targetID)
	if err != nil {
		t.Fatalf("Unfollow: %v", err)
	}

	// Verify not following
	following, _ := service.GetFollowing(followerID)
	if len(following) != 0 {
		t.Fatalf("expected 0 following, got %d", len(following))
	}
}

// --- Auth Service: GetUserByID ---

func TestAuthService_GetUserByID(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewAuthService(repo)

	// Create test user
	userID := utils.NewID()
	user := models.User{ID: userID, Username: "getbyid", Email: "getbyid@test.com"}
	db.Create(&user)

	// Get user
	got, err := service.GetUserByID(userID)
	if err != nil {
		t.Fatalf("GetUserByID: %v", err)
	}
	if got.ID != userID {
		t.Fatalf("expected ID %s, got %s", userID, got.ID)
	}
}

// --- Auth Service: GetUserByID not found ---

func TestAuthService_GetUserByIDNotFound(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewAuthService(repo)

	_, err := service.GetUserByID(utils.NewID())
	if err == nil {
		t.Fatal("expected error for non-existent user")
	}
}

// --- User Service: GetUser ---

func TestUserService_GetUser(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewUserService(repo)

	// Create test user
	userID := utils.NewID()
	user := models.User{ID: userID, Username: "getuser", Email: "getuser@test.com"}
	db.Create(&user)

	// Get user
	got, err := service.GetUser(userID)
	if err != nil {
		t.Fatalf("GetUser: %v", err)
	}
	if got.ID != userID {
		t.Fatalf("expected ID %s, got %s", userID, got.ID)
	}
}

// --- User Service: GetUser not found ---

func TestUserService_GetUserNotFound(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewUserService(repo)

	_, err := service.GetUser(utils.NewID())
	if err == nil {
		t.Fatal("expected error for non-existent user")
	}
}

// --- User Service: UpdateUser ---

func TestUserService_UpdateUser(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewUserService(repo)

	// Create test user
	userID := utils.NewID()
	user := models.User{ID: userID, Username: "upduser", Email: "upduser@test.com"}
	db.Create(&user)

	// Update user
	updated, err := service.UpdateUser(userID, models.UpdateUserInput{Bio: "new bio"})
	if err != nil {
		t.Fatalf("UpdateUser: %v", err)
	}
	if updated.Bio != "new bio" {
		t.Fatalf("expected bio %q, got %q", "new bio", updated.Bio)
	}
}

// --- User Service: DeleteUser ---

func TestUserService_DeleteUser(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewUserService(repo)

	// Create test user
	userID := utils.NewID()
	user := models.User{ID: userID, Username: "deluser", Email: "deluser@test.com"}
	db.Create(&user)

	// Delete user
	err := service.DeleteUser(userID)
	if err != nil {
		t.Fatalf("DeleteUser: %v", err)
	}

	// Verify deleted
	_, err = service.GetUser(userID)
	if err == nil {
		t.Fatal("expected error for deleted user")
	}
}

// --- Notification Service: GetUnread ---

func TestNotificationService_GetUnread(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewNotificationRepositories(db)
	pubsub := repositories.NewNotificationPubSub(sharedRDB)
	service := services.NewNotificationService(repo, pubsub)

	// Create test user
	userID := utils.NewID()
	user := models.User{ID: userID, Username: "notifuser", Email: "notif@test.com"}
	db.Create(&user)

	// Send notification
	service.SendNotification(userID, "", utils.NewID(), "actor", "test", "test content", "")

	// Get unread
	notifs, err := service.GetUnread(userID)
	if err != nil {
		t.Fatalf("GetUnread: %v", err)
	}
	if len(notifs) != 1 {
		t.Fatalf("expected 1 notification, got %d", len(notifs))
	}
}

// --- Notification Service: MarkAllRead ---

func TestNotificationService_MarkAllRead(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewNotificationRepositories(db)
	pubsub := repositories.NewNotificationPubSub(sharedRDB)
	service := services.NewNotificationService(repo, pubsub)

	// Create test user
	userID := utils.NewID()
	user := models.User{ID: userID, Username: "markread", Email: "markread@test.com"}
	db.Create(&user)

	// Send notifications
	service.SendNotification(userID, "", utils.NewID(), "actor1", "test", "content1", "")
	service.SendNotification(userID, "", utils.NewID(), "actor2", "test", "content2", "")

	// Mark all read
	err := service.MarkAllRead(userID)
	if err != nil {
		t.Fatalf("MarkAllRead: %v", err)
	}

	// Verify
	notifs, _ := service.GetUnread(userID)
	if len(notifs) != 0 {
		t.Fatalf("expected 0 unread, got %d", len(notifs))
	}
}

// --- Notification Service: MarkRead ---

func TestNotificationService_MarkRead(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewNotificationRepositories(db)
	pubsub := repositories.NewNotificationPubSub(sharedRDB)
	service := services.NewNotificationService(repo, pubsub)

	// Create test user
	userID := utils.NewID()
	user := models.User{ID: userID, Username: "markreadsingle", Email: "markreadsingle@test.com"}
	db.Create(&user)

	// Send notification
	service.SendNotification(userID, "", utils.NewID(), "actor", "test", "content", "")

	// Get notification
	notifs, _ := service.GetUnread(userID)
	if len(notifs) != 1 {
		t.Fatalf("expected 1 notification")
	}

	// Mark read
	err := service.MarkRead(userID, notifs[0].ID)
	if err != nil {
		t.Fatalf("MarkRead: %v", err)
	}

	// Verify
	notifs, _ = service.GetUnread(userID)
	if len(notifs) != 0 {
		t.Fatalf("expected 0 unread, got %d", len(notifs))
	}
}
