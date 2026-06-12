package test

import (
	"testing"

	"ft_transcendence/backend/internal/models"
	"ft_transcendence/backend/internal/repositories"
	"ft_transcendence/backend/internal/services"
	"ft_transcendence/backend/internal/utils"
)

func TestPostService_GetPostsFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	userID := utils.NewID()
	user := models.User{ID: userID, Username: "getpostsfinal", Email: "getpostsfinal@test.com"}
	post1 := models.Post{ID: utils.NewID(), AuthorID: userID, Content: "post1"}
	post2 := models.Post{ID: utils.NewID(), AuthorID: userID, Content: "post2"}
	db.Create(&user)
	repo.Create(&post1)
	repo.Create(&post2)

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

func TestPostService_GetPostFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	userID := utils.NewID()
	postID := utils.NewID()
	user := models.User{ID: userID, Username: "getpostfinal", Email: "getpostfinal@test.com"}
	post := models.Post{ID: postID, AuthorID: userID, Content: "post content"}
	db.Create(&user)
	repo.Create(&post)

	got, err := service.GetPost(postID)
	if err != nil {
		t.Fatalf("GetPost: %v", err)
	}
	if got.Content != "post content" {
		t.Fatalf("expected 'post content', got %q", got.Content)
	}
}

func TestPostService_GetPostNotFoundFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	_, err := service.GetPost(utils.NewID())
	if err == nil {
		t.Fatal("expected error for non-existent post")
	}
}

func TestPostService_GetPostsByAuthorFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	userID := utils.NewID()
	user := models.User{ID: userID, Username: "getpostsbyauthorfinal", Email: "getpostsbyauthorfinal@test.com"}
	post1 := models.Post{ID: utils.NewID(), AuthorID: userID, Content: "post1"}
	post2 := models.Post{ID: utils.NewID(), AuthorID: userID, Content: "post2"}
	db.Create(&user)
	repo.Create(&post1)
	repo.Create(&post2)

	posts, err := service.GetPostsByAuthor(userID)
	if err != nil {
		t.Fatalf("GetPostsByAuthor: %v", err)
	}
	if len(posts) != 2 {
		t.Fatalf("expected 2 posts, got %d", len(posts))
	}
}

func TestPostService_GetTrendsFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	userID := utils.NewID()
	user := models.User{ID: userID, Username: "gettrendsfinal", Email: "gettrendsfinal@test.com"}
	post1 := models.Post{ID: utils.NewID(), AuthorID: userID, Content: "post1", Tags: []string{"#golang"}}
	post2 := models.Post{ID: utils.NewID(), AuthorID: userID, Content: "post2", Tags: []string{"#golang"}}
	db.Create(&user)
	repo.Create(&post1)
	repo.Create(&post2)

	trends, err := service.GetTrends(10)
	if err != nil {
		t.Fatalf("GetTrends: %v", err)
	}
	if len(trends) < 1 {
		t.Fatalf("expected at least 1 trend, got %d", len(trends))
	}
}

func TestPostService_UpdatePostFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	userID := utils.NewID()
	postID := utils.NewID()
	user := models.User{ID: userID, Username: "updpostfinal", Email: "updpostfinal@test.com"}
	post := models.Post{ID: postID, AuthorID: userID, Content: "original"}
	db.Create(&user)
	repo.Create(&post)

	updated, err := service.UpdatePost(postID, models.UpdatePostInput{Content: "updated"}, userID)
	if err != nil {
		t.Fatalf("UpdatePost: %v", err)
	}
	if updated.Content != "updated" {
		t.Fatalf("expected 'updated', got %q", updated.Content)
	}
}

func TestPostService_UpdatePostNotFoundFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	_, err := service.UpdatePost(utils.NewID(), models.UpdatePostInput{Content: "updated"}, utils.NewID())
	if err == nil {
		t.Fatal("expected error for non-existent post")
	}
}

func TestPostService_UpdatePostForbiddenFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	userID := utils.NewID()
	otherUserID := utils.NewID()
	postID := utils.NewID()
	user := models.User{ID: userID, Username: "updpostforbiddenfinal", Email: "updpostforbiddenfinal@test.com"}
	otherUser := models.User{ID: otherUserID, Username: "otherupdpostfinal", Email: "otherupdpostfinal@test.com"}
	post := models.Post{ID: postID, AuthorID: userID, Content: "post"}
	db.Create(&user)
	db.Create(&otherUser)
	repo.Create(&post)

	_, err := service.UpdatePost(postID, models.UpdatePostInput{Content: "hacked"}, otherUserID)
	if err == nil {
		t.Fatal("expected error for updating other user's post")
	}
}

func TestPostService_DeletePostFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	userID := utils.NewID()
	postID := utils.NewID()
	user := models.User{ID: userID, Username: "delpostfinal", Email: "delpostfinal@test.com"}
	post := models.Post{ID: postID, AuthorID: userID, Content: "to be deleted"}
	db.Create(&user)
	repo.Create(&post)

	err := service.DeletePost(postID, userID)
	if err != nil {
		t.Fatalf("DeletePost: %v", err)
	}

	_, err = repo.GetByID(postID)
	if err == nil {
		t.Fatal("expected error for deleted post")
	}
}

func TestPostService_DeletePostNotFoundFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	err := service.DeletePost(utils.NewID(), utils.NewID())
	if err == nil {
		t.Fatal("expected error for non-existent post")
	}
}

func TestPostService_DeletePostForbiddenFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	userID := utils.NewID()
	otherUserID := utils.NewID()
	postID := utils.NewID()
	user := models.User{ID: userID, Username: "delpostforbiddenfinal", Email: "delpostforbiddenfinal@test.com"}
	otherUser := models.User{ID: otherUserID, Username: "otherdelpostfinal", Email: "otherdelpostfinal@test.com"}
	post := models.Post{ID: postID, AuthorID: userID, Content: "post"}
	db.Create(&user)
	db.Create(&otherUser)
	repo.Create(&post)

	err := service.DeletePost(postID, otherUserID)
	if err == nil {
		t.Fatal("expected error for deleting other user's post")
	}
}

func TestUserService_GetUsersFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewUserService(repo)

	user1ID := utils.NewID()
	user2ID := utils.NewID()
	user1 := models.User{ID: user1ID, Username: "getusers1final", Email: "gu1final@test.com"}
	user2 := models.User{ID: user2ID, Username: "getusers2final", Email: "gu2final@test.com"}
	db.Create(&user1)
	db.Create(&user2)

	users, err := service.GetUsers()
	if err != nil {
		t.Fatalf("GetUsers: %v", err)
	}
	if len(users) < 2 {
		t.Fatalf("expected at least 2 users, got %d", len(users))
	}
}

func TestUserService_GetUserByUsernameFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewUserService(repo)

	userID := utils.NewID()
	user := models.User{ID: userID, Username: "getuserbynamefinal", Email: "getuserbynamefinal@test.com"}
	db.Create(&user)

	got, err := service.GetUserByUsername("getuserbynamefinal")
	if err != nil {
		t.Fatalf("GetUserByUsername: %v", err)
	}
	if got.ID != userID {
		t.Fatalf("expected ID %s, got %s", userID, got.ID)
	}
}

func TestUserService_GetUserByUsernameNotFoundFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewUserService(repo)

	_, err := service.GetUserByUsername("nonexistent")
	if err == nil {
		t.Fatal("expected error for non-existent username")
	}
}

func TestUserService_GetUserFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewUserService(repo)

	userID := utils.NewID()
	user := models.User{ID: userID, Username: "getuserfinal", Email: "getuserfinal@test.com"}
	db.Create(&user)

	got, err := service.GetUser(userID)
	if err != nil {
		t.Fatalf("GetUser: %v", err)
	}
	if got.ID != userID {
		t.Fatalf("expected ID %s, got %s", userID, got.ID)
	}
}

func TestUserService_GetUserNotFoundFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewUserService(repo)

	_, err := service.GetUser(utils.NewID())
	if err == nil {
		t.Fatal("expected error for non-existent user")
	}
}

func TestUserService_UpdateUserFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewUserService(repo)

	userID := utils.NewID()
	user := models.User{ID: userID, Username: "upduserfinal", Email: "upduserfinal@test.com"}
	db.Create(&user)

	updated, err := service.UpdateUser(userID, models.UpdateUserInput{Bio: "new bio"})
	if err != nil {
		t.Fatalf("UpdateUser: %v", err)
	}
	if updated.Bio != "new bio" {
		t.Fatalf("expected bio 'new bio', got %q", updated.Bio)
	}
}

func TestUserService_UpdateUserNotFoundFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewUserService(repo)

	_, err := service.UpdateUser(utils.NewID(), models.UpdateUserInput{Bio: "bio"})
	if err == nil {
		t.Fatal("expected error for non-existent user")
	}
}

func TestUserService_DeleteUserFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewUserService(repo)

	userID := utils.NewID()
	user := models.User{ID: userID, Username: "deluserfinal", Email: "deluserfinal@test.com"}
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

func TestUserService_DeleteUserNotFoundFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewUserService(repo)

	err := service.DeleteUser(utils.NewID())
	if err == nil {
		t.Fatal("expected error for non-existent user")
	}
}

func TestAuthService_GetUserByIDFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewAuthService(repo)

	userID := utils.NewID()
	user := models.User{ID: userID, Username: "getbyidfinal", Email: "getbyidfinal@test.com"}
	db.Create(&user)

	got, err := service.GetUserByID(userID)
	if err != nil {
		t.Fatalf("GetUserByID: %v", err)
	}
	if got.ID != userID {
		t.Fatalf("expected ID %s, got %s", userID, got.ID)
	}
}

func TestAuthService_GetUserByIDNotFoundFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewAuthService(repo)

	_, err := service.GetUserByID(utils.NewID())
	if err == nil {
		t.Fatal("expected error for non-existent user")
	}
}

func TestNotificationService_GetUnreadFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewNotificationRepositories(db)
	pubsub := repositories.NewNotificationPubSub(sharedRDB)
	service := services.NewNotificationService(repo, pubsub)

	userID := utils.NewID()
	user := models.User{ID: userID, Username: "notifuserfinal", Email: "notifuserfinal@test.com"}
	db.Create(&user)

	service.SendNotification(userID, "", utils.NewID(), "actor", "test", "content", "")

	notifs, err := service.GetUnread(userID)
	if err != nil {
		t.Fatalf("GetUnread: %v", err)
	}
	if len(notifs) != 1 {
		t.Fatalf("expected 1 notification, got %d", len(notifs))
	}
}

func TestNotificationService_MarkAllReadFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewNotificationRepositories(db)
	pubsub := repositories.NewNotificationPubSub(sharedRDB)
	service := services.NewNotificationService(repo, pubsub)

	userID := utils.NewID()
	user := models.User{ID: userID, Username: "markallreadfinal", Email: "markallreadfinal@test.com"}
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

func TestNotificationService_MarkReadFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewNotificationRepositories(db)
	pubsub := repositories.NewNotificationPubSub(sharedRDB)
	service := services.NewNotificationService(repo, pubsub)

	userID := utils.NewID()
	user := models.User{ID: userID, Username: "markreadfinal", Email: "markreadfinal@test.com"}
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

func TestFriendService_SendRequestFinal(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	senderID := utils.NewID()
	receiverID := utils.NewID()
	sender := models.User{ID: senderID, Username: "senderfinal", Email: "senderfinal@test.com"}
	receiver := models.User{ID: receiverID, Username: "receiverfinal", Email: "receiverfinal@test.com"}
	db.Create(&sender)
	db.Create(&receiver)

	err := service.SendRequest(senderID, receiverID)
	if err != nil {
		t.Fatalf("SendRequest: %v", err)
	}
}

func TestFriendService_SendRequestToSelfFinal(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	userID := utils.NewID()
	user := models.User{ID: userID, Username: "selffinal", Email: "selffinal@test.com"}
	db.Create(&user)

	err := service.SendRequest(userID, userID)
	if err != nil {
		t.Fatalf("sending request to self is a no-op: %v", err)
	}
}

func TestFriendService_FollowFinal(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	followerID := utils.NewID()
	targetID := utils.NewID()
	follower := models.User{ID: followerID, Username: "followerfinal", Email: "followerfinal@test.com"}
	target := models.User{ID: targetID, Username: "targetfinal", Email: "targetfinal@test.com"}
	db.Create(&follower)
	db.Create(&target)

	err := service.Follow(followerID, targetID)
	if err != nil {
		t.Fatalf("Follow: %v", err)
	}
}

func TestFriendService_FollowSelfFinal(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	userID := utils.NewID()
	user := models.User{ID: userID, Username: "selfollowfinal", Email: "selfollowfinal@test.com"}
	db.Create(&user)

	err := service.Follow(userID, userID)
	if err != nil {
		t.Fatalf("following self is a no-op: %v", err)
	}
}

func TestFriendService_UnfollowFinal(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	followerID := utils.NewID()
	targetID := utils.NewID()
	follower := models.User{ID: followerID, Username: "unfollowerfinal", Email: "unfollowerfinal@test.com"}
	target := models.User{ID: targetID, Username: "untargetfinal", Email: "untargetfinal@test.com"}
	db.Create(&follower)
	db.Create(&target)

	service.Follow(followerID, targetID)
	err := service.Unfollow(followerID, targetID)
	if err != nil {
		t.Fatalf("Unfollow: %v", err)
	}
}

func TestFriendService_UnfollowNotFollowingFinal(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	userID := utils.NewID()
	targetID := utils.NewID()
	user := models.User{ID: userID, Username: "notfollowingfinal", Email: "notfollowingfinal@test.com"}
	target := models.User{ID: targetID, Username: "nftargetfinal", Email: "nftargetfinal@test.com"}
	db.Create(&user)
	db.Create(&target)

	err := service.Unfollow(userID, targetID)
	if err != nil {
		t.Fatalf("unfollowing without following is a no-op: %v", err)
	}
}

func TestFriendService_RemoveFriendFinal(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	user1ID := utils.NewID()
	user2ID := utils.NewID()
	user1 := models.User{ID: user1ID, Username: "rmfriend1final", Email: "rmfriend1final@test.com"}
	user2 := models.User{ID: user2ID, Username: "rmfriend2final", Email: "rmfriend2final@test.com"}
	db.Create(&user1)
	db.Create(&user2)

	service.SendRequest(user1ID, user2ID)
	service.AcceptRequest(user2ID, user1ID)
	err := service.RemoveFriend(user1ID, user2ID)
	if err != nil {
		t.Fatalf("RemoveFriend: %v", err)
	}
}

func TestFriendService_RemoveFriendNotFriendsFinal(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	user1ID := utils.NewID()
	user2ID := utils.NewID()
	user1 := models.User{ID: user1ID, Username: "notfriends1final", Email: "notfriends1final@test.com"}
	user2 := models.User{ID: user2ID, Username: "notfriends2final", Email: "notfriends2final@test.com"}
	db.Create(&user1)
	db.Create(&user2)

	err := service.RemoveFriend(user1ID, user2ID)
	if err != nil {
		t.Fatalf("removing a non-friend is a no-op: %v", err)
	}
}

func TestFriendService_CountFollowersFinal(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	user1ID := utils.NewID()
	user2ID := utils.NewID()
	targetID := utils.NewID()
	user1 := models.User{ID: user1ID, Username: "countfollower1final", Email: "cf1final@test.com"}
	user2 := models.User{ID: user2ID, Username: "countfollower2final", Email: "cf2final@test.com"}
	target := models.User{ID: targetID, Username: "counttargetfinal", Email: "ctfinal@test.com"}
	db.Create(&user1)
	db.Create(&user2)
	db.Create(&target)

	service.Follow(user1ID, targetID)
	service.Follow(user2ID, targetID)

	count, err := service.CountFollowers(targetID)
	if err != nil {
		t.Fatalf("CountFollowers: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected 2 followers, got %d", count)
	}
}

func TestFriendService_CountFollowingFinal(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	userID := utils.NewID()
	user1ID := utils.NewID()
	user2ID := utils.NewID()
	user := models.User{ID: userID, Username: "countfollowingfinal", Email: "cffinal@test.com"}
	user1 := models.User{ID: user1ID, Username: "countfollowing1final", Email: "cf1final@test.com"}
	user2 := models.User{ID: user2ID, Username: "countfollowing2final", Email: "cf2final@test.com"}
	db.Create(&user)
	db.Create(&user1)
	db.Create(&user2)

	service.Follow(userID, user1ID)
	service.Follow(userID, user2ID)

	count, err := service.CountFollowing(userID)
	if err != nil {
		t.Fatalf("CountFollowing: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected 2 following, got %d", count)
	}
}
