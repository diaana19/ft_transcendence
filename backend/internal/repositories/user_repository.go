package repositories

import (
	"gorm.io/gorm"

	"ft_transcendence/backend/internal/models"
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
	GetByIdentifier(identifier string) (*models.User, error)
	GetByGithubID(githubID string) (*models.User, error)
	LinkGithub(userID, githubID string) error
	SetTwoFASecret(userID, secret string) error
	SetTwoFAEnabled(userID string, enabled bool) error
	ClearTwoFA(userID string) error
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
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}

	updates := make(map[string]any)
	if input.Name != "" {
		updates["name"] = input.Name
	}
	if input.Username != "" {
		updates["username"] = input.Username
	}
	if input.Email != "" {
		updates["email"] = input.Email
	}
	if input.Bio != "" {
		updates["bio"] = input.Bio
	}
	if input.Avatar != nil {
		updates["avatar"] = *input.Avatar
	}
	if input.Wallpaper != nil {
		updates["wallpaper"] = *input.Wallpaper
	}

	if len(updates) > 0 {
		if err := r.db.Model(&user).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Delete(id string) error {
	result := r.db.Delete(&models.User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
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

func (r *userRepository) GetByIdentifier(identifier string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ? OR username = ?", identifier, identifier).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByGithubID(githubID string) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, "github_id = ?", githubID)
	return &user, result.Error
}

func (r *userRepository) LinkGithub(userID, githubID string) error {
	result := r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]any{
			"github_id": githubID,
			"provider":  "github",
		})
	return result.Error
}

func (r *userRepository) SetTwoFASecret(userID, secret string) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("two_fa_secret", secret).Error
}

func (r *userRepository) SetTwoFAEnabled(userID string, enabled bool) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("two_fa_enabled", enabled).Error
}

func (r *userRepository) ClearTwoFA(userID string) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]any{
			"two_fa_secret":  nil,
			"two_fa_enabled": false,
		}).Error
}
