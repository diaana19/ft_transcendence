package services

import (
	"errors"
	"testing"
)

// UserService list/get/update/delete success paths, not-found cases and partial
// updates are covered end-to-end through the /users endpoints. Only the
// repository error branch (unreachable over HTTP) is unit-tested here.

func TestGetUsers_RepoError(t *testing.T) {
	repo := newMockRepo()
	repo.err = errors.New("db error")
	svc := NewUserService(repo)

	if _, err := svc.GetUsers(); err == nil {
		t.Fatal("should propagate repository error")
	}
}
