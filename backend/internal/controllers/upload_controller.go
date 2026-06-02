package controllers

import (
	"errors"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"ft_transcendence/backend/internal/models"
	"ft_transcendence/backend/internal/services"
)

type UploadController struct {
	Service       *services.UploadService
	FriendService *services.FriendService
}

// UploadFile godoc
// @Summary  Upload a file
// @Tags     uploads
// @Security BearerAuth
// @Accept   multipart/form-data
// @Produce  json
// @Param    file       formData file   true  "file to upload"
// @Param    visibility query    string false "file visibility (public, friends, private)" default(public)
// @Success  200 {object} map[string]interface{}
// @Failure  400 {object} map[string]string
// @Failure  401 {object} map[string]string
// @Router   /upload [post]
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

// ServeFile godoc
// @Summary  Serve a stored file
// @Description Public files are accessible to anyone; friends/private files require an authenticated, authorised user.
// @Tags     uploads
// @Produce  octet-stream
// @Param    id path string true "file ID"
// @Success  200 {file} binary
// @Failure  401 {object} map[string]string
// @Failure  403 {object} map[string]string
// @Failure  404 {object} map[string]string
// @Failure  500 {object} map[string]string
// @Router   /files/{id} [get]
func (uc *UploadController) ServeFile(c *gin.Context) {
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

	if _, statErr := os.Stat(file.Path); os.IsNotExist(statErr) {
		c.JSON(http.StatusNotFound, gin.H{"error": "file missing on disk"})
		return
	}

	if file.Visibility == models.FileVisibilityPublic {
		uc.streamFile(c, file)
		return
	}

	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDRaw.(string)

	if !uc.canAccess(file, userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	uc.streamFile(c, file)
}

// streamFile writes the file to the response with inline content headers.
func (uc *UploadController) streamFile(c *gin.Context, file *models.File) {
	c.Header("Content-Type", file.MimeType)
	c.Header("Content-Disposition", `inline; filename="`+file.Filename+`"`)
	c.File(file.Path)
}

// canAccess reports whether userID may read file, based on ownership and the
// file's visibility (public, friends, or private).
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
