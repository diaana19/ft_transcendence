package test

import (
	"strings"
	"testing"

	"ft_transcendence/backend/internal/models"
	"ft_transcendence/backend/internal/repositories"
	"ft_transcendence/backend/internal/services"
	"ft_transcendence/backend/internal/utils"
)

func TestPostService_CreateCommentFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	userID := utils.NewID()
	postID := utils.NewID()
	user := models.User{ID: userID, Username: "createcommentfinal", Email: "createcommentfinal@test.com"}
	post := models.Post{ID: postID, AuthorID: userID, Content: "post"}
	db.Create(&user)
	repo.Create(&post)

	comment, err := service.CreateComment("new comment", userID, postID, nil)
	if err != nil {
		t.Fatalf("CreateComment: %v", err)
	}
	if comment.Content != "new comment" {
		t.Fatalf("expected 'new comment', got %q", comment.Content)
	}
}

func TestPostService_CreateCommentEmptyFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	_, err := service.CreateComment("", utils.NewID(), utils.NewID(), nil)
	if err == nil {
		t.Fatal("expected error for empty comment")
	}
}

func TestPostService_CreateCommentTooLongFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	var longContent strings.Builder
	for range 300 {
		longContent.WriteString("a")
	}

	_, err := service.CreateComment(longContent.String(), utils.NewID(), utils.NewID(), nil)
	if err == nil {
		t.Fatal("expected error for comment over 280 chars")
	}
}

func TestPostService_CreateCommentNotFoundFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	_, err := service.CreateComment("comment", utils.NewID(), utils.NewID(), nil)
	if err == nil {
		t.Fatal("expected error for non-existent post")
	}
}

func TestPostService_UpdateCommentFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	userID := utils.NewID()
	postID := utils.NewID()
	user := models.User{ID: userID, Username: "updcommentfinal", Email: "updcommentfinal@test.com"}
	post := models.Post{ID: postID, AuthorID: userID, Content: "post"}
	db.Create(&user)
	repo.Create(&post)

	comment, _ := service.CreateComment("original", userID, postID, nil)

	updated, err := service.UpdateComment(comment.ID, models.UpdateCommentInput{Content: "updated"}, userID)
	if err != nil {
		t.Fatalf("UpdateComment: %v", err)
	}
	if updated.Content != "updated" {
		t.Fatalf("expected 'updated', got %q", updated.Content)
	}
}

func TestPostService_UpdateCommentNotFoundFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	_, err := service.UpdateComment(utils.NewID(), models.UpdateCommentInput{Content: "updated"}, utils.NewID())
	if err == nil {
		t.Fatal("expected error for non-existent comment")
	}
}

func TestPostService_UpdateCommentForbiddenFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	userID := utils.NewID()
	otherUserID := utils.NewID()
	postID := utils.NewID()
	user := models.User{ID: userID, Username: "updcforbiddenfinal", Email: "updcforbiddenfinal@test.com"}
	otherUser := models.User{ID: otherUserID, Username: "otherupdcfinal", Email: "otherupdcfinal@test.com"}
	post := models.Post{ID: postID, AuthorID: userID, Content: "post"}
	db.Create(&user)
	db.Create(&otherUser)
	repo.Create(&post)

	comment, _ := service.CreateComment("comment", userID, postID, nil)

	_, err := service.UpdateComment(comment.ID, models.UpdateCommentInput{Content: "hacked"}, otherUserID)
	if err == nil {
		t.Fatal("expected error for updating other user's comment")
	}
}

func TestPostService_DeleteCommentFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	userID := utils.NewID()
	postID := utils.NewID()
	user := models.User{ID: userID, Username: "delcommentfinal", Email: "delcommentfinal@test.com"}
	post := models.Post{ID: postID, AuthorID: userID, Content: "post"}
	db.Create(&user)
	repo.Create(&post)

	comment, _ := service.CreateComment("to delete", userID, postID, nil)

	err := service.DeleteComment(comment.ID, userID)
	if err != nil {
		t.Fatalf("DeleteComment: %v", err)
	}
}

func TestPostService_DeleteCommentNotFoundFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	err := service.DeleteComment(utils.NewID(), utils.NewID())
	if err == nil {
		t.Fatal("expected error for non-existent comment")
	}
}

func TestPostService_DeleteCommentForbiddenFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	userID := utils.NewID()
	otherUserID := utils.NewID()
	postID := utils.NewID()
	user := models.User{ID: userID, Username: "delcforbiddenfinal", Email: "delcforbiddenfinal@test.com"}
	otherUser := models.User{ID: otherUserID, Username: "otherdelcfinal", Email: "otherdelcfinal@test.com"}
	post := models.Post{ID: postID, AuthorID: userID, Content: "post"}
	db.Create(&user)
	db.Create(&otherUser)
	repo.Create(&post)

	comment, _ := service.CreateComment("comment", userID, postID, nil)

	err := service.DeleteComment(comment.ID, otherUserID)
	if err == nil {
		t.Fatal("expected error for deleting other user's comment")
	}
}

func TestPostService_ReactToPostFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	userID := utils.NewID()
	postID := utils.NewID()
	user := models.User{ID: userID, Username: "reactpostfinal", Email: "reactpostfinal@test.com"}
	post := models.Post{ID: postID, AuthorID: utils.NewID(), Content: "post"}
	db.Create(&user)
	repo.Create(&post)

	value, _, err := service.ReactToPost(userID, postID, 1)
	if err != nil {
		t.Fatalf("ReactToPost like: %v", err)
	}
	if value != 1 {
		t.Fatalf("expected 1, got %d", value)
	}

	value, _, err = service.ReactToPost(userID, postID, 1)
	if err != nil {
		t.Fatalf("ReactToPost unlike: %v", err)
	}
	if value != 0 {
		t.Fatalf("expected 0, got %d", value)
	}
}

func TestPostService_ReactToPostNotFoundFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	_, _, err := service.ReactToPost(utils.NewID(), utils.NewID(), 1)
	if err == nil {
		t.Fatal("expected error for non-existent post")
	}
}

func TestPostService_ReactToCommentFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	userID := utils.NewID()
	postID := utils.NewID()
	commentID := utils.NewID()
	user := models.User{ID: userID, Username: "reactcommentfinal", Email: "reactcommentfinal@test.com"}
	post := models.Post{ID: postID, AuthorID: utils.NewID(), Content: "post"}
	comment := models.Reply{ID: commentID, PostID: postID, AuthorID: utils.NewID(), Content: "comment"}
	db.Create(&user)
	repo.Create(&post)
	repo.CreateComment(&comment)

	value, _, err := service.ReactToComment(userID, commentID, 1)
	if err != nil {
		t.Fatalf("ReactToComment: %v", err)
	}
	if value != 1 {
		t.Fatalf("expected 1, got %d", value)
	}
}

func TestPostService_ReactToCommentNotFoundFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	_, _, err := service.ReactToComment(utils.NewID(), utils.NewID(), 1)
	if err == nil {
		t.Fatal("expected error for non-existent comment")
	}
}

func TestPostService_GetPostReactionFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	userID := utils.NewID()
	postID := utils.NewID()
	user := models.User{ID: userID, Username: "getreactpostfinal", Email: "getreactpostfinal@test.com"}
	post := models.Post{ID: postID, AuthorID: utils.NewID(), Content: "post"}
	db.Create(&user)
	repo.Create(&post)

	repo.SetPostReaction(userID, postID, 1)

	val, err := service.GetPostReaction(userID, postID)
	if err != nil {
		t.Fatalf("GetPostReaction: %v", err)
	}
	if val != 1 {
		t.Fatalf("expected 1, got %d", val)
	}
}

func TestPostService_GetCommentReactionFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	userID := utils.NewID()
	postID := utils.NewID()
	commentID := utils.NewID()
	user := models.User{ID: userID, Username: "getreactcommentfinal", Email: "getreactcommentfinal@test.com"}
	post := models.Post{ID: postID, AuthorID: utils.NewID(), Content: "post"}
	comment := models.Reply{ID: commentID, PostID: postID, AuthorID: utils.NewID(), Content: "comment"}
	db.Create(&user)
	repo.Create(&post)
	repo.CreateComment(&comment)

	repo.SetReplyReaction(userID, commentID, 1)

	val, err := service.GetCommentReaction(userID, commentID)
	if err != nil {
		t.Fatalf("GetCommentReaction: %v", err)
	}
	if val != 1 {
		t.Fatalf("expected 1, got %d", val)
	}
}

func TestAuthService_LoginAuthUserServiceFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewAuthService(repo)

	userID := utils.NewID()
	password := "StrongPass123!"
	hash, _ := utils.HashString(password)
	user := models.User{ID: userID, Username: "loginfinal", Email: "loginfinal@test.com", Password: &hash}
	db.Create(&user)

	_, err := service.LoginAuthUserService("loginfinal@test.com", password)
	if err != nil {
		t.Fatalf("LoginAuthUserService with email: %v", err)
	}

	_, err = service.LoginAuthUserService("loginfinal", password)
	if err != nil {
		t.Fatalf("LoginAuthUserService with username: %v", err)
	}
}

func TestAuthService_LoginAuthUserServiceWrongPasswordFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewAuthService(repo)

	userID := utils.NewID()
	password := "StrongPass123!"
	hash, _ := utils.HashString(password)
	user := models.User{ID: userID, Username: "loginwpfinal", Email: "loginwpfinal@test.com", Password: &hash}
	db.Create(&user)

	_, err := service.LoginAuthUserService("loginwpfinal@test.com", "wrongpassword")
	if err == nil {
		t.Fatal("expected error for wrong password")
	}
}

func TestAuthService_LoginAuthUserServiceNotFoundFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewAuthService(repo)

	_, err := service.LoginAuthUserService("nonexistent@test.com", "password")
	if err == nil {
		t.Fatal("expected error for non-existent user")
	}
}

func TestAuthService_CreateAuthUserServiceFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewAuthService(repo)

	password := "StrongPass123!"
	hash, _ := utils.HashString(password)
	user := &models.User{
		Username: "createfinal",
		Email:    "createfinal@test.com",
		Password: &hash,
	}
	created, err := service.CreateAuthUserService(user)
	if err != nil {
		t.Fatalf("CreateAuthUserService: %v", err)
	}
	if created.Username != "createfinal" {
		t.Fatalf("expected username 'createfinal', got %q", created.Username)
	}
}

func TestAuthService_CreateAuthUserServiceDuplicateFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewAuthService(repo)

	password := "StrongPass123!"
	hash, _ := utils.HashString(password)
	user1 := &models.User{
		Username: "dupfinal",
		Email:    "dupfinal@test.com",
		Password: &hash,
	}
	_, err := service.CreateAuthUserService(user1)
	if err != nil {
		t.Fatalf("CreateAuthUserService first: %v", err)
	}

	user2 := &models.User{
		Username: "dupfinal",
		Email:    "dupfinal2@test.com",
		Password: &hash,
	}
	_, err = service.CreateAuthUserService(user2)
	if err == nil {
		t.Fatal("expected error for duplicate username")
	}
}

func TestFriendService_AcceptRequestFinal(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	senderID := utils.NewID()
	receiverID := utils.NewID()
	sender := models.User{ID: senderID, Username: "acceptsenderfinal", Email: "asfinal@test.com"}
	receiver := models.User{ID: receiverID, Username: "acceptreceiverfinal", Email: "arfinal@test.com"}
	db.Create(&sender)
	db.Create(&receiver)

	_, err := service.SendRequest(senderID, receiverID)
	if err != nil {
		t.Fatalf("SendRequest: %v", err)
	}

	err = service.AcceptRequest(receiverID, senderID)
	if err != nil {
		t.Fatalf("AcceptRequest: %v", err)
	}

	areFriends, _ := service.AreFriends(senderID, receiverID)
	if !areFriends {
		t.Fatal("expected friends after accept")
	}
}

func TestFriendService_AcceptRequestWithoutPendingFinal(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	user1ID := utils.NewID()
	user2ID := utils.NewID()
	user1 := models.User{ID: user1ID, Username: "nopendingfinal1", Email: "npfinal1@test.com"}
	user2 := models.User{ID: user2ID, Username: "nopendingfinal2", Email: "npfinal2@test.com"}
	db.Create(&user1)
	db.Create(&user2)

	err := service.AcceptRequest(user1ID, user2ID)
	if err == nil {
		t.Fatal("expected error for accept without pending")
	}
}

func TestFriendService_RejectRequestFinal(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	senderID := utils.NewID()
	receiverID := utils.NewID()
	sender := models.User{ID: senderID, Username: "rejectsenderfinal", Email: "rsfinal@test.com"}
	receiver := models.User{ID: receiverID, Username: "rejectreceiverfinal", Email: "rrfinal@test.com"}
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

func TestFriendService_RejectRequestWithoutPendingFinal(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	user1ID := utils.NewID()
	user2ID := utils.NewID()
	user1 := models.User{ID: user1ID, Username: "norejectfinal1", Email: "nr1final@test.com"}
	user2 := models.User{ID: user2ID, Username: "norejectfinal2", Email: "nr2final@test.com"}
	db.Create(&user1)
	db.Create(&user2)

	err := service.RejectRequest(user1ID, user2ID)
	if err != nil {
		t.Fatalf("reject without pending is a no-op: %v", err)
	}
}

func TestFriendService_GetFollowersFinal(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	user1ID := utils.NewID()
	user2ID := utils.NewID()
	targetID := utils.NewID()
	user1 := models.User{ID: user1ID, Username: "follower1final", Email: "f1final@test.com"}
	user2 := models.User{ID: user2ID, Username: "follower2final", Email: "f2final@test.com"}
	target := models.User{ID: targetID, Username: "targetfinal", Email: "targetfinal@test.com"}
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

func TestFriendService_GetFollowingFinal(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	userID := utils.NewID()
	user1ID := utils.NewID()
	user2ID := utils.NewID()
	user := models.User{ID: userID, Username: "followinguserfinal", Email: "fufinal@test.com"}
	user1 := models.User{ID: user1ID, Username: "following1final", Email: "f1final@test.com"}
	user2 := models.User{ID: user2ID, Username: "following2final", Email: "f2final@test.com"}
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

func TestFriendService_GetFriendsFinal(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	user1ID := utils.NewID()
	user2ID := utils.NewID()
	user1 := models.User{ID: user1ID, Username: "friend1final", Email: "fr1final@test.com"}
	user2 := models.User{ID: user2ID, Username: "friend2final", Email: "fr2final@test.com"}
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

func TestFriendService_AreFriendsFinal(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	user1ID := utils.NewID()
	user2ID := utils.NewID()
	user1 := models.User{ID: user1ID, Username: "arefriends1final", Email: "af1final@test.com"}
	user2 := models.User{ID: user2ID, Username: "arefriends2final", Email: "af2final@test.com"}
	db.Create(&user1)
	db.Create(&user2)

	areFriends, err := service.AreFriends(user1ID, user2ID)
	if err != nil {
		t.Fatalf("AreFriends: %v", err)
	}
	if areFriends {
		t.Fatal("expected not friends initially")
	}

	service.SendRequest(user1ID, user2ID)
	service.AcceptRequest(user2ID, user1ID)

	areFriends, err = service.AreFriends(user1ID, user2ID)
	if err != nil {
		t.Fatalf("AreFriends after accept: %v", err)
	}
	if !areFriends {
		t.Fatal("expected friends after accept")
	}
}
