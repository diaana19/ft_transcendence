package models

import (
	"time"

	"gorm.io/gorm"
)

// User is the account of a person.
type User struct {
	ID        string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	GithubID  *string        `gorm:"type:varchar(255);uniqueIndex" json:"github_id"`
	Provider  string         `gorm:"type:varchar(50);default:'local'" json:"provider"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitzero" gorm:"index"`

	DisplayName string  `json:"displayname"`
	Username    string  `json:"username" binding:"required" gorm:"unique;not null"`
	Email       string  `json:"email" binding:"required,email" gorm:"unique;not null"`
	Password    *string `json:"-"`

	DateOfBirth *time.Time `json:"dateOfBirth"`

	TwoFASecret  *string `gorm:"type:varchar(255)" json:"-"`
	TwoFAEnabled bool    `gorm:"default:false" json:"two_fa_enabled"`

	Wallpaper *string `json:"wallpaper"`
	Avatar    *string `json:"avatar"`
	Bio       string  `json:"bio"`
}

// UpdateUserInput holds the fields a user can change in the profile.
type UpdateUserInput struct {
	DisplayName string  `json:"displayname"`
	Username    string  `json:"username"`
	Email       string  `json:"email"`
	Bio         string  `json:"bio"`
	Avatar      *string `json:"avatar"`
	Wallpaper   *string `json:"wallpaper"`
}

// UserResponse is the user data we send to the client.
type UserResponse struct {
	ID             string    `json:"id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	DisplayName    string    `json:"displayname,omitempty"`
	Bio            string    `json:"bio,omitempty"`
	Avatar         *string   `json:"avatar,omitempty"`
	Wallpaper      *string   `json:"wallpaper,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	TwoFAEnabled   bool      `json:"two_fa_enabled"`
	FollowersCount int64     `json:"followers_count"`
	FollowingCount int64     `json:"following_count"`
}

// ToResponse converts the User into a UserResponse.
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:           u.ID,
		Username:     u.Username,
		Email:        u.Email,
		DisplayName:  u.DisplayName,
		Bio:          u.Bio,
		Avatar:       u.Avatar,
		Wallpaper:    u.Wallpaper,
		CreatedAt:    u.CreatedAt,
		TwoFAEnabled: u.TwoFAEnabled,
	}
}

// Friend is the relation between two users.
type Friend struct {
	ID       string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID   string `gorm:"type:uuid;not null;index" json:"user_id"`
	FriendID string `gorm:"type:uuid;not null;index" json:"friend_id"`
	Status   string `json:"status"`
}
