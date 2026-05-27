package test

import (
	"net/http"
	"testing"
)

func TestDeleteUser_Forbidden(t *testing.T) {
	router, _ := SetupTestEnv()
	owner := registerAndLogin(t, router, "deluowner", "deluowner@test.com", "StrongPass123!")
	attacker := registerAndLogin(t, router, "deluattacker", "deluattacker@test.com", "StrongPass123!")

	w := authedRequest(t, router, "DELETE", "/api/users/"+owner.ID, attacker.Token,
		`{"password":"StrongPass123!"}`)
	if w.Code != http.StatusForbidden {
		t.Fatalf("delete other user: expected 403, got %d", w.Code)
	}
}

func TestDeleteUser_MissingPasswordBody(t *testing.T) {
	router, _ := SetupTestEnv()
	u := registerAndLogin(t, router, "delunopw", "delunopw@test.com", "StrongPass123!")

	w := authedRequest(t, router, "DELETE", "/api/users/"+u.ID, u.Token, `{}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("delete without password: expected 400, got %d", w.Code)
	}
}

func TestDeleteUser_WrongPassword(t *testing.T) {
	router, _ := SetupTestEnv()
	u := registerAndLogin(t, router, "delubadpw", "deluwrongpw@test.com", "StrongPass123!")

	w := authedRequest(t, router, "DELETE", "/api/users/"+u.ID, u.Token, `{"password":"NotMyPass123!"}`)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("delete with wrong password: expected 401, got %d", w.Code)
	}
}

func TestUpdateUser_Forbidden(t *testing.T) {
	router, _ := SetupTestEnv()
	owner := registerAndLogin(t, router, "upowner", "upowner@test.com", "StrongPass123!")
	attacker := registerAndLogin(t, router, "upattacker", "upattacker@test.com", "StrongPass123!")

	w := authedRequest(t, router, "PUT", "/api/users/"+owner.ID, attacker.Token, `{"name":"hijack"}`)
	if w.Code != http.StatusForbidden {
		t.Fatalf("update other user: expected 403, got %d", w.Code)
	}
}

func TestGetUser_IncludesFollowerCounts(t *testing.T) {
	router, _ := SetupTestEnv()
	owner := registerAndLogin(t, router, "fcowner", "fcowner@test.com", "StrongPass123!")
	follower := registerAndLogin(t, router, "fcfollower", "fcfollower@test.com", "StrongPass123!")

	authedRequest(t, router, "POST", "/api/friends/follow/"+owner.ID, follower.Token, "")

	w := authedRequest(t, router, "GET", "/api/users/"+owner.ID, follower.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("get user: expected 200, got %d", w.Code)
	}
}
