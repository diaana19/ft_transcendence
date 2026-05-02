package tests

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"

	"github.com/Transcendence/config"
	"github.com/Transcendence/models"
	"github.com/Transcendence/routes"
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	const (
		dbName = "app_test"
		dbUser = "test"
		dbPass = "test"
	)

	container, err := tcpostgres.Run(ctx, "postgres:15",
		tcpostgres.WithDatabase(dbName),
		tcpostgres.WithUsername(dbUser),
		tcpostgres.WithPassword(dbPass),
		tcpostgres.WithInitScripts("../../database/init.sql"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		log.Fatalf("failed to start postgres container: %v", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		log.Fatalf("failed to get container host: %v", err)
	}
	port, err := container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		log.Fatalf("failed to get mapped port: %v", err)
	}

	os.Setenv("DB_HOST", host)
	os.Setenv("DB_PORT", port.Port())
	os.Setenv("DB_USER", dbUser)
	os.Setenv("DB_PASSWORD", dbPass)
	os.Setenv("DB_NAME", dbName)

	code := m.Run()

	if err := testcontainers.TerminateContainer(container); err != nil {
		log.Printf("failed to terminate container: %v", err)
	}

	os.Exit(code)
}

func SetupTestEnv() (*gin.Engine, *gorm.DB) {
	gin.SetMode(gin.TestMode)

	db, err := config.ConnectDB()
	if err != nil {
		panic(fmt.Errorf("connect test db: %w", err))
	}

	db.Exec("DROP TABLE IF EXISTS messages CASCADE")
	db.Exec("DROP TABLE IF EXISTS users CASCADE")
	db.AutoMigrate(&models.User{}, &models.Message{})

	router := gin.Default()
	routes.SetupRoutes(router, db)

	return router, db
}
