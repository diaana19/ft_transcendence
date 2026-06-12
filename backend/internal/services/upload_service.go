package services

import (
	"errors"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"ft_transcendence/backend/internal/models"
	"ft_transcendence/backend/internal/repositories"
	"ft_transcendence/backend/internal/utils"
)

// MaxFileSize is the max allowed upload size in bytes (25MB).
const (
	MaxFileSize = 25 * 1024 * 1024
)

var allowedMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
	"video/mp4":  true,
	"video/webm": true,
}

var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
	".mp4":  true,
	".webm": true,
}

// UploadService validates and saves uploaded files.
type UploadService struct {
	FileRepo  repositories.FileRepository
	UploadDir string
}

// NewUploadService creates a new UploadService.
func NewUploadService(fileRepo repositories.FileRepository, uploadDir string) *UploadService {
	return &UploadService{FileRepo: fileRepo, UploadDir: uploadDir}
}

// ValidateFile checks the size, extension and mime type of the file. It returns the mime type.
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

// SaveFile validates the file, writes it to disk with saveFn and tracks it in the database.
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
	// filepath.Join cleans the path (strips "./"), so compare against the cleaned base.
	baseDir := filepath.Clean(s.UploadDir)
	dst := filepath.Join(baseDir, safeName)

	if !strings.HasPrefix(dst, baseDir+string(filepath.Separator)) {
		return nil, errors.New("invalid path")
	}

	if err := os.MkdirAll(s.UploadDir, 0o755); err != nil {
		return nil, errors.New("failed to create upload directory")
	}

	if err := saveFn(fileHeader, dst); err != nil {
		return nil, errors.New("failed to save file")
	}

	file := &models.File{
		ID:         fileID,
		OwnerID:    ownerID,
		Path:       dst,
		Filename:   fileHeader.Filename,
		MimeType:   mime,
		Size:       fileHeader.Size,
		Visibility: visibility,
	}
	if err := s.FileRepo.Create(file); err != nil {
		_ = os.Remove(dst)
		return nil, errors.New("failed to track file in database")
	}

	return file, nil
}
