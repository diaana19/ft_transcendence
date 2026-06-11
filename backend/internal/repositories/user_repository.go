package repositories

import (
	"gorm.io/gorm"

	"ft_transcendence/backend/internal/models"
)

type userRepository struct {
	db *gorm.DB
}

// UserRepository handles the users in the database.
type UserRepository interface {
	// GetAll returns all the users.
	GetAll() ([]models.User, error)
	// GetByID returns the user with this id.
	GetByID(id string) (*models.User, error)
	// Update changes the user fields from the input and returns the updated user.
	Update(id string, input models.UpdateUserInput) (*models.User, error)
	// Delete removes the user with this id.
	Delete(id string) error
	// CreateUser saves a new user.
	CreateUser(user *models.User) error
	// GetByEmail returns the user with this email.
	GetByEmail(email string) (*models.User, error)
	// GetByUsername returns the user with this username.
	GetByUsername(username string) (*models.User, error)
	// SearchUsers returns a page of users matching the query and the total count.
	SearchUsers(query string, limit, offset int) ([]models.User, int64, error)
	// GetByIdentifier returns the user matching the email or the username.
	GetByIdentifier(identifier string) (*models.User, error)
	// GetByGithubID returns the user with this GitHub id.
	GetByGithubID(githubID string) (*models.User, error)
	// LinkGithub links a GitHub id to the user.
	LinkGithub(userID, githubID string) error
	// SetTwoFASecret saves the 2FA secret of the user.
	SetTwoFASecret(userID, secret string) error
	// SetTwoFAEnabled enables or disables 2FA for the user.
	SetTwoFAEnabled(userID string, enabled bool) error
	// ClearTwoFA removes the 2FA secret and disables 2FA for the user.
	ClearTwoFA(userID string) error
	// UpdatePassword changes the password of the user.
	UpdatePassword(userID, hashedPassword string) error
}

// NewUserRepository creates a new UserRepository using the given database.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// GetAll returns all the users.
func (r *userRepository) GetAll() ([]models.User, error) {
	var users []models.User
	result := r.db.Find(&users)
	return users, result.Error
}

// GetByID returns the user with this id.
func (r *userRepository) GetByID(id string) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, "id = ?", id)
	return &user, result.Error
}

// Update changes the user fields from the input and returns the updated user.
func (r *userRepository) Update(id string, input models.UpdateUserInput) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}

	updates := make(map[string]any)
	if input.DisplayName != "" {
		updates["display_name"] = input.DisplayName
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

// Delete scrubs the user's personal data and then soft deletes the row, so the username
// and email are freed and the account no longer exists.
func (r *userRepository) Delete(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		scrub := tx.Model(&models.User{}).Where("id = ?", id).Updates(models.AnonymizedFields(id))
		if scrub.Error != nil {
			return scrub.Error
		}
		if scrub.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		return tx.Delete(&models.User{}, "id = ?", id).Error
	})
}

// CreateUser saves a new user.
func (r *userRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

// GetByEmail returns the user with this email.
func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, "email = ?", email)
	return &user, result.Error
}

// GetByUsername returns the user with this username.
func (r *userRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, "username = ?", username)
	return &user, result.Error
}

// SearchUsers returns a page of users whose username or display name matches the
// query, ordered by username, together with the total number of matches.
func (r *userRepository) SearchUsers(query string, limit, offset int) ([]models.User, int64, error) {
	like := "%" + query + "%"
	cond := r.db.Model(&models.User{}).Where("username ILIKE ? OR display_name ILIKE ?", like, like)

	var total int64
	if err := cond.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var users []models.User
	err := cond.Order("username ASC").Limit(limit).Offset(offset).Find(&users).Error
	return users, total, err
}

// UpdatePassword changes the password of the user.
func (r *userRepository) UpdatePassword(userID, hashedPassword string) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("password", hashedPassword).Error
}

// GetByIdentifier returns the user matching the email or the username.
func (r *userRepository) GetByIdentifier(identifier string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ? OR username = ?", identifier, identifier).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByGithubID returns the user with this GitHub id.
func (r *userRepository) GetByGithubID(githubID string) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, "github_id = ?", githubID)
	return &user, result.Error
}

// LinkGithub links a GitHub id to the user.
func (r *userRepository) LinkGithub(userID, githubID string) error {
	result := r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]any{
			"github_id": githubID,
			"provider":  "github",
		})
	return result.Error
}

// SetTwoFASecret saves the 2FA secret of the user.
func (r *userRepository) SetTwoFASecret(userID, secret string) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("two_fa_secret", secret).Error
}

// SetTwoFAEnabled enables or disables 2FA for the user.
func (r *userRepository) SetTwoFAEnabled(userID string, enabled bool) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("two_fa_enabled", enabled).Error
}

// ClearTwoFA removes the 2FA secret and disables 2FA for the user.
func (r *userRepository) ClearTwoFA(userID string) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]any{
			"two_fa_secret":  nil,
			"two_fa_enabled": false,
		}).Error
}
