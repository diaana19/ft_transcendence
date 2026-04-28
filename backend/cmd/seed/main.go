package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Transcendence/config"
	"github.com/Transcendence/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func hashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

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

func seedPosts(db *gorm.DB) {
	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		fmt.Println("Error fetching users:", err)
		return
	}

	if len(users) == 0 {
		fmt.Println("No users found, skipping post seeding")
		return
	}

	for contentIdx, content := range postContents {
		post := models.Post{
			ID:        uuid.NewString(),
			Content:   content,
			AuthorID:  users[contentIdx%len(users)].ID,
			CreatedAt: time.Now().Add(-time.Duration((len(postContents)-contentIdx)*24) * time.Hour),
			UpdatedAt: time.Now().Add(-time.Duration((len(postContents)-contentIdx)*24) * time.Hour),
		}

		var existing models.Post
		if err := db.Where("content = ? AND author_id = ?", post.Content, post.AuthorID).First(&existing).Error; err == nil {
			fmt.Println("Post already exists for user:", users[contentIdx%len(users)].Username)
			continue
		} 

		if err := db.Create(&post).Error; err != nil {
			fmt.Println("Error inserting post:", err)
		} else {
			fmt.Printf("✓ Inserted post for user: %s\n", users[contentIdx%len(users)].Username)
		}
	}
}

func ensureSchema(db *gorm.DB) error {
	fmt.Println("🔧 Ensuring schema is up to date...")
	return db.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Like{},
		&models.Reply{},
		&models.Repost{},
	)
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

func (s seedUserInput) toUser() models.User {
	hashed := hashPassword(s.Password)
	dob := s.DateOfBirth
	wallpaper := s.Wallpaper
	avatar := s.Avatar
	now := time.Now()

	return models.User{
		ID:          uuid.NewString(),
		Name:        s.Name,
		Username:    s.Username,
		Email:       s.Email,
		Password:    &hashed,
		DateOfBirth: &dob,
		Bio:         s.Bio,
		Wallpaper:   &wallpaper,
		Avatar:      &avatar,
		Provider:    "local",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func main() {
	db, err := config.ConnectDB()
	if err != nil {
		panic(err)
	}

	if err := ensureSchema(db); err != nil {
		panic(fmt.Errorf("schema migration failed: %w", err))
	}

	fmt.Println("\n🌱 Seeding users...")
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
		err := db.Where("email = ? OR username = ?", in.Email, in.Username).First(&existing).Error

		if err == nil {
			fmt.Println("User already exists:", in.Email)
			continue
		}

		user := in.toUser()
		if err := db.Create(&user).Error; err != nil {
			fmt.Println("Error inserting user:", err)
		} else {
			fmt.Println("✓ Inserted user:", user.Email)
		}
	}

	fmt.Println("\n🌱 Seeding posts...")
	seedPosts(db)

	fmt.Println("\n✅ Seeding finished!")
}
