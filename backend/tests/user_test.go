package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Transcendence/config"
	"github.com/Transcendence/controllers"
	"github.com/Transcendence/repositories"
	"github.com/Transcendence/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load("../../infra/.env")
}

func setupRouter(t *testing.T) *gin.Engine {
	db, err := config.ConnectDB()
	if err != nil {
		t.Fatalf("could not connect to db: %v", err)
	}

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// mirror your real SetupRoutes
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	authController := controllers.NewAuthController(authService)
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	api := r.Group("/api")
	{
		api.POST("/auth/register", authController.RegisterUser)
		api.GET("/users", userController.GetUsers)
		api.GET("/users/:id", userController.GetUser)
		api.PUT("/users/:id", userController.UpdateUser)
		api.DELETE("/users/:id", userController.DeleteUser)
	}

	return r
}

func TestRegisterUser(t *testing.T) {
	router := setupRouter(t)

	body := map[string]interface{}{
		"username":    "testuser",
		"email":       "test@test.com",
		"password":    "lestd3wasd@DF8",
		"dateOfBirth": "1995-06-15",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestGetUsers(t *testing.T) {
	router := setupRouter(t)

	req, _ := http.NewRequest("GET", "/api/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestGetUser(t *testing.T) {
	router := setupRouter(t)

	// create first
	body := map[string]interface{}{
		"username":    "getuser",
		"email":       "getuser@test.com",
		"password":    "getadw2@userdsadwTef32t3gb",
		"dateOfBirth": "1990-01-01",
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)

	id, ok := result["id"].(string)
	if !ok || id == "" {
		t.Fatalf("register failed, could not get id — body: %s", w.Body.String())
	}

	// then fetch
	req2, _ := http.NewRequest("GET", "/api/users/"+id, nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — body: %s", w2.Code, w2.Body.String())
	}
}

func TestDeleteUser(t *testing.T) {
	router := setupRouter(t)

	// create first
	body := map[string]interface{}{
		"username":    "deleteuser",
		"email":       "delete@test.com",
		"password":    "secreTtq123",
		"dateOfBirth": "1995-06-15",
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	id, ok := result["id"].(string)
	if !ok || id == "" {
		t.Fatalf("register failed — body: %s", w.Body.String())
	}
	// then delete
	req2, _ := http.NewRequest("DELETE", "/api/users/"+id, nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — body: %s", w2.Code, w2.Body.String())
	}
}
