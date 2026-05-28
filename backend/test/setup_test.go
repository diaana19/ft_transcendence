package test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"

	"ft_transcendence/backend/internal/config"
	"ft_transcendence/backend/internal/models"
	"ft_transcendence/backend/internal/routes"
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	// Ryuk's resource reaper fails to start under rootless podman, so disable it;
	// containers are torn down explicitly via TerminateContainer below.
	os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")

	// Effectively disable rate limiting so the shared test client IP is not throttled.
	os.Setenv("RATE_LIMIT_MAX", "1000000")

	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")

	postgresContainer, err := tcpostgres.Run(ctx, "postgres:15",
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

	postgresHost, err := postgresContainer.Host(ctx)
	if err != nil {
		log.Fatalf("failed to get postgres host: %v", err)
	}
	postgresPort, err := postgresContainer.MappedPort(ctx, "5432/tcp")
	if err != nil {
		log.Fatalf("failed to get postgres mapped port: %v", err)
	}

	os.Setenv("DB_HOST", postgresHost)
	os.Setenv("DB_PORT", postgresPort.Port())

	redisContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "redis:7-alpine",
			ExposedPorts: []string{"6379/tcp"},
			WaitingFor:   wait.ForLog("Ready to accept connections").WithStartupTimeout(60 * time.Second),
		},
		Started: true,
	})
	if err != nil {
		log.Fatalf("failed to start redis container: %v", err)
	}

	redisHost, err := redisContainer.Host(ctx)
	if err != nil {
		log.Fatalf("failed to get redis host: %v", err)
	}
	redisPort, err := redisContainer.MappedPort(ctx, "6379/tcp")
	if err != nil {
		log.Fatalf("failed to get redis mapped port: %v", err)
	}

	os.Setenv("REDIS_HOST", redisHost)
	os.Setenv("REDIS_PORT", redisPort.Port())

	code := m.Run()

	if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
		log.Printf("failed to terminate postgres container: %v", err)
	}
	if err := testcontainers.TerminateContainer(redisContainer); err != nil {
		log.Printf("failed to terminate redis container: %v", err)
	}

	os.Exit(code)
}

// sharedDB and sharedRDB are initialised once and reused across every test.
// Tests run serially, so a single connection pool avoids exhausting Postgres'
// connection limit (each cfg.Postgres.Connect call would otherwise open a pool
// that is never closed).
var (
	sharedDB  *gorm.DB
	sharedRDB *goredis.Client
)

func SetupTestEnv() (*gin.Engine, *gorm.DB) {
	gin.SetMode(gin.TestMode)

	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	if sharedDB == nil {
		db, err := cfg.Postgres.Connect()
		if err != nil {
			panic(fmt.Errorf("connect test postgres: %w", err))
		}
		sharedDB = db
	}
	db := sharedDB

	if sharedRDB == nil {
		rdb, err := cfg.Redis.Connect()
		if err != nil {
			panic(err)
		}
		sharedRDB = rdb
	}
	rdb := sharedRDB

	rdb.FlushDB(context.Background())

	db.Exec("DROP TABLE IF EXISTS messages CASCADE")
	db.Exec("DROP TABLE IF EXISTS users CASCADE")
	db.AutoMigrate(&models.User{}, &models.Message{})
	db.Exec("TRUNCATE TABLE posts, likes, replies, reposts, notifications, friends RESTART IDENTITY CASCADE")

	router := gin.Default()
	ctrl := routes.Wire(db, rdb, cfg)
	routes.SetupRoutes(router, ctrl, rdb)

	return router, db
}
