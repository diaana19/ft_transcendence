package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/Transcendence/config"
	"github.com/Transcendence/models"
)

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

// Mapping exact des utilisateurs aux fichiers réels
var avatarMap = map[string]string{
	"alice01":  "av1.jpg",
	"carlos02": "av2.jpg",
	"sophie03": "av3.jpg",
	"john04":   "av4.avif",
	"emma05":   "av5.jpg",
	"lucas06":  "av6.jpeg",
	"mia07":    "av7.jpg",
	"ahmed08":  "av1.jpg",
	"olivia09": "av2.jpg",
	"noah10":   "av3.jpg",
}

func hashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

const (
	mimeJPEG        = "image/jpeg"
	mimePNG         = "image/png"
	mimeGIF         = "image/gif"
	mimeAVIF        = "image/avif"
	mimeWebP        = "image/webp"
	mimeOctetStream = "application/octet-stream"
)

func getMimeType(filename string) string {
	ext := filepath.Ext(filename)
	switch ext {
	case ".jpg", ".jpeg":
		return mimeJPEG
	case ".png":
		return mimePNG
	case ".gif":
		return mimeGIF
	case ".avif":
		return mimeAVIF
	case ".webp":
		return mimeWebP
	default:
		return mimeOctetStream
	}
}

// createFileRecord : crée un File record qui pointe vers un fichier existent
// N'utilise PAS l'UUID comme nom de fichier - garde le nom original
func createFileRecord(db *gorm.DB, userID, sourceFileName, subDir string) (string, error) {
	if sourceFileName == "" {
		return "", nil
	}

	fileID := uuid.NewString()

	// Le chemin en BDD pointe vers le fichier EXISTENT
	dbPath := filepath.Join("uploads", subDir, sourceFileName)

	fileRecord := models.File{
		ID:         fileID,
		Filename:   sourceFileName,
		Path:       dbPath,
		MimeType:   getMimeType(sourceFileName),
		Size:       0,
		OwnerID:    userID,
		Visibility: models.FileVisibilityPublic,
	}

	if err := db.Create(&fileRecord).Error; err != nil {
		return "", err
	}

	fmt.Printf("    ✅ Registered %s (%s) -> /api/files/%s\n", subDir, sourceFileName, fileID)
	return "/api/files/" + fileID, nil
}

type seedUserInput struct {
	Name        string `json:"name"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	DateOfBirth string `json:"dateOfBirth"`
	Bio         string `json:"bio"`
}

func seedUsers(db *gorm.DB) []models.User {
	fmt.Println("Seeding users...")
	var users []models.User

	data, err := os.ReadFile("users.json")
	if err != nil {
		panic(err)
	}

	var inputs []seedUserInput
	if err := json.Unmarshal(data, &inputs); err != nil {
		panic(err)
	}

	for _, input := range inputs {
		user := models.User{
			ID:        uuid.NewString(),
			Name:      input.Name,
			Username:  input.Username,
			Email:     input.Email,
			Bio:       input.Bio,
			Provider:  "local",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if input.Password != "" {
			hashed := hashPassword(input.Password)
			user.Password = &hashed
		}

		var existing models.User
		if err := db.Where("email = ? OR username = ?", user.Email, user.Username).First(&existing).Error; err == nil {
			fmt.Printf("  - User already exists: %s\n", user.Email)
			users = append(users, existing)
			continue
		}

		if err := db.Create(&user).Error; err != nil {
			fmt.Printf("❌ Failed to create user %s: %v\n", user.Username, err)
			continue
		}
		fmt.Printf("  + Inserted user: %s\n", user.Username)

		// Assigner l'avatar depuis la map
		if avatarFile, ok := avatarMap[user.Username]; ok {
			avatarURL, _ := createFileRecord(db, user.ID, avatarFile, "avatars")
			if avatarURL != "" {
				db.Model(&user).Update("avatar", avatarURL)
			}
		}

		// Assigner le wallpaper (toujours default.jpg)
		wallpaperURL, _ := createFileRecord(db, user.ID, "default.jpg", "wallpapers")
		if wallpaperURL != "" {
			db.Model(&user).Update("wallpaper", wallpaperURL)
		}

		users = append(users, user)
	}

	return users
}

func ensureSchema(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&models.User{}, &models.File{}, &models.Post{}, &models.Like{},
		&models.Reply{}, &models.Friend{}, &models.Notification{},
	); err != nil {
		return fmt.Errorf("auto-migrate: %w", err)
	}
	return nil
}

func seedFriendships(db *gorm.DB, users []models.User) {
	fmt.Println("Seeding friendships...")
	n := len(users)
	if n < 3 {
		fmt.Println("  - Need at least 3 users, skipping")
		return
	}

	accepted := 0
	pending := 0
	for i := range n {
		acceptedTarget := users[(i+1)%n]
		var count int64
		db.Model(&models.Friend{}).
			Where("user_id = ? AND friend_id = ?", users[i].ID, acceptedTarget.ID).
			Count(&count)
		if count == 0 {
			db.Create(&models.Friend{
				ID:       uuid.NewString(),
				UserID:   users[i].ID,
				FriendID: acceptedTarget.ID,
				Status:   "accepted",
			})
			accepted++
		}

		pendingTarget := users[(i+2)%n]
		db.Model(&models.Friend{}).
			Where("user_id = ? AND friend_id = ?", users[i].ID, pendingTarget.ID).
			Count(&count)
		if count == 0 {
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
	fmt.Println("Seeding follows...")
	n := len(users)
	if n < 2 {
		fmt.Println("  - Need at least 2 users, skipping")
		return
	}

	created := 0
	for i := range n {
		for offset := 1; offset <= 2 && offset < n; offset++ {
			target := users[(i+offset)%n]
			var count int64
			db.Model(&models.Friend{}).
				Where("user_id = ? AND friend_id = ? AND status = ?", users[i].ID, target.ID, "follow").
				Count(&count)
			if count == 0 {
				db.Create(&models.Friend{
					ID:       uuid.NewString(),
					UserID:   users[i].ID,
					FriendID: target.ID,
					Status:   "follow",
				})
				created++
			}
		}
	}

	fmt.Printf("  + Created %d follow relationships\n", created)
}

func seedPosts(db *gorm.DB, users []models.User) []models.Post {
	fmt.Println("Seeding posts...")
	if len(users) == 0 {
		fmt.Println("  - No users, skipping")
		return nil
	}

	var posts []models.Post
	for contentIdx, content := range postContents {
		user := users[contentIdx%len(users)]
		post := models.Post{
			ID:       uuid.NewString(),
			Content:  content,
			AuthorID: user.ID,
		}

		var existing models.Post
		if err := db.Where("content = ? AND author_id = ?", post.Content, post.AuthorID).First(&existing).Error; err == nil {
			continue
		}

		if err := db.Create(&post).Error; err != nil {
			fmt.Printf("  - Error creating post: %v\n", err)
			continue
		}
		posts = append(posts, post)
	}

	fmt.Printf("  + Created %d posts\n", len(posts))
	return posts
}

func seedLikes(db *gorm.DB, users []models.User, posts []models.Post) {
	fmt.Println("Seeding likes...")
	if len(users) == 0 || len(posts) == 0 {
		fmt.Println("  - Need users and posts, skipping")
		return
	}

	created := 0
	for _, post := range posts {
		for _, user := range users {
			if user.ID == post.AuthorID {
				continue
			}
			if (len(post.ID)+len(user.ID))%2 == 0 {
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
	fmt.Println("Seeding replies (comments)...")
	if len(users) < 2 || len(posts) == 0 {
		fmt.Println("  - Need users and posts, skipping")
		return
	}

	created := 0
	for i, post := range posts {
		count := 1 + (i % 2)
		for j := range count {
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
	fmt.Println("Seeding notifications...")
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

func main() {
	db, err := config.ConnectDB()
	if err != nil {
		panic(err)
	}

	if err := ensureSchema(db); err != nil {
		panic(err)
	}

	users := seedUsers(db)
	seedFriendships(db, users)
	seedFollows(db, users)
	posts := seedPosts(db, users)
	seedLikes(db, users, posts)
	seedReplies(db, users, posts)
	seedNotifications(db, users)

	fmt.Println("\n✅ Seeding finished successfully!")
}
