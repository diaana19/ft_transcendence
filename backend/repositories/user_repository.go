package repositories

import (
	"github.com/Transcendence/models"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

type UserRepository interface {
	GetAll() ([]models.User, error)
	GetByID(id string) (*models.User, error) 
	Update(id string, input models.UpdateUserInput) (*models.User, error)
	Delete(id string) error
	CreateUser(user *models.User) error
	GetByEmail(email string) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetAll() ([]models.User, error) { 
	var users []models.User
	result := r.db.Find(&users)
	return users, result.Error
}

func (r *userRepository) GetByID(id string) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, "id = ?", id)
	return &user, result.Error
}

func (r *userRepository) Update(id string, input models.UpdateUserInput) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	result = r.db.Model(&user).Updates(input)
	return &user, result.Error
}

func (r *userRepository) Delete(id string) error {
	result := r.db.Delete(&models.User{}, "id = ?", id)
	return result.Error
}

func (r *userRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, "email = ?", email)
	return &user, result.Error
}

func (r *userRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, "username = ?", username)
	return &user, result.Error
}