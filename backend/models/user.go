package models

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	ID        string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	Name     string `json:"name"`
	Username string `json:"username" binding:"required" gorm:"unique;not null"`
	Email    string `json:"email" binding:"required,email" gorm:"unique;not null"`
	Password string `json:"password" binding:"required,min=8"`

	DateOfBirth time.Time `json:"dateOfBirth" binding:"required"`

	Wallpaper string `json:"wallpaper"`
	Avatar    string `json:"avatar"`
	Bio       string `json:"bio"`
}

type UpdateUserInput struct {
	Name      string `json:"name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Bio       string `json:"bio"`
	Avatar    string `json:"avatar"`
	Wallpaper string `json:"wallpaper"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Name      string    `json:"name,omitempty"`
	Bio       string    `json:"bio,omitempty"`
	Avatar    string    `json:"avatar,omitempty"`
	Wallpaper string    `json:"wallpaper,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Name:      u.Name,
		Bio:       u.Bio,
		Avatar:    u.Avatar,
		Wallpaper: u.Wallpaper,
		CreatedAt: u.CreatedAt,
	}
}

type Friend struct {
    ID        uint `gorm:"primaryKey"`
    UserID    uint
    FriendID  uint
    Status    string // "pending", "accepted"
}
