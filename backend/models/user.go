package models

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	ID        string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	GithubID  *string        `gorm:"type:varchar(255);uniqueIndex" json:"github_id"`
	Provider  string         `gorm:"type:varchar(50);default:'local'" json:"provider"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	Name     string `json:"name"`
	Username string `json:"username" binding:"required" gorm:"unique;not null"`
	Email    string `json:"email" binding:"required,email" gorm:"unique;not null"`
	Password *string `json:"password"`

	DateOfBirth *time.Time `json:"dateOfBirth"`

	Wallpaper *string `json:"wallpaper"`
	Avatar    *string `json:"avatar"`
	Bio       string `json:"bio"`
}

type UpdateUserInput struct {
	Name      string `json:"name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Bio       string `json:"bio"`
	Avatar    *string `json:"avatar"`
	Wallpaper *string `json:"wallpaper"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Name      string    `json:"name,omitempty"`
	Bio       string    `json:"bio,omitempty"`
	Avatar    *string    `json:"avatar,omitempty"`
	Wallpaper *string    `json:"wallpaper,omitempty"`
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
	ID       string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID   string `gorm:"type:uuid;not null;index"`
	FriendID string `gorm:"type:uuid;not null;index"`
	Status   string
}
