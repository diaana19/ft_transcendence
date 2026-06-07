package services

import (
	"errors"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"ft_transcendence/backend/internal/utils"

	"ft_transcendence/backend/internal/models"
	"ft_transcendence/backend/internal/repositories"
)

const (
	MaxFileSize = 25 * 1024 * 1024
)

// Only image types are accepted — videos are intentionally disallowed.
var allowedMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
}

type UploadService struct {
	FileRepo repositories.FileRepository
}

func NewUploadService(fileRepo repositories.FileRepository) *UploadService {
	return &UploadService{FileRepo: fileRepo}
}

func (s *UploadService) ValidateFile(file *multipart.FileHeader) (string, error) {
	if file.Size > MaxFileSize {
		return "", errors.New("file too large (max 25MB)")
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExtensions[ext] {
		return "", errors.New("file extension not allowed")
	}

	src, err := file.Open()
	if err != nil {
		return "", errors.New("could not open file")
	}
	defer func() { _ = src.Close() }()

	buffer := make([]byte, 512)
	if _, err := src.Read(buffer); err != nil {
		return "", errors.New("could not read file")
	}

	mime := http.DetectContentType(buffer)
	if !allowedMimeTypes[mime] {
		return "", errors.New("file type not allowed: " + mime)
	}

	return mime, nil
}

func (s *UploadService) SaveFile(
	fileHeader *multipart.FileHeader,
	ownerID string,
	visibility string,
	saveFn func(*multipart.FileHeader, string) error,
) (*models.File, error) {
	if visibility != models.FileVisibilityPublic &&
		visibility != models.FileVisibilityFriends &&
		visibility != models.FileVisibilityPrivate {
		return nil, errors.New("invalid visibility")
	}

	mime, err := s.ValidateFile(fileHeader)
	if err != nil {
		return nil, err
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	fileID := utils.NewID()
	safeName := fileID + ext
	diskPath := filepath.Join("./uploads", safeName)

	absUploads, _ := filepath.Abs("./uploads")
	absPath, _ := filepath.Abs(diskPath)
	if !strings.HasPrefix(absPath, absUploads) {
		return nil, errors.New("invalid path")
	}

	if err := saveFn(fileHeader, diskPath); err != nil {
		return nil, errors.New("failed to save file")
	}

	file := &models.File{
		ID:         fileID,
		OwnerID:    ownerID,
		Path:       diskPath,
		Filename:   fileHeader.Filename,
		MimeType:   mime,
		Size:       fileHeader.Size,
		Visibility: visibility,
	}
	if err := s.FileRepo.Create(file); err != nil {
		_ = os.Remove(diskPath)
		return nil, errors.New("failed to track file in database")
	}

	return file, nil
}
