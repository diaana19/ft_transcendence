package models

import "time"

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Wallpaper string    `json:"wallpaper"`
	Avatar    string    `json:"avatar"`
	Bio       string    `json:"bio"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type UpdateUserInput struct {
	Name      string `json:"name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Bio       string `json:"bio"`
	Avatar    string `json:"avatar"`
	Wallpaper string `json:"wallpaper"`
}
