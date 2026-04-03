package services

import (
	"github.com/Transcendence/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}	
}

type CreateUserInput struct {
    Name     string `json:"name"     binding:"required"`
    Username string `json:"username" binding:"required"`
    Email    string `json:"email"    binding:"required"`
    Password string `json:"password" binding:"required"`
    Bio      string `json:"bio"`
}

type UpdateUserInput struct {
	Name     string	`json:"name"`
	Username string	`json:"username"`
	Email    string	`json:"email"`
	Bio      string	`json:"bio"`
	Avatar   string	`json:"avatar"`
	Wallpaper string `json:"wallpaper"`
}

func (s *UserService) GetUsers() ([]models.User, error) {
	var users []models.User
	result := s.db.Find(&users)
	return users, result.Error
}

func (s *UserService) GetUser(id string) (*models.User, error) {
	var user models.User
	result := s.db.First(&user, "id = ?", id)
	return &user, result.Error
}

func (s *UserService) CreateUser(input CreateUserInput) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if (err != nil) {
		return nil, err
	}
	user := models.User {
		Name :	input.Name,
		Username: input.Username,
		Email: input.Email,
		Password: string(hashedPassword),
	}
	err = s.db.Create(&user).Error
	if (err != nil) {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) UpdateUser(id string, input UpdateUserInput) (*models.User, error) {
	var user models.User
	result := s.db.First(&user, "id = ?", id)
	if (result.Error != nil) {
		return nil, result.Error
	}
	result = s.db.Model(&user).Updates(input)
	return &user, result.Error
}

func (s *UserService) DeleteUser(id string) error {
	result := s.db.Delete(&models.User{}, "id = ?", id)
	return result.Error
}

