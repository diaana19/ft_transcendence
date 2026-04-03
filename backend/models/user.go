package models

import "time"

type User struct {
	ID      	string    `json:"id"`
	Name    	string    `json:"name"`
	Username	string    `json:"username"`
	Password	string    `json:"-"`
	Email    	string    `json:"email"`
	Wallpaper	string    `json:"wallpaper"`
	Avatar   	string    `json:"avatar"`
	Bio      	string    `json:"bio"`
	CreatedAt	time.Time `json:"created_at"`
	UpdatedAt	time.Time `json:"updated_at"`
	DeletedAt	time.Time `json:"deleted_at"`
}