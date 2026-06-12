package test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"ft_transcendence/backend/internal/config"
	"ft_transcendence/backend/internal/routes"
)

func brokenDB(t *testing.T) *gorm.DB {
	t.Helper()
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("config.Load: %v", err)
	}
	db, err := cfg.Postgres.Connect()
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("db handle: %v", err)
	}
	_ = sqlDB.Close()
	return db
}

func brokenControllers(t *testing.T) *routes.Controllers {
	t.Helper()
	SetupTestEnv()
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("config.Load: %v", err)
	}
	return routes.Wire(brokenDB(t), sharedRDB, cfg)
}

func testCtx(method, target, body string) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	var r *http.Request
	if body != "" {
		r = httptest.NewRequest(method, target, strings.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
	} else {
		r = httptest.NewRequest(method, target, nil)
	}
	c.Request = r
	return c, w
}

func setParam(c *gin.Context, key, value string) {
	c.Params = append(c.Params, gin.Param{Key: key, Value: value})
}

func TestNotificationController_Unauthorized(t *testing.T) {
	ctrl := brokenControllers(t)

	c, w := testCtx(http.MethodGet, "/", "")
	ctrl.Notification.GetNotifications(c)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("GetNotifications: expected 401, got %d", w.Code)
	}

	c, w = testCtx(http.MethodPatch, "/", "")
	ctrl.Notification.MarkAllRead(c)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("MarkAllRead: expected 401, got %d", w.Code)
	}

	c, w = testCtx(http.MethodPatch, "/", "")
	ctrl.Notification.MarkRead(c)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("MarkRead: expected 401, got %d", w.Code)
	}
}

func TestNotificationController_DBErrors(t *testing.T) {
	ctrl := brokenControllers(t)

	c, w := testCtx(http.MethodGet, "/", "")
	c.Set("user_id", "u1")
	ctrl.Notification.GetNotifications(c)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("GetNotifications: expected 500, got %d", w.Code)
	}

	c, w = testCtx(http.MethodPatch, "/", "")
	c.Set("user_id", "u1")
	ctrl.Notification.MarkAllRead(c)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("MarkAllRead: expected 500, got %d", w.Code)
	}

	c, w = testCtx(http.MethodPatch, "/", "")
	c.Set("user_id", "u1")
	setParam(c, "id", "n1")
	ctrl.Notification.MarkRead(c)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("MarkRead: expected 500, got %d", w.Code)
	}
}

func TestGDPRController_Unauthorized(t *testing.T) {
	ctrl := brokenControllers(t)

	c, w := testCtx(http.MethodGet, "/", "")
	ctrl.GDPR.ExportUserData(c)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("ExportUserData: expected 401, got %d", w.Code)
	}

	c, w = testCtx(http.MethodDelete, "/", "")
	ctrl.GDPR.DeleteUserData(c)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("DeleteUserData: expected 401, got %d", w.Code)
	}
}

func TestGDPRController_DBErrors(t *testing.T) {
	ctrl := brokenControllers(t)

	c, w := testCtx(http.MethodGet, "/", "")
	c.Set("user_id", "u1")
	ctrl.GDPR.ExportUserData(c)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("ExportUserData: expected 500, got %d", w.Code)
	}

	c, w = testCtx(http.MethodDelete, "/", "")
	c.Set("user_id", "u1")
	ctrl.GDPR.DeleteUserData(c)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("DeleteUserData: expected 500, got %d", w.Code)
	}
}

func TestGamificationController_Branches(t *testing.T) {
	ctrl := brokenControllers(t)

	c, w := testCtx(http.MethodGet, "/", "")
	setParam(c, "id", "")
	ctrl.Gamification.Get(c)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("Get empty id: expected 400, got %d", w.Code)
	}

	c, w = testCtx(http.MethodGet, "/", "")
	setParam(c, "id", "u1")
	ctrl.Gamification.Get(c)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("Get db error: expected 500, got %d", w.Code)
	}

	c, w = testCtx(http.MethodGet, "/", "")
	ctrl.Gamification.GetLeaderboard(c)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("GetLeaderboard db error: expected 500, got %d", w.Code)
	}
}

func TestFriendController_Unauthorized(t *testing.T) {
	ctrl := brokenControllers(t)

	for _, h := range []func(*gin.Context){
		ctrl.Friend.SendFriendRequest,
		ctrl.Friend.AcceptFriend,
		ctrl.Friend.FollowUser,
		ctrl.Friend.UnfollowUser,
		ctrl.Friend.RemoveFriend,
		ctrl.Friend.RejectFriendRequest,
	} {
		c, w := testCtx(http.MethodPost, "/", "")
		setParam(c, "id", "target")
		h(c)
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401 without user_id, got %d", w.Code)
		}
	}

	for _, h := range []func(*gin.Context){
		ctrl.Friend.SendFriendRequest,
		ctrl.Friend.AcceptFriend,
	} {
		c, w := testCtx(http.MethodPost, "/", "")
		c.Set("user_id", "u1")
		setParam(c, "id", "target")
		h(c)
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401 without username, got %d", w.Code)
		}
	}
}

func TestFriendController_DBErrors(t *testing.T) {
	ctrl := brokenControllers(t)

	for _, h := range []func(*gin.Context){
		ctrl.Friend.GetFollowers,
		ctrl.Friend.GetFollowing,
		ctrl.Friend.GetFriends,
	} {
		c, w := testCtx(http.MethodGet, "/", "")
		setParam(c, "id", "u1")
		h(c)
		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500 on db error, got %d", w.Code)
		}
	}
}

func TestUserController_DBErrors(t *testing.T) {
	ctrl := brokenControllers(t)

	c, w := testCtx(http.MethodGet, "/?username=someone", "")
	ctrl.User.GetUsers(c)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("GetUsers by username: expected 500, got %d", w.Code)
	}

	c, w = testCtx(http.MethodGet, "/", "")
	ctrl.User.GetUsers(c)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("GetUsers list: expected 500, got %d", w.Code)
	}

	c, w = testCtx(http.MethodGet, "/", "")
	setParam(c, "id", "u1")
	ctrl.User.GetUser(c)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("GetUser: expected 500, got %d", w.Code)
	}
}

func TestUserController_Unauthorized(t *testing.T) {
	ctrl := brokenControllers(t)

	c, w := testCtx(http.MethodPut, "/", `{"display_name":"x"}`)
	setParam(c, "id", "u1")
	ctrl.User.UpdateUser(c)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("UpdateUser: expected 401, got %d", w.Code)
	}

	c, w = testCtx(http.MethodDelete, "/", `{"password":"x"}`)
	setParam(c, "id", "u1")
	ctrl.User.DeleteUser(c)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("DeleteUser: expected 401, got %d", w.Code)
	}
}

func TestUserController_Forbidden(t *testing.T) {
	ctrl := brokenControllers(t)

	c, w := testCtx(http.MethodPut, "/", `{"display_name":"x"}`)
	c.Set("user_id", "other")
	setParam(c, "id", "u1")
	ctrl.User.UpdateUser(c)
	if w.Code != http.StatusForbidden {
		t.Fatalf("UpdateUser forbidden: expected 403, got %d", w.Code)
	}

	c, w = testCtx(http.MethodDelete, "/", `{"password":"x"}`)
	c.Set("user_id", "other")
	setParam(c, "id", "u1")
	ctrl.User.DeleteUser(c)
	if w.Code != http.StatusForbidden {
		t.Fatalf("DeleteUser forbidden: expected 403, got %d", w.Code)
	}
}

func TestPostController_TrendsClamp(t *testing.T) {
	router, _ := SetupTestEnv()

	req := httptest.NewRequest(http.MethodGet, "/api/trends?limit=-5", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("trends limit<=0: expected 200, got %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/trends?limit=999", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("trends limit>50: expected 200, got %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/posts?limit=0&offset=-3", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("posts limit=0: expected 200, got %d", w.Code)
	}
}

func TestPostController_Unauthorized(t *testing.T) {
	ctrl := brokenControllers(t)

	c, w := testCtx(http.MethodPost, "/", `{"content":"hello"}`)
	ctrl.Post.CreatePost(c)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("CreatePost: expected 401, got %d", w.Code)
	}

	c, w = testCtx(http.MethodPut, "/", `{"content":"hi"}`)
	setParam(c, "id", "p1")
	ctrl.Post.UpdatePost(c)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("UpdatePost: expected 401, got %d", w.Code)
	}

	c, w = testCtx(http.MethodDelete, "/", "")
	setParam(c, "id", "p1")
	ctrl.Post.DeletePost(c)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("DeletePost: expected 401, got %d", w.Code)
	}

	c, w = testCtx(http.MethodPost, "/", `{"value":1}`)
	setParam(c, "id", "p1")
	ctrl.Post.React(c)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("React: expected 401, got %d", w.Code)
	}

	c, w = testCtx(http.MethodPost, "/", `{"value":1}`)
	setParam(c, "commentId", "c1")
	ctrl.Post.ReactComment(c)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("ReactComment: expected 401, got %d", w.Code)
	}

	c, w = testCtx(http.MethodPut, "/", `{"content":"x"}`)
	setParam(c, "commentId", "c1")
	ctrl.Post.UpdateComment(c)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("UpdateComment: expected 401, got %d", w.Code)
	}

	c, w = testCtx(http.MethodDelete, "/", "")
	setParam(c, "commentId", "c1")
	ctrl.Post.DeleteComment(c)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("DeleteComment: expected 401, got %d", w.Code)
	}
}

func TestPostController_CreateCommentValidation(t *testing.T) {
	ctrl := brokenControllers(t)

	c, w := testCtx(http.MethodPost, "/", "")
	setParam(c, "id", "p1")
	ctrl.Post.CreateComment(c)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("CreateComment empty: expected 400, got %d", w.Code)
	}

	form := urlValues("content", strings.Repeat("a", 300))
	c, w = testCtx(http.MethodPost, "/", "")
	c.Request = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(form))
	c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	setParam(c, "id", "p1")
	ctrl.Post.CreateComment(c)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("CreateComment too long: expected 400, got %d", w.Code)
	}

	form = urlValues("content", "nice")
	c, w = testCtx(http.MethodPost, "/", "")
	c.Request = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(form))
	c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	setParam(c, "id", "p1")
	ctrl.Post.CreateComment(c)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("CreateComment no user: expected 401, got %d", w.Code)
	}
}

func TestPostController_DBErrors(t *testing.T) {
	ctrl := brokenControllers(t)

	c, w := testCtx(http.MethodGet, "/", "")
	setParam(c, "id", "p1")
	ctrl.Post.GetPost(c)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("GetPost: expected 500, got %d", w.Code)
	}

	c, w = testCtx(http.MethodGet, "/", "")
	setParam(c, "userId", "u1")
	ctrl.Post.GetPostsByUser(c)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("GetPostsByUser: expected 500, got %d", w.Code)
	}

	c, w = testCtx(http.MethodPost, "/", `{"value":1}`)
	c.Set("user_id", "u1")
	setParam(c, "id", "p1")
	ctrl.Post.React(c)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("React: expected 500, got %d", w.Code)
	}

	c, w = testCtx(http.MethodPost, "/", `{"value":1}`)
	c.Set("user_id", "u1")
	setParam(c, "commentId", "c1")
	ctrl.Post.ReactComment(c)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("ReactComment: expected 500, got %d", w.Code)
	}
}

func urlValues(key, value string) string {
	return key + "=" + strings.ReplaceAll(value, " ", "+")
}
