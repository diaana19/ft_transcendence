package controllers

import (
	"net/http"

	"github.com/Transcendence/services"

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

    path := "./uploads/" + file.Filename

    if err := c.SaveUploadedFile(file, path); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "upload failed"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "file uploaded",
        "path":    path,
    })
}
