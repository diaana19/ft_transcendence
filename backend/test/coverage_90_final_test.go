package test

import (
	"testing"

	"ft_transcendence/backend/internal/models"
	"ft_transcendence/backend/internal/repositories"
	"ft_transcendence/backend/internal/services"
	"ft_transcendence/backend/internal/utils"
)

// --- Post Service: GetPosts ---

func TestPostService_GetPostsFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	// Create test user and posts
	userID := utils.NewID()
	user := models.User{ID: userID, Username: "getpostsfinal", Email: "getpostsfinal@test.com"}
	post1 := models.Post{ID: utils.NewID(), AuthorID: userID, Content: "post1"}
	post2 := models.Post{ID: utils.NewID(), AuthorID: userID, Content: "post2"}
	db.Create(&user)
	repo.Create(&post1)
	repo.Create(&post2)

	// Get posts
	posts, total, err := service.GetPosts(10, 0)
	if err != nil {
		t.Fatalf("GetPosts: %v", err)
	}
	if len(posts) != 2 {
		t.Fatalf("expected 2 posts, got %d", len(posts))
	}
	if total != 2 {
		t.Fatalf("expected total 2, got %d", total)
	}
}

// --- Post Service: GetPost ---

func TestPostService_GetPostFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	// Create test user and post
	userID := utils.NewID()
	postID := utils.NewID()
	user := models.User{ID: userID, Username: "getpostfinal", Email: "getpostfinal@test.com"}
	post := models.Post{ID: postID, AuthorID: userID, Content: "post content"}
	db.Create(&user)
	repo.Create(&post)

	// Get post
	got, err := service.GetPost(postID)
	if err != nil {
		t.Fatalf("GetPost: %v", err)
	}
	if got.Content != "post content" {
		t.Fatalf("expected 'post content', got %q", got.Content)
	}
}

// --- Post Service: GetPost not found ---

func TestPostService_GetPostNotFoundFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	_, err := service.GetPost(utils.NewID())
	if err == nil {
		t.Fatal("expected error for non-existent post")
	}
}

// --- Post Service: GetPostsByAuthor ---

func TestPostService_GetPostsByAuthorFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	// Create test user and posts
	userID := utils.NewID()
	user := models.User{ID: userID, Username: "getpostsbyauthorfinal", Email: "getpostsbyauthorfinal@test.com"}
	post1 := models.Post{ID: utils.NewID(), AuthorID: userID, Content: "post1"}
	post2 := models.Post{ID: utils.NewID(), AuthorID: userID, Content: "post2"}
	db.Create(&user)
	repo.Create(&post1)
	repo.Create(&post2)

	// Get posts by author
	posts, err := service.GetPostsByAuthor(userID)
	if err != nil {
		t.Fatalf("GetPostsByAuthor: %v", err)
	}
	if len(posts) != 2 {
		t.Fatalf("expected 2 posts, got %d", len(posts))
	}
}

// --- Post Service: GetTrends ---

func TestPostService_GetTrendsFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	// Create test user and posts with tags
	userID := utils.NewID()
	user := models.User{ID: userID, Username: "gettrendsfinal", Email: "gettrendsfinal@test.com"}
	post1 := models.Post{ID: utils.NewID(), AuthorID: userID, Content: "post1", Tags: []string{"#golang"}}
	post2 := models.Post{ID: utils.NewID(), AuthorID: userID, Content: "post2", Tags: []string{"#golang"}}
	db.Create(&user)
	repo.Create(&post1)
	repo.Create(&post2)

	// Get trends
	trends, err := service.GetTrends(10)
	if err != nil {
		t.Fatalf("GetTrends: %v", err)
	}
	if len(trends) < 1 {
		t.Fatalf("expected at least 1 trend, got %d", len(trends))
	}
}

// --- Post Service: UpdatePost ---

func TestPostService_UpdatePostFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	// Create test user and post
	userID := utils.NewID()
	postID := utils.NewID()
	user := models.User{ID: userID, Username: "updpostfinal", Email: "updpostfinal@test.com"}
	post := models.Post{ID: postID, AuthorID: userID, Content: "original"}
	db.Create(&user)
	repo.Create(&post)

	// Update post
	updated, err := service.UpdatePost(postID, models.UpdatePostInput{Content: "updated"}, userID)
	if err != nil {
		t.Fatalf("UpdatePost: %v", err)
	}
	if updated.Content != "updated" {
		t.Fatalf("expected 'updated', got %q", updated.Content)
	}
}

// --- Post Service: UpdatePost not found ---

func TestPostService_UpdatePostNotFoundFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	_, err := service.UpdatePost(utils.NewID(), models.UpdatePostInput{Content: "updated"}, utils.NewID())
	if err == nil {
		t.Fatal("expected error for non-existent post")
	}
}

// --- Post Service: UpdatePost forbidden ---

func TestPostService_UpdatePostForbiddenFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	// Create test user and post
	userID := utils.NewID()
	otherUserID := utils.NewID()
	postID := utils.NewID()
	user := models.User{ID: userID, Username: "updpostforbiddenfinal", Email: "updpostforbiddenfinal@test.com"}
	otherUser := models.User{ID: otherUserID, Username: "otherupdpostfinal", Email: "otherupdpostfinal@test.com"}
	post := models.Post{ID: postID, AuthorID: userID, Content: "post"}
	db.Create(&user)
	db.Create(&otherUser)
	repo.Create(&post)

	// Try to update other user's post
	_, err := service.UpdatePost(postID, models.UpdatePostInput{Content: "hacked"}, otherUserID)
	if err == nil {
		t.Fatal("expected error for updating other user's post")
	}
}

// --- Post Service: DeletePost ---

func TestPostService_DeletePostFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	// Create test user and post
	userID := utils.NewID()
	postID := utils.NewID()
	user := models.User{ID: userID, Username: "delpostfinal", Email: "delpostfinal@test.com"}
	post := models.Post{ID: postID, AuthorID: userID, Content: "to be deleted"}
	db.Create(&user)
	repo.Create(&post)

	// Delete post
	err := service.DeletePost(postID, userID)
	if err != nil {
		t.Fatalf("DeletePost: %v", err)
	}

	// Verify deleted
	_, err = repo.GetByID(postID)
	if err == nil {
		t.Fatal("expected error for deleted post")
	}
}

// --- Post Service: DeletePost not found ---

func TestPostService_DeletePostNotFoundFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	err := service.DeletePost(utils.NewID(), utils.NewID())
	if err == nil {
		t.Fatal("expected error for non-existent post")
	}
}

// --- Post Service: DeletePost forbidden ---

func TestPostService_DeletePostForbiddenFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	// Create test user and post
	userID := utils.NewID()
	otherUserID := utils.NewID()
	postID := utils.NewID()
	user := models.User{ID: userID, Username: "delpostforbiddenfinal", Email: "delpostforbiddenfinal@test.com"}
	otherUser := models.User{ID: otherUserID, Username: "otherdelpostfinal", Email: "otherdelpostfinal@test.com"}
	post := models.Post{ID: postID, AuthorID: userID, Content: "post"}
	db.Create(&user)
	db.Create(&otherUser)
	repo.Create(&post)

	// Try to delete other user's post
	err := service.DeletePost(postID, otherUserID)
	if err == nil {
		t.Fatal("expected error for deleting other user's post")
	}
}

// --- User Service: GetUsers ---

func TestUserService_GetUsersFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewUserService(repo)

	// Create test users
	user1ID := utils.NewID()
	user2ID := utils.NewID()
	user1 := models.User{ID: user1ID, Username: "getusers1final", Email: "gu1final@test.com"}
	user2 := models.User{ID: user2ID, Username: "getusers2final", Email: "gu2final@test.com"}
	db.Create(&user1)
	db.Create(&user2)

	// Get users
	users, err := service.GetUsers()
	if err != nil {
		t.Fatalf("GetUsers: %v", err)
	}
	if len(users) < 2 {
		t.Fatalf("expected at least 2 users, got %d", len(users))
	}
}

// --- User Service: GetUserByUsername ---

func TestUserService_GetUserByUsernameFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewUserService(repo)

	// Create test user
	userID := utils.NewID()
	user := models.User{ID: userID, Username: "getuserbynamefinal", Email: "getuserbynamefinal@test.com"}
	db.Create(&user)

	// Get user by username
	got, err := service.GetUserByUsername("getuserbynamefinal")
	if err != nil {
		t.Fatalf("GetUserByUsername: %v", err)
	}
	if got.ID != userID {
		t.Fatalf("expected ID %s, got %s", userID, got.ID)
	}
}

// --- User Service: GetUserByUsername not found ---

func TestUserService_GetUserByUsernameNotFoundFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewUserService(repo)

	_, err := service.GetUserByUsername("nonexistent")
	if err == nil {
		t.Fatal("expected error for non-existent username")
	}
}

// --- User Service: GetUser ---

func TestUserService_GetUserFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewUserService(repo)

	// Create test user
	userID := utils.NewID()
	user := models.User{ID: userID, Username: "getuserfinal", Email: "getuserfinal@test.com"}
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

func TestUserService_GetUserNotFoundFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewUserService(repo)

	_, err := service.GetUser(utils.NewID())
	if err == nil {
		t.Fatal("expected error for non-existent user")
	}
}

// --- User Service: UpdateUser ---

func TestUserService_UpdateUserFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewUserService(repo)

	// Create test user
	userID := utils.NewID()
	user := models.User{ID: userID, Username: "upduserfinal", Email: "upduserfinal@test.com"}
	db.Create(&user)

	// Update user
	updated, err := service.UpdateUser(userID, models.UpdateUserInput{Bio: "new bio"})
	if err != nil {
		t.Fatalf("UpdateUser: %v", err)
	}
	if updated.Bio != "new bio" {
		t.Fatalf("expected bio 'new bio', got %q", updated.Bio)
	}
}

// --- User Service: UpdateUser not found ---

func TestUserService_UpdateUserNotFoundFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewUserService(repo)

	_, err := service.UpdateUser(utils.NewID(), models.UpdateUserInput{Bio: "bio"})
	if err == nil {
		t.Fatal("expected error for non-existent user")
	}
}

// --- User Service: DeleteUser ---

func TestUserService_DeleteUserFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewUserService(repo)

	// Create test user
	userID := utils.NewID()
	user := models.User{ID: userID, Username: "deluserfinal", Email: "deluserfinal@test.com"}
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

// --- User Service: DeleteUser not found ---

func TestUserService_DeleteUserNotFoundFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewUserService(repo)

	err := service.DeleteUser(utils.NewID())
	if err == nil {
		t.Fatal("expected error for non-existent user")
	}
}

// --- Auth Service: GetUserByID ---

func TestAuthService_GetUserByIDFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewAuthService(repo)

	// Create test user
	userID := utils.NewID()
	user := models.User{ID: userID, Username: "getbyidfinal", Email: "getbyidfinal@test.com"}
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

func TestAuthService_GetUserByIDNotFoundFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewAuthService(repo)

	_, err := service.GetUserByID(utils.NewID())
	if err == nil {
		t.Fatal("expected error for non-existent user")
	}
}

// --- Notification Service: GetUnread ---

func TestNotificationService_GetUnreadFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewNotificationRepositories(db)
	pubsub := repositories.NewNotificationPubSub(sharedRDB)
	service := services.NewNotificationService(repo, pubsub)

	// Create test user
	userID := utils.NewID()
	user := models.User{ID: userID, Username: "notifuserfinal", Email: "notifuserfinal@test.com"}
	db.Create(&user)

	// Send notification
	service.SendNotification(userID, "", utils.NewID(), "actor", "test", "content", "")

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

func TestNotificationService_MarkAllReadFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewNotificationRepositories(db)
	pubsub := repositories.NewNotificationPubSub(sharedRDB)
	service := services.NewNotificationService(repo, pubsub)

	// Create test user
	userID := utils.NewID()
	user := models.User{ID: userID, Username: "markallreadfinal", Email: "markallreadfinal@test.com"}
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

func TestNotificationService_MarkReadFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewNotificationRepositories(db)
	pubsub := repositories.NewNotificationPubSub(sharedRDB)
	service := services.NewNotificationService(repo, pubsub)

	// Create test user
	userID := utils.NewID()
	user := models.User{ID: userID, Username: "markreadfinal", Email: "markreadfinal@test.com"}
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

// --- Friend Service: SendRequest ---

func TestFriendService_SendRequestFinal(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	// Create test users
	senderID := utils.NewID()
	receiverID := utils.NewID()
	sender := models.User{ID: senderID, Username: "senderfinal", Email: "senderfinal@test.com"}
	receiver := models.User{ID: receiverID, Username: "receiverfinal", Email: "receiverfinal@test.com"}
	db.Create(&sender)
	db.Create(&receiver)

	// Send request
	err := service.SendRequest(senderID, receiverID)
	if err != nil {
		t.Fatalf("SendRequest: %v", err)
	}
}

// --- Friend Service: SendRequest to self ---

func TestFriendService_SendRequestToSelfFinal(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	// Create test user
	userID := utils.NewID()
	user := models.User{ID: userID, Username: "selffinal", Email: "selffinal@test.com"}
	db.Create(&user)

	// Send request to self
	err := service.SendRequest(userID, userID)
	if err == nil {
		t.Fatal("expected error for sending request to self")
	}
}

// --- Friend Service: Follow ---

func TestFriendService_FollowFinal(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	// Create test users
	followerID := utils.NewID()
	targetID := utils.NewID()
	follower := models.User{ID: followerID, Username: "followerfinal", Email: "followerfinal@test.com"}
	target := models.User{ID: targetID, Username: "targetfinal", Email: "targetfinal@test.com"}
	db.Create(&follower)
	db.Create(&target)

	// Follow
	err := service.Follow(followerID, targetID)
	if err != nil {
		t.Fatalf("Follow: %v", err)
	}
}

// --- Friend Service: Follow self ---

func TestFriendService_FollowSelfFinal(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	// Create test user
	userID := utils.NewID()
	user := models.User{ID: userID, Username: "selfollowfinal", Email: "selfollowfinal@test.com"}
	db.Create(&user)

	// Follow self
	err := service.Follow(userID, userID)
	if err == nil {
		t.Fatal("expected error for following self")
	}
}

// --- Friend Service: Unfollow ---

func TestFriendService_UnfollowFinal(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	// Create test users
	followerID := utils.NewID()
	targetID := utils.NewID()
	follower := models.User{ID: followerID, Username: "unfollowerfinal", Email: "unfollowerfinal@test.com"}
	target := models.User{ID: targetID, Username: "untargetfinal", Email: "untargetfinal@test.com"}
	db.Create(&follower)
	db.Create(&target)

	// Follow then unfollow
	service.Follow(followerID, targetID)
	err := service.Unfollow(followerID, targetID)
	if err != nil {
		t.Fatalf("Unfollow: %v", err)
	}
}

// --- Friend Service: Unfollow not following ---

func TestFriendService_UnfollowNotFollowingFinal(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	// Create test users
	userID := utils.NewID()
	targetID := utils.NewID()
	user := models.User{ID: userID, Username: "notfollowingfinal", Email: "notfollowingfinal@test.com"}
	target := models.User{ID: targetID, Username: "nftargetfinal", Email: "nftargetfinal@test.com"}
	db.Create(&user)
	db.Create(&target)

	// Unfollow without following
	err := service.Unfollow(userID, targetID)
	if err == nil {
		t.Fatal("expected error for unfollowing without following")
	}
}

// --- Friend Service: RemoveFriend ---

func TestFriendService_RemoveFriendFinal(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	// Create test users
	user1ID := utils.NewID()
	user2ID := utils.NewID()
	user1 := models.User{ID: user1ID, Username: "rmfriend1final", Email: "rmfriend1final@test.com"}
	user2 := models.User{ID: user2ID, Username: "rmfriend2final", Email: "rmfriend2final@test.com"}
	db.Create(&user1)
	db.Create(&user2)

	// Make friends then remove
	service.SendRequest(user1ID, user2ID)
	service.AcceptRequest(user2ID, user1ID)
	err := service.RemoveFriend(user1ID, user2ID)
	if err != nil {
		t.Fatalf("RemoveFriend: %v", err)
	}
}

// --- Friend Service: RemoveFriend not friends ---

func TestFriendService_RemoveFriendNotFriendsFinal(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	// Create test users
	user1ID := utils.NewID()
	user2ID := utils.NewID()
	user1 := models.User{ID: user1ID, Username: "notfriends1final", Email: "notfriends1final@test.com"}
	user2 := models.User{ID: user2ID, Username: "notfriends2final", Email: "notfriends2final@test.com"}
	db.Create(&user1)
	db.Create(&user2)

	// Remove without being friends
	err := service.RemoveFriend(user1ID, user2ID)
	if err == nil {
		t.Fatal("expected error for removing non-friend")
	}
}

// --- Friend Service: CountFollowers ---

func TestFriendService_CountFollowersFinal(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	// Create test users
	user1ID := utils.NewID()
	user2ID := utils.NewID()
	targetID := utils.NewID()
	user1 := models.User{ID: user1ID, Username: "countfollower1final", Email: "cf1final@test.com"}
	user2 := models.User{ID: user2ID, Username: "countfollower2final", Email: "cf2final@test.com"}
	target := models.User{ID: targetID, Username: "counttargetfinal", Email: "ctfinal@test.com"}
	db.Create(&user1)
	db.Create(&user2)
	db.Create(&target)

	// Create follows
	service.Follow(user1ID, targetID)
	service.Follow(user2ID, targetID)

	// Count followers
	count, err := service.CountFollowers(targetID)
	if err != nil {
		t.Fatalf("CountFollowers: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected 2 followers, got %d", count)
	}
}

// --- Friend Service: CountFollowing ---

func TestFriendService_CountFollowingFinal(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	// Create test users
	userID := utils.NewID()
	user1ID := utils.NewID()
	user2ID := utils.NewID()
	user := models.User{ID: userID, Username: "countfollowingfinal", Email: "cffinal@test.com"}
	user1 := models.User{ID: user1ID, Username: "countfollowing1final", Email: "cf1final@test.com"}
	user2 := models.User{ID: user2ID, Username: "countfollowing2final", Email: "cf2final@test.com"}
	db.Create(&user)
	db.Create(&user1)
	db.Create(&user2)

	// Create follows
	service.Follow(userID, user1ID)
	service.Follow(userID, user2ID)

	// Count following
	count, err := service.CountFollowing(userID)
	if err != nil {
		t.Fatalf("CountFollowing: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected 2 following, got %d", count)
	}
}
