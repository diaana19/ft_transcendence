package services

import (
    "errors"
    "mime/multipart"
)

type UploadService struct{}

func (s *UploadService) ValidateFile(file *multipart.FileHeader) error {
    if file.Size > 5*1024*1024 {
        return errors.New("file too large")
    }

    return nil
}
