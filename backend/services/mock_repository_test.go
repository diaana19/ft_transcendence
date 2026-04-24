package services

import (
	"errors"

	"github.com/Transcendence/models"
)

type mockUserRepository struct {
	users  map[string]*models.User
	err    error
	nextID int
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
		return nil, errors.New("record not found")
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
	if input.Name != "" {
		user.Name = input.Name
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
	// Avatar and Wallpaper are now *string to allow null values.
	// A non-nil pointer means the client wants to set/update the field.
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

// GetByGithubID added to satisfy the UserRepository interface after OAuth work.
func (m *mockUserRepository) GetByGithubID(githubID string) (*models.User, error) {
	for _, u := range m.users {
		if u.GithubID != nil && *u.GithubID == githubID {
			return u, nil
		}
	}
	return nil, errors.New("record not found")
}
