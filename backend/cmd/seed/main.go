package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Transcendence/config"
	"github.com/Transcendence/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// --- Data ---

var postContents = []string{
	"Just finished an amazing project! Feeling proud 🎉",
	"Anyone else love coding at 3am? Coffee is life ☕",
	"Check out this cool new tech stack I'm learning about!",
	"Weekend is here! Time to build something cool 🚀",
	"Finally deployed to production! No more bugs (hopefully) 😅",
	"Open source contributions are the best way to learn",
	"Just discovered this incredible library, game changer!",
	"Working on something secret, can't wait to share soon...",
	"The debugging journey never ends... but that's what makes it fun!",
	"New blog post is live! Check it out, feedback welcome 📝",
	"Sometimes the best code is the code you delete 🗑️",
	"Excited to announce we're hiring! Great team, great mission 💼",
}

var commentContents = []string{
	"Great post!",
	"Totally agree.",
	"Could you share more details?",
	"Thanks for sharing!",
	"This is awesome!",
	"Same here haha",
	"Love this perspective.",
	"Interesting take!",
}

// --- Helpers ---

func hashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

type seedUserInput struct {
	Name        string    `json:"name"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	DateOfBirth time.Time `json:"dateOfBirth"`
	Bio         string    `json:"bio"`
	Wallpaper   string    `json:"wallpaper"`
	Avatar      string    `json:"avatar"`
}

// toUser builds a User without avatar/wallpaper.
// Those are set later by seedFiles (which uploads the actual files into the storage system).
func (s seedUserInput) toUser() models.User {
	hashed := hashPassword(s.Password)
	dob := s.DateOfBirth
	now := time.Now()

	return models.User{
		ID:          uuid.NewString(),
		Name:        s.Name,
		Username:    s.Username,
		Email:       s.Email,
		Password:    &hashed,
		DateOfBirth: &dob,
		Bio:         s.Bio,
		Provider:    "local",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// --- File seeding ---

func seedFileForUser(db *gorm.DB, userID, sourcePath string) (string, error) {
	if sourcePath == "" {
		return "", nil
	}

	relativeSourcePath := strings.TrimPrefix(sourcePath, "/")
	if relativeSourcePath == "" {
		return "", nil
	}

	srcInfo, err := os.Stat(relativeSourcePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("    Source file not found: %s (skipping)\n", relativeSourcePath)
			return "", nil
		}
		return "", err
	}

	fileID := uuid.NewString()
	ext := strings.ToLower(filepath.Ext(relativeSourcePath))
	destName := fileID + ext
	destPath := filepath.Join("uploads", destName)

	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	if err := copyFile(relativeSourcePath, destPath); err != nil {
		return "", fmt.Errorf("copy file: %w", err)
	}

	file := models.File{
		ID:         fileID,
		OwnerID:    userID,
		Path:       destPath,
		Filename:   filepath.Base(relativeSourcePath),
		MimeType:   mimeType,
		Size:       srcInfo.Size(),
		Visibility: models.FileVisibilityPublic,
		CreatedAt:  time.Now(),
	}
	if err := db.Create(&file).Error; err != nil {
		return "", fmt.Errorf("create file row: %w", err)
	}

	url := "/api/files/" + fileID
	return url, nil
}

func copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return nil
}

func seedUserMedia(db *gorm.DB, user *models.User, sourceAvatar, sourceWallpaper string) {
	if user.Avatar == nil && sourceAvatar != "" {
		url, err := seedFileForUser(db, user.ID, sourceAvatar)
		if err != nil {
			fmt.Printf("    Error seeding avatar for %s: %v\n", user.Username, err)
		} else if url != "" {
			db.Model(user).Update("avatar", url)
			fmt.Printf("    + Avatar: %s -> %s\n", sourceAvatar, url)
		}
	}

	if user.Wallpaper == nil && sourceWallpaper != "" {
		url, err := seedFileForUser(db, user.ID, sourceWallpaper)
		if err != nil {
			fmt.Printf("    Error seeding wallpaper for %s: %v\n", user.Username, err)
		} else if url != "" {
			db.Model(user).Update("wallpaper", url)
			fmt.Printf("    + Wallpaper: %s -> %s\n", sourceWallpaper, url)
		}
	}
}

// --- Seeders ---

func ensureSchema(db *gorm.DB) error {
	fmt.Println("Ensuring schema is up to date...")
	return db.AutoMigrate(
		&models.User{},
		&models.Friend{},
		&models.Post{},
		&models.Like{},
		&models.Reply{},
		&models.Repost{},
		&models.Message{},
		&models.Notification{},
		&models.File{},
		&models.FileAccess{},
	)
}

func seedUsers(db *gorm.DB) []models.User {
	fmt.Println("\nSeeding users...")

	file, err := os.ReadFile("users.json")
	if err != nil {
		panic(err)
	}

	var inputs []seedUserInput
	if err := json.Unmarshal(file, &inputs); err != nil {
		panic(err)
	}

	for _, in := range inputs {
		var existing models.User
		if err := db.Where("email = ? OR username = ?", in.Email, in.Username).First(&existing).Error; err == nil {
			fmt.Println("  - User already exists:", in.Email)
			seedUserMedia(db, &existing, in.Avatar, in.Wallpaper)
			continue
		}

		user := in.toUser()
		if err := db.Create(&user).Error; err != nil {
			fmt.Println("  - Error inserting user:", err)
			continue
		}
		fmt.Printf("  + Inserted user: %s\n", user.Username)

		seedUserMedia(db, &user, in.Avatar, in.Wallpaper)
	}

	var users []models.User
	db.Find(&users)
	return users
}

func seedFriendships(db *gorm.DB, users []models.User) {
	fmt.Println("\nSeeding friendships...")

	n := len(users)
	if n < 3 {
		fmt.Println("  - Need at least 3 users, skipping")
		return
	}

	accepted := 0
	pending := 0
	for i := 0; i < n; i++ {
		acceptedTarget := users[(i+1)%n]
		if !friendExists(db, users[i].ID, acceptedTarget.ID) {
			db.Create(&models.Friend{
				ID:       uuid.NewString(),
				UserID:   users[i].ID,
				FriendID: acceptedTarget.ID,
				Status:   "accepted",
			})
			accepted++
		}

		pendingTarget := users[(i+2)%n]
		if !friendExists(db, users[i].ID, pendingTarget.ID) {
			db.Create(&models.Friend{
				ID:       uuid.NewString(),
				UserID:   users[i].ID,
				FriendID: pendingTarget.ID,
				Status:   "pending",
			})
			pending++
		}
	}

	fmt.Printf("  + Created %d accepted friendships, %d pending requests\n", accepted, pending)
}

func seedFollows(db *gorm.DB, users []models.User) {
	fmt.Println("\nSeeding follows...")

	n := len(users)
	if n < 2 {
		fmt.Println("  - Need at least 2 users, skipping")
		return
	}

	created := 0
	for i := 0; i < n; i++ {
		for offset := 1; offset <= 2 && offset < n; offset++ {
			target := users[(i+offset)%n]
			if friendExistsWithStatus(db, users[i].ID, target.ID, "follow") {
				continue
			}
			db.Create(&models.Friend{
				ID:       uuid.NewString(),
				UserID:   users[i].ID,
				FriendID: target.ID,
				Status:   "follow",
			})
			created++
		}
	}

	fmt.Printf("  + Created %d follow relationships\n", created)
}

func seedPosts(db *gorm.DB, users []models.User) []models.Post {
	fmt.Println("\nSeeding posts...")

	if len(users) == 0 {
		fmt.Println("  - No users, skipping")
		return nil
	}

	for contentIdx, content := range postContents {
		user := users[contentIdx%len(users)]
		post := models.Post{
			ID:        uuid.NewString(),
			Content:   content,
			AuthorID:  user.ID,
			CreatedAt: time.Now().Add(-time.Duration((len(postContents)-contentIdx)*24) * time.Hour),
			UpdatedAt: time.Now().Add(-time.Duration((len(postContents)-contentIdx)*24) * time.Hour),
		}

		var existing models.Post
		if err := db.Where("content = ? AND author_id = ?", post.Content, post.AuthorID).First(&existing).Error; err == nil {
			continue
		}

		if err := db.Create(&post).Error; err != nil {
			fmt.Println("  - Error inserting post:", err)
		}
	}

	var posts []models.Post
	db.Find(&posts)
	fmt.Printf("  + Total posts: %d\n", len(posts))
	return posts
}

func seedLikes(db *gorm.DB, users []models.User, posts []models.Post) {
	fmt.Println("\nSeeding likes...")

	if len(users) == 0 || len(posts) == 0 {
		fmt.Println("  - Need users and posts, skipping")
		return
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	created := 0
	for _, post := range posts {
		for _, user := range users {
			if user.ID == post.AuthorID {
				continue
			}
			if rng.Float64() > 0.5 {
				continue
			}

			var existing models.Like
			if err := db.Where("post_id = ? AND user_id = ?", post.ID, user.ID).First(&existing).Error; err == nil {
				continue
			}

			db.Create(&models.Like{
				ID:     uuid.NewString(),
				PostID: post.ID,
				UserID: user.ID,
			})
			created++

			db.Model(&models.Post{}).
				Where("id = ?", post.ID).
				UpdateColumn("likes_count", gorm.Expr("likes_count + ?", 1))
		}
	}

	fmt.Printf("  + Created %d likes\n", created)
}

func seedReplies(db *gorm.DB, users []models.User, posts []models.Post) {
	fmt.Println("\nSeeding replies (comments)...")

	if len(users) < 2 || len(posts) == 0 {
		fmt.Println("  - Need users and posts, skipping")
		return
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano() + 1))
	created := 0
	for i, post := range posts {
		count := 1 + rng.Intn(3)
		for j := 0; j < count; j++ {
			authorIdx := (i + j + 1) % len(users)
			author := users[authorIdx]
			if author.ID == post.AuthorID {
				authorIdx = (authorIdx + 1) % len(users)
				author = users[authorIdx]
			}

			content := commentContents[(i+j)%len(commentContents)]

			var existing models.Reply
			if err := db.Where("post_id = ? AND author_id = ? AND content = ?",
				post.ID, author.ID, content).First(&existing).Error; err == nil {
				continue
			}

			db.Create(&models.Reply{
				ID:       uuid.NewString(),
				PostID:   post.ID,
				AuthorID: author.ID,
				Content:  content,
			})
			created++

			db.Model(&models.Post{}).
				Where("id = ?", post.ID).
				UpdateColumn("comments_count", gorm.Expr("comments_count + ?", 1))
		}
	}

	fmt.Printf("  + Created %d replies\n", created)
}

func seedNotifications(db *gorm.DB, users []models.User) {
	fmt.Println("\nSeeding notifications...")

	if len(users) < 2 {
		fmt.Println("  - Need at least 2 users, skipping")
		return
	}

	created := 0
	for i, user := range users {
		actor := users[(i+len(users)-1)%len(users)]

		var existing models.Notification
		if err := db.Where("user_id = ? AND actor_id = ? AND type = ?",
			user.ID, actor.ID, "friend_request").First(&existing).Error; err == nil {
			continue
		}

		db.Create(&models.Notification{
			ID:            uuid.NewString(),
			CreatedAt:     time.Now().Add(-time.Hour),
			UserID:        user.ID,
			UserUsername:  user.Username,
			ActorID:       actor.ID,
			ActorUsername: actor.Username,
			Type:          "friend_request",
			Content:       actor.Username + " sent you a friend request",
			Read:          false,
		})
		created++
	}

	fmt.Printf("  + Created %d notifications\n", created)
}

// --- Existence checks ---

func friendExists(db *gorm.DB, userID, friendID string) bool {
	var count int64
	db.Model(&models.Friend{}).
		Where("user_id = ? AND friend_id = ?", userID, friendID).
		Count(&count)
	return count > 0
}

func friendExistsWithStatus(db *gorm.DB, userID, friendID, status string) bool {
	var count int64
	db.Model(&models.Friend{}).
		Where("user_id = ? AND friend_id = ? AND status = ?", userID, friendID, status).
		Count(&count)
	return count > 0
}

// --- Main ---

func main() {
	db, err := config.ConnectDB()
	if err != nil {
		panic(err)
	}

	if err := ensureSchema(db); err != nil {
		panic(fmt.Errorf("schema migration failed: %w", err))
	}

	users := seedUsers(db)
	seedFriendships(db, users)
	seedFollows(db, users)
	posts := seedPosts(db, users)
	seedLikes(db, users, posts)
	seedReplies(db, users, posts)
	seedNotifications(db, users)

	fmt.Println("\nSeeding finished!")
}
