package services

import (
	"errors"
	"mime/multipart"
	"net/http"
)

type UploadService struct{}

func (s *UploadService) ValidateFile(file *multipart.FileHeader) error {
	// Size check
	if file.Size > 5*1024*1024 {
		return errors.New("file too large (max 5MB)")
	}

	// Real MIME check (not just extension)
	src, err := file.Open()
	if err != nil {
		return errors.New("could not open file")
	}
	defer src.Close()

	// Read the first 512 bytes — http.DetectContentType only needs that
	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil {
		return errors.New("could not read file")
	}

	mime := http.DetectContentType(buffer)
	allowed := []string{
		"image/jpeg",
		"image/png",
		"image/gif",
		"image/webp",
	}

	for _, m := range allowed {
		if mime == m {
			return nil
		}
	}

	return errors.New("file type not allowed: " + mime)
}
