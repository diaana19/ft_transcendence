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
