package test

import (
	"bytes"
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func pngBytes(t *testing.T) []byte {
	t.Helper()
	img := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.NRGBA{R: 10, G: 20, B: 30, A: 255})
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("encode png: %v", err)
	}
	return buf.Bytes()
}

func uploadFile(t *testing.T, router *gin.Engine, token, filename, visibility string, content []byte) *httptest.ResponseRecorder {
	t.Helper()
	var body bytes.Buffer
	mw := multipart.NewWriter(&body)
	part, err := mw.CreateFormFile("file", filename)
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	part.Write(content)
	mw.Close()

	url := "/api/upload"
	if visibility != "" {
		url += "?visibility=" + visibility
	}
	req, _ := http.NewRequest("POST", url, &body)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func uploadAndGetID(t *testing.T, router *gin.Engine, token, visibility string) string {
	t.Helper()
	w := uploadFile(t, router, token, "pic.png", visibility, pngBytes(t))
	if w.Code != http.StatusOK {
		t.Fatalf("upload: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
	var resp struct {
		ID string `json:"id"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.ID == "" {
		t.Fatal("expected file id in upload response")
	}
	return resp.ID
}

func TestUpload_Success(t *testing.T) {
	t.Cleanup(func() { os.RemoveAll("./uploads") })
	router, _ := SetupTestEnv()
	u := registerAndLogin(t, router, "upl-ok", "upl-ok@test.com", "StrongPass123!")

	w := uploadFile(t, router, u.Token, "pic.png", "public", pngBytes(t))
	if w.Code != http.StatusOK {
		t.Fatalf("upload: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
	var resp struct {
		ID       string `json:"id"`
		URL      string `json:"url"`
		MimeType string `json:"mime_type"`
		Size     int64  `json:"size"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.ID == "" || resp.MimeType != "image/png" || resp.Size == 0 {
		t.Fatalf("unexpected upload response: %+v", resp)
	}
	if resp.URL != "/api/files/"+resp.ID {
		t.Fatalf("unexpected url: %s", resp.URL)
	}
}

func TestUpload_RequiresAuth(t *testing.T) {
	t.Cleanup(func() { os.RemoveAll("./uploads") })
	router, _ := SetupTestEnv()

	w := uploadFile(t, router, "", "pic.png", "public", pngBytes(t))
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("upload without token: expected 401, got %d", w.Code)
	}
}

func TestUpload_MissingFile(t *testing.T) {
	router, _ := SetupTestEnv()
	u := registerAndLogin(t, router, "upl-nofile", "upl-nofile@test.com", "StrongPass123!")

	var body bytes.Buffer
	mw := multipart.NewWriter(&body)
	mw.WriteField("nothing", "here")
	mw.Close()
	req, _ := http.NewRequest("POST", "/api/upload", &body)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+u.Token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("missing file: expected 400, got %d", w.Code)
	}
}

func TestUpload_InvalidExtension(t *testing.T) {
	t.Cleanup(func() { os.RemoveAll("./uploads") })
	router, _ := SetupTestEnv()
	u := registerAndLogin(t, router, "upl-ext", "upl-ext@test.com", "StrongPass123!")

	w := uploadFile(t, router, u.Token, "notes.txt", "public", []byte("hello world"))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("invalid extension: expected 400, got %d - body: %s", w.Code, w.Body.String())
	}
}

func TestUpload_MimeMismatch(t *testing.T) {
	t.Cleanup(func() { os.RemoveAll("./uploads") })
	router, _ := SetupTestEnv()
	u := registerAndLogin(t, router, "upl-mime", "upl-mime@test.com", "StrongPass123!")

	w := uploadFile(t, router, u.Token, "fake.png", "public", []byte("this is not really a png"))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("mime mismatch: expected 400, got %d - body: %s", w.Code, w.Body.String())
	}
}

func TestUpload_InvalidVisibility(t *testing.T) {
	t.Cleanup(func() { os.RemoveAll("./uploads") })
	router, _ := SetupTestEnv()
	u := registerAndLogin(t, router, "upl-vis", "upl-vis@test.com", "StrongPass123!")

	w := uploadFile(t, router, u.Token, "pic.png", "everyone", pngBytes(t))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("invalid visibility: expected 400, got %d - body: %s", w.Code, w.Body.String())
	}
}

func TestServeFile_NotFound(t *testing.T) {
	router, _ := SetupTestEnv()
	u := registerAndLogin(t, router, "srv-nf", "srv-nf@test.com", "StrongPass123!")

	w := authedRequest(t, router, "GET", "/api/files/550e8400-e29b-41d4-a716-446655440000", u.Token, "")
	if w.Code != http.StatusNotFound {
		t.Fatalf("serve missing file: expected 404, got %d", w.Code)
	}
}

func TestServeFile_PublicAccessGranted(t *testing.T) {
	t.Cleanup(func() { os.RemoveAll("./uploads") })
	router, _ := SetupTestEnv()
	owner := registerAndLogin(t, router, "srv-pub-o", "srv-pub-o@test.com", "StrongPass123!")
	other := registerAndLogin(t, router, "srv-pub-x", "srv-pub-x@test.com", "StrongPass123!")

	fileID := uploadAndGetID(t, router, owner.Token, "public")

	w := authedRequest(t, router, "GET", "/api/files/"+fileID, other.Token, "")
	if w.Code == http.StatusForbidden {
		t.Fatalf("public file should be accessible to others, got 403")
	}
}

func TestServeFile_PrivateForbiddenForOthers(t *testing.T) {
	t.Cleanup(func() { os.RemoveAll("./uploads") })
	router, _ := SetupTestEnv()
	owner := registerAndLogin(t, router, "srv-pri-o", "srv-pri-o@test.com", "StrongPass123!")
	other := registerAndLogin(t, router, "srv-pri-x", "srv-pri-x@test.com", "StrongPass123!")

	fileID := uploadAndGetID(t, router, owner.Token, "private")

	w := authedRequest(t, router, "GET", "/api/files/"+fileID, other.Token, "")
	if w.Code != http.StatusForbidden {
		t.Fatalf("private file to non-owner: expected 403, got %d", w.Code)
	}

	w = authedRequest(t, router, "GET", "/api/files/"+fileID, owner.Token, "")
	if w.Code == http.StatusForbidden {
		t.Fatalf("owner must access own private file, got 403")
	}
}

func TestServeFile_FriendsVisibility(t *testing.T) {
	t.Cleanup(func() { os.RemoveAll("./uploads") })
	router, _ := SetupTestEnv()
	owner := registerAndLogin(t, router, "srv-fr-o", "srv-fr-o@test.com", "StrongPass123!")
	stranger := registerAndLogin(t, router, "srv-fr-s", "srv-fr-s@test.com", "StrongPass123!")
	friend := registerAndLogin(t, router, "srv-fr-f", "srv-fr-f@test.com", "StrongPass123!")

	fileID := uploadAndGetID(t, router, owner.Token, "friends")

	w := authedRequest(t, router, "GET", "/api/files/"+fileID, stranger.Token, "")
	if w.Code != http.StatusForbidden {
		t.Fatalf("friends-only file to stranger: expected 403, got %d", w.Code)
	}

	authedRequest(t, router, "POST", "/api/friends/request/"+friend.ID, owner.Token, "")
	authedRequest(t, router, "POST", "/api/friends/accept/"+owner.ID, friend.Token, "")

	w = authedRequest(t, router, "GET", "/api/files/"+fileID, friend.Token, "")
	if w.Code == http.StatusForbidden {
		t.Fatalf("friends-only file to accepted friend should be granted, got 403")
	}
}

func anonRequest(t *testing.T, router *gin.Engine, path string) *httptest.ResponseRecorder {
	t.Helper()
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestServeFile_PublicAllowsAnonymous(t *testing.T) {
	t.Cleanup(func() { os.RemoveAll("./uploads") })
	router, _ := SetupTestEnv()
	owner := registerAndLogin(t, router, "srv-pub-anon-o", "srv-pub-anon-o@test.com", "StrongPass123!")
	fileID := uploadAndGetID(t, router, owner.Token, "public")

	w := anonRequest(t, router, "/api/files/"+fileID)
	if w.Code == http.StatusUnauthorized {
		t.Fatalf("public file should be reachable without auth, got 401 - body: %s", w.Body.String())
	}
}

func TestServeFile_NonPublicRejectsAnonymous(t *testing.T) {
	t.Cleanup(func() { os.RemoveAll("./uploads") })
	router, _ := SetupTestEnv()
	owner := registerAndLogin(t, router, "srv-anon-rej-o", "srv-anon-rej-o@test.com", "StrongPass123!")

	cases := map[string]string{
		"private": uploadAndGetID(t, router, owner.Token, "private"),
		"friends": uploadAndGetID(t, router, owner.Token, "friends"),
	}
	for label, id := range cases {
		w := anonRequest(t, router, "/api/files/"+id)
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("%s file to anonymous: expected 401, got %d - body: %s", label, w.Code, w.Body.String())
		}
	}
}
