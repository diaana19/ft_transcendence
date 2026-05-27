package controllers

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"ft_transcendence/backend/models"
	"ft_transcendence/backend/services"
)

type UploadController struct {
	Service       *services.UploadService
	FriendService *services.FriendService
}

func (uc *UploadController) UploadFile(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	ownerID := userIDRaw.(string)

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file required"})
		return
	}

	visibility := c.DefaultQuery("visibility", models.FileVisibilityPublic)

	saveFn := func(fh *multipart.FileHeader, dst string) error {
		return c.SaveUploadedFile(fh, dst)
	}

	file, err := uc.Service.SaveFile(fileHeader, ownerID, visibility, saveFn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":        file.ID,
		"url":       "/api/files/" + file.ID,
		"mime_type": file.MimeType,
		"size":      file.Size,
	})
}

func (uc *UploadController) ServeFile(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDRaw.(string)

	fileID := c.Param("id")
	file, err := uc.Service.FileRepo.GetByID(fileID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !uc.canAccess(file, userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	// Utiliser le chemin complet depuis la BD
	realPath := filepath.Join("/root", file.Path)

	fmt.Printf("[ServeFile] Trying real path: %s\n", realPath)

	if _, err := os.Stat(realPath); err == nil {
		fmt.Printf("[ServeFile] ✅ SUCCESS serving: %s\n", realPath)
		c.Header("Content-Type", file.MimeType)
		c.Header("Content-Disposition", `inline; filename="`+file.Filename+`"`)
		c.File(realPath)
		return
	}

	// Fallback si jamais
	fmt.Printf("[ServeFile] ❌ Still not found at %s\n", realPath)
	c.JSON(http.StatusNotFound, gin.H{"error": "file not found on disk"})
}

func (uc *UploadController) canAccess(file *models.File, userID string) bool {
	if file.OwnerID == userID {
		return true
	}

	switch file.Visibility {
	case models.FileVisibilityPublic:
		return true

	case models.FileVisibilityFriends:
		if uc.FriendService == nil {
			return false
		}
		isFriend, err := uc.FriendService.AreFriends(userID, file.OwnerID)
		return err == nil && isFriend

	case models.FileVisibilityPrivate:
		hasAccess, err := uc.Service.FileRepo.HasAccess(file.ID, userID)
		return err == nil && hasAccess

	default:
		return false
	}
}
