package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model          `gorm:"embedded"`
	Name      string    `json:"name"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"password"`
	Wallpaper string    `json:"wallpaper"`
	Avatar    string    `json:"avatar"`
	Bio       string    `json:"bio"`
}

type UpdateUserInput struct {
	Name      string `json:"name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Bio       string `json:"bio"`
	Avatar    string `json:"avatar"`
	Wallpaper string `json:"wallpaper"`
}
