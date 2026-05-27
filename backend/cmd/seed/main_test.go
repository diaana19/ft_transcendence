package main

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"

	"ft_transcendence/backend/internal/config"
	"ft_transcendence/backend/internal/models"
	"ft_transcendence/backend/internal/utils"
)

// TestMain starts a throwaway Postgres so the seeding CLI can be exercised
// against a real database (its logic is not reachable through the HTTP API).
func TestMain(m *testing.M) {
	ctx := context.Background()
	os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")

	dbName := envOr("DB_NAME", "app_db")
	dbUser := envOr("DB_USER", "app")
	dbPass := envOr("DB_PASSWORD", "app")

	pg, err := tcpostgres.Run(ctx, "postgres:15",
		tcpostgres.WithDatabase(dbName),
		tcpostgres.WithUsername(dbUser),
		tcpostgres.WithPassword(dbPass),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		log.Fatalf("failed to start postgres: %v", err)
	}

	host, _ := pg.Host(ctx)
	port, _ := pg.MappedPort(ctx, "5432/tcp")
	os.Setenv("DB_HOST", host)
	os.Setenv("DB_PORT", port.Port())
	os.Setenv("DB_NAME", dbName)
	os.Setenv("DB_USER", dbUser)
	os.Setenv("DB_PASSWORD", dbPass)

	code := m.Run()

	if err := testcontainers.TerminateContainer(pg); err != nil {
		log.Printf("failed to terminate postgres: %v", err)
	}
	os.Exit(code)
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func connect(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := config.ConnectDB()
	if err != nil {
		t.Fatalf("connect db: %v", err)
	}
	if err := ensureSchema(db); err != nil {
		t.Fatalf("ensure schema: %v", err)
	}
	return db
}

func TestHashPassword(t *testing.T) {
	hash := hashPassword("StrongPass123!")
	if hash == "" || hash == "StrongPass123!" {
		t.Fatal("password must be hashed")
	}
	if !utils.CheckHashString("StrongPass123!", hash) {
		t.Fatal("hashed password should verify against the original")
	}
}

func TestGetMimeType(t *testing.T) {
	cases := map[string]string{
		"a.jpg":     "image/jpeg",
		"a.jpeg":    "image/jpeg",
		"a.png":     "image/png",
		"a.gif":     "image/gif",
		"a.avif":    "image/avif",
		"a.webp":    "image/webp",
		"a.unknown": "application/octet-stream",
		"noext":     "application/octet-stream",
	}
	for name, want := range cases {
		if got := getMimeType(name); got != want {
			t.Errorf("getMimeType(%q) = %q, want %q", name, got, want)
		}
	}
}

func TestCreateFileRecord_EmptyName(t *testing.T) {
	db := connect(t)
	url, err := createFileRecord(db, "00000000-0000-0000-0000-000000000000", "", "avatars")
	if err != nil || url != "" {
		t.Fatalf("empty source name should return empty url and no error, got (%q,%v)", url, err)
	}
}

// TestSeedPipeline runs the full seeding sequence against the real database and
// asserts the expected data was created. Running twice also exercises the
// idempotency branches (records that already exist are skipped).
func TestSeedPipeline(t *testing.T) {
	db := connect(t)

	users := seedUsers(db)
	if len(users) == 0 {
		t.Fatal("expected users to be seeded from users.json")
	}

	seedFriendships(db, users)
	seedFollows(db, users)
	posts := seedPosts(db, users)
	if len(posts) == 0 {
		t.Fatal("expected posts to be seeded")
	}
	seedLikes(db, users, posts)
	seedReplies(db, users, posts)
	seedNotifications(db, users)

	// second pass: every seeder should be idempotent (no duplicates, no panic)
	users2 := seedUsers(db)
	seedFriendships(db, users2)
	seedFollows(db, users2)
	posts2 := seedPosts(db, users2)
	seedLikes(db, users2, posts2)
	seedReplies(db, users2, posts2)
	seedNotifications(db, users2)
}

// TestRunMain exercises the top-level orchestration exactly as the CLI does.
func TestRunMain(_ *testing.T) {
	main()
}

func TestSeeders_EmptyInputsAreNoOps(t *testing.T) {
	db := connect(t)
	var none []models.User
	// guard branches for too-few users / no posts must not panic
	seedFriendships(db, none)
	seedFollows(db, none)
	seedLikes(db, none, nil)
	seedReplies(db, none, nil)
	seedNotifications(db, none)
	if got := seedPosts(db, none); got != nil {
		t.Fatalf("seedPosts with no users should return nil, got %v", got)
	}
}
