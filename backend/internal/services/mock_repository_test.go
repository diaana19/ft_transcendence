package services

import (
	"errors"

	"gorm.io/gorm"

	"ft_transcendence/backend/internal/models"
)

type mockUserRepository struct {
	users map[string]*models.User
	err   error
}

func newMockRepo() *mockUserRepository {
	return &mockUserRepository{
		users: make(map[string]*models.User),
	}
}

func (m *mockUserRepository) GetAll() ([]models.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	var result []models.User
	for _, u := range m.users {
		result = append(result, *u)
	}
	return result, nil
}

func (m *mockUserRepository) GetByID(id string) (*models.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	user, ok := m.users[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return user, nil
}

func (m *mockUserRepository) Update(id string, input models.UpdateUserInput) (*models.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	user, ok := m.users[id]
	if !ok {
		return nil, errors.New("record not found")
	}
	if input.DisplayName != "" {
		user.DisplayName = input.DisplayName
	}
	if input.Username != "" {
		user.Username = input.Username
	}
	if input.Email != "" {
		user.Email = input.Email
	}
	if input.Bio != "" {
		user.Bio = input.Bio
	}

	if input.Avatar != nil {
		user.Avatar = input.Avatar
	}
	if input.Wallpaper != nil {
		user.Wallpaper = input.Wallpaper
	}
	return user, nil
}

func (m *mockUserRepository) Delete(id string) error {
	if m.err != nil {
		return m.err
	}
	if _, ok := m.users[id]; !ok {
		return errors.New("record not found")
	}
	delete(m.users, id)
	return nil
}

func (m *mockUserRepository) CreateUser(user *models.User) error {
	if m.err != nil {
		return m.err
	}
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepository) GetByEmail(email string) (*models.User, error) {
	for _, u := range m.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, errors.New("record not found")
}

func (m *mockUserRepository) GetByUsername(username string) (*models.User, error) {
	for _, u := range m.users {
		if u.Username == username {
			return u, nil
		}
	}
	return nil, errors.New("record not found")
}

func (m *mockUserRepository) GetByIdentifier(identifier string) (*models.User, error) {
	for _, u := range m.users {
		if u.Email == identifier || u.Username == identifier {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *mockUserRepository) GetByGithubID(githubID string) (*models.User, error) {
	for _, u := range m.users {
		if u.GithubID != nil && *u.GithubID == githubID {
			return u, nil
		}
	}
	return nil, errors.New("record not found")
}

func (m *mockUserRepository) LinkGithub(userID, githubID string) error {
	if m.err != nil {
		return m.err
	}
	user, ok := m.users[userID]
	if !ok {
		return errors.New("record not found")
	}
	gid := githubID
	user.GithubID = &gid
	user.Provider = "github"
	return nil
}

func (m *mockUserRepository) SetTwoFASecret(userID, secret string) error {
	if m.err != nil {
		return m.err
	}
	user, ok := m.users[userID]
	if !ok {
		return errors.New("record not found")
	}
	s := secret
	user.TwoFASecret = &s
	return nil
}

func (m *mockUserRepository) SetTwoFAEnabled(userID string, enabled bool) error {
	if m.err != nil {
		return m.err
	}
	user, ok := m.users[userID]
	if !ok {
		return errors.New("record not found")
	}
	user.TwoFAEnabled = enabled
	return nil
}

func (m *mockUserRepository) ClearTwoFA(userID string) error {
	if m.err != nil {
		return m.err
	}
	user, ok := m.users[userID]
	if !ok {
		return errors.New("record not found")
	}
	user.TwoFASecret = nil
	user.TwoFAEnabled = false
	return nil
}
