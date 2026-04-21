package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Transcendence/models"
	"github.com/Transcendence/config"
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
			ID:       uuid.NewString(),
			Content:  content,
			AuthorID: users[contentIdx%len(users)].ID,
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

func main() {
	db, err := config.ConnectDB()
	if err != nil {
		panic(err)
	}

	// Seed Users
	fmt.Println("🌱 Seeding users...")
	file, err := os.ReadFile("users.json")
	if err != nil {
		panic(err)
	}

	var users []models.User
	if err := json.Unmarshal(file, &users); err != nil {
		panic(err)
	}

	for i := range users {
		users[i].ID = uuid.NewString()
		users[i].Password = hashPassword(users[i].Password)

		if users[i].CreatedAt.IsZero() {
			users[i].CreatedAt = time.Now()
		}
		if users[i].UpdatedAt.IsZero() {
			users[i].UpdatedAt = time.Now()
		}

		var existing models.User
		err := db.Where("email = ? OR username = ?", users[i].Email, users[i].Username).First(&existing).Error

		if err == nil {
			fmt.Println("User already exists:", users[i].Email)
			continue
		}

		if err := db.Create(&users[i]).Error; err != nil {
			fmt.Println("Error inserting user:", err)
		} else {
			fmt.Println("✓ Inserted user:", users[i].Email)
		}
	}

	// Seed Posts
	fmt.Println("\n🌱 Seeding posts...")
	seedPosts(db)

	fmt.Println("\n✅ Seeding finished!")
}