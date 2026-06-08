package test

import (
	"testing"

	"ft_transcendence/backend/internal/models"
	"ft_transcendence/backend/internal/repositories"
	"ft_transcendence/backend/internal/utils"
)

func TestFileRepository_DeleteByOwner(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewFileRepository(db)

	userID := utils.NewID()
	user := models.User{ID: userID, Username: "fileowner", Email: "file@test.com"}
	db.Create(&user)

	fileID := utils.NewID()
	file := models.File{ID: fileID, OwnerID: userID, Path: "/uploads/file1.png", Filename: "file1.png", MimeType: "image/png", Size: 100, Visibility: "public"}
	repo.Create(&file)

	err := repo.DeleteByOwner(fileID, userID)
	if err != nil {
		t.Fatalf("DeleteByOwner: %v", err)
	}

	_, err = repo.GetByID(fileID)
	if err == nil {
		t.Fatal("expected error for deleted file")
	}
}

func TestFileRepository_DeleteByOwnerNotFound(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewFileRepository(db)

	err := repo.DeleteByOwner(utils.NewID(), utils.NewID())
	if err == nil {
		t.Fatal("expected error for non-existent file")
	}
}

func TestFileRepository_GrantAccess(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewFileRepository(db)

	userID := utils.NewID()
	fileID := utils.NewID()
	user := models.User{ID: userID, Username: "accessuser", Email: "access@test.com"}
	db.Create(&user)

	file := models.File{ID: fileID, OwnerID: utils.NewID(), Path: "/uploads/test.png", Filename: "test.png", MimeType: "image/png", Size: 100, Visibility: "private"}
	repo.Create(&file)

	err := repo.GrantAccess(fileID, userID)
	if err != nil {
		t.Fatalf("GrantAccess: %v", err)
	}

	hasAccess, err := repo.HasAccess(fileID, userID)
	if err != nil {
		t.Fatalf("HasAccess: %v", err)
	}
	if !hasAccess {
		t.Fatal("expected user to have access")
	}
}

func TestFileRepository_HasAccessNoAccess(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewFileRepository(db)

	fileID := utils.NewID()
	file := models.File{ID: fileID, OwnerID: utils.NewID(), Path: "/uploads/test.png", Filename: "test.png", MimeType: "image/png", Size: 100, Visibility: "private"}
	repo.Create(&file)

	hasAccess, err := repo.HasAccess(fileID, utils.NewID())
	if err != nil {
		t.Fatalf("HasAccess: %v", err)
	}
	if hasAccess {
		t.Fatal("expected user to not have access")
	}
}

func TestFileRepository_HasAccessGranted(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewFileRepository(db)

	userID := utils.NewID()
	fileID := utils.NewID()
	user := models.User{ID: userID, Username: "grantuser", Email: "grant@test.com"}
	file := models.File{ID: fileID, OwnerID: utils.NewID(), Path: "/uploads/test.png", Filename: "test.png", MimeType: "image/png", Size: 100, Visibility: "private"}
	db.Create(&user)
	repo.Create(&file)

	err := repo.GrantAccess(fileID, userID)
	if err != nil {
		t.Fatalf("GrantAccess: %v", err)
	}

	hasAccess, err := repo.HasAccess(fileID, userID)
	if err != nil {
		t.Fatalf("HasAccess: %v", err)
	}
	if !hasAccess {
		t.Fatal("expected user to have access after grant")
	}
}
