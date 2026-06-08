package services

import (
	"errors"
	"testing"
)

func TestGetUsers_RepoError(t *testing.T) {
	repo := newMockRepo()
	repo.err = errors.New("db error")
	svc := NewUserService(repo)

	if _, err := svc.GetUsers(); err == nil {
		t.Fatal("should propagate repository error")
	}
}
