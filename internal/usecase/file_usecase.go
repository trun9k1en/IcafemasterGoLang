package usecase

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"icafe-registration/internal/config"
	"icafe-registration/internal/domain"

	"github.com/google/uuid"
)

type fileUsecase struct {
	fileRepo       domain.FileRepository
	uploadConfig   *config.UploadConfig
	contextTimeout time.Duration
}

// NewFileUsecase creates a new file usecase
func NewFileUsecase(repo domain.FileRepository, uploadConfig *config.UploadConfig, timeout time.Duration) domain.FileUsecase {
	return &fileUsecase{
		fileRepo:       repo,
		uploadConfig:   uploadConfig,
		contextTimeout: timeout,
	}
}

// Upload uploads a file
func (u *fileUsecase) Upload(ctx context.Context, fileHeader *multipart.FileHeader, fileType domain.FileType) (*domain.File, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// Validate file size
	if fileHeader.Size > u.uploadConfig.MaxFileSize {
		return nil, domain.ErrFileTooLarge
	}

	// Get content type
	contentType := fileHeader.Header.Get("Content-Type")
	if !u.isAllowedType(contentType) {
		return nil, domain.ErrInvalidFileType
	}

	// Open the uploaded file
	src, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// Generate unique filename
	ext := filepath.Ext(fileHeader.Filename)
	uniqueFileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// Determine subdirectory based on file type
	subDir := "files"
	if fileType == domain.FileTypeVideo {
		subDir = "videos"
	}

	// Create directory if not exists
	uploadDir := filepath.Join(u.uploadConfig.Path, subDir)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, err
	}

	// Create destination file
	filePath := filepath.Join(uploadDir, uniqueFileName)
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, src); err != nil {
		return nil, err
	}

	// Create file URL
	fileURL := fmt.Sprintf("%s/api/v1/%s/%s", u.uploadConfig.BaseURL, subDir, uniqueFileName)

	// Create file record
	file := &domain.File{
		FileName:     uniqueFileName,
		OriginalName: fileHeader.Filename,
		FilePath:     filePath,
		FileType:     fileType,
		MimeType:     contentType,
		Size:         fileHeader.Size,
		URL:          fileURL,
	}

	if err := u.fileRepo.Create(ctx, file); err != nil {
		// Remove uploaded file if database insert fails
		os.Remove(filePath)
		return nil, err
	}

	return file, nil
}

// GetByID gets a file by ID
func (u *fileUsecase) GetByID(ctx context.Context, id string) (*domain.File, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.fileRepo.GetByID(ctx, id)
}

// GetAll gets all files with pagination
func (u *fileUsecase) GetAll(ctx context.Context, fileType domain.FileType, limit, offset int64) ([]*domain.File, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	files, err := u.fileRepo.GetAll(ctx, fileType, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := u.fileRepo.Count(ctx, fileType)
	if err != nil {
		return nil, 0, err
	}

	return files, total, nil
}

// Delete deletes a file
func (u *fileUsecase) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// Get file info first
	file, err := u.fileRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Delete from database
	if err := u.fileRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Delete physical file
	os.Remove(file.FilePath)

	return nil
}

// isAllowedType checks if the content type is allowed
func (u *fileUsecase) isAllowedType(contentType string) bool {
	for _, allowed := range u.uploadConfig.AllowedTypes {
		if strings.EqualFold(contentType, allowed) {
			return true
		}
	}
	return false
}
