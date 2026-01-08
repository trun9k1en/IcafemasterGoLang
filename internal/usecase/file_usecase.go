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
)

type fileUsecase struct {
	fileRepo       domain.FileRepository
	uploadConfig   *config.UploadConfig
	contextTimeout time.Duration
}

// NewFileUsecase creates a new file usecase
func NewFileUsecase(
	repo domain.FileRepository,
	uploadConfig *config.UploadConfig,
	timeout time.Duration,
) domain.FileUsecase {
	return &fileUsecase{
		fileRepo:       repo,
		uploadConfig:   uploadConfig,
		contextTimeout: timeout,
	}
}

// Upload uploads a file (DOCUMENT / VIDEO) với TÊN GỐC, KHÔNG UUID
func (u *fileUsecase) Upload(
	ctx context.Context,
	fileHeader *multipart.FileHeader,
	fileType domain.FileType,
) (*domain.File, error) {

	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// Validate file size
	if fileHeader.Size > u.uploadConfig.MaxFileSize {
		return nil, domain.ErrFileTooLarge
	}

	// Validate content type
	contentType := fileHeader.Header.Get("Content-Type")
	if !u.isAllowedType(contentType) {
		return nil, domain.ErrInvalidFileType
	}

	// Open source file
	src, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// Sử dụng tên gốc
	fileName := fileHeader.Filename

	// Determine sub directory
	subDir := "files"
	if fileType == domain.FileTypeVideo {
		subDir = "videos"
	}

	// ===== ABSOLUTE PATH (ghi file ra disk) =====
	uploadDir := filepath.Join(u.uploadConfig.Path, subDir)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, err
	}

	diskPath := filepath.Join(uploadDir, fileName)

	// Nếu file đã tồn tại, ghi đè
	dst, err := os.Create(diskPath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return nil, err
	}

	// ===== RELATIVE PATH (lưu DB) =====
	dbPath := filepath.Join(subDir, fileName)

	// Build file URL
	fileURL := fmt.Sprintf(
		"%s/%s/serve/%s",
		strings.TrimRight(u.uploadConfig.BaseURL, "/"),
		subDir,
		fileName,
	)

	// Create domain file
	file := &domain.File{
		FileName:     fileName,
		OriginalName: fileHeader.Filename,
		FilePath:     dbPath, // ✅ chỉ lưu relative path
		FileType:     fileType,
		MimeType:     contentType,
		Size:         fileHeader.Size,
		URL:          fileURL,
	}

	// Lưu vào DB
	if err := u.fileRepo.Create(ctx, file); err != nil {
		_ = os.Remove(diskPath)
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
func (u *fileUsecase) GetAll(
	ctx context.Context,
	fileType domain.FileType,
	limit, offset int64,
) ([]*domain.File, int64, error) {

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

// Delete deletes file (DB + physical file)
func (u *fileUsecase) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	file, err := u.fileRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := u.fileRepo.Delete(ctx, id); err != nil {
		return err
	}

	// build absolute path before delete
	absPath := filepath.Join(u.uploadConfig.Path, file.FilePath)
	_ = os.Remove(absPath)

	return nil
}

// isAllowedType checks if content type is allowed
func (u *fileUsecase) isAllowedType(contentType string) bool {
	for _, allowed := range u.uploadConfig.AllowedTypes {
		if strings.EqualFold(contentType, allowed) {
			return true
		}
	}
	return false
}
