package models

import "testing"

// Repost has no HTTP route yet, so Repost.ToResponse is unreachable through the
// e2e suite. It is a pure mapping, so unit-test it directly: the embedded author
// must be projected through User.ToResponse and the scalar fields copied across.
func TestRepostToResponse(t *testing.T) {
	avatar := "https://example.com/a.png"
	r := Repost{
		ID:       "repost-1",
		PostID:   "post-1",
		AuthorID: "user-1",
		Author: User{
			ID:       "user-1",
			Username: "alice",
			Email:    "alice@example.com",
			Avatar:   &avatar,
		},
	}

	res := r.ToResponse()

	if res.ID != r.ID || res.PostID != r.PostID || res.AuthorID != r.AuthorID {
		t.Errorf("scalar fields not copied: got %+v", res)
	}
	if res.CreatedAt != r.CreatedAt {
		t.Errorf("CreatedAt = %v, want %v", res.CreatedAt, r.CreatedAt)
	}
	if res.Author.ID != "user-1" || res.Author.Username != "alice" {
		t.Errorf("author not projected via User.ToResponse: got %+v", res.Author)
	}
	if res.Author.Avatar == nil || *res.Author.Avatar != avatar {
		t.Errorf("author avatar = %v, want %q", res.Author.Avatar, avatar)
	}
}
