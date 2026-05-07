package controllers

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/Transcendence/services"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type UploadController struct {
	Service *services.UploadService
}

func (uc *UploadController) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file required"})
		return
	}

	if err := uc.Service.ValidateFile(file); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !isAllowedExtension(ext) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file type not allowed"})
		return
	}

	safeName := uuid.New().String() + ext
	path := filepath.Join("./uploads", safeName)

	absUploads, _ := filepath.Abs("./uploads")
	absPath, _ := filepath.Abs(path)

	if !strings.HasPrefix(absPath, absUploads) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid path"})
		return
	}

	if err := c.SaveUploadedFile(file, path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "upload failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "file uploaded",
		"path":    "/uploads/" + safeName,
	})
}

func isAllowedExtension(ext string) bool {
	allowed := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	return allowed[ext]
}
