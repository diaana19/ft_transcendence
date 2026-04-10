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
)

func hashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func main() {
	db, err := config.ConnectDB()
	if err != nil {
		panic(err)
	}

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
			fmt.Println("Error inserting:", err)
		} else {
			fmt.Println("Inserted:", users[i].Email)
		}
	}

	fmt.Println("✅ Seeding Finished")
}