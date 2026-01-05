package domain

import (
	"context"
	"mime/multipart"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FileType represents the type of file
type FileType string

const (
	FileTypeDocument FileType = "document"
	FileTypeVideo    FileType = "video"
	FileTypeImage    FileType = "image"
)

// File represents the file entity
type File struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FileName    string             `json:"file_name" bson:"file_name"`
	OriginalName string            `json:"original_name" bson:"original_name"`
	FilePath    string             `json:"file_path" bson:"file_path"`
	FileType    FileType           `json:"file_type" bson:"file_type"`
	MimeType    string             `json:"mime_type" bson:"mime_type"`
	Size        int64              `json:"size" bson:"size"`
	URL         string             `json:"url" bson:"url"`
	CreatedOn   time.Time          `json:"created_on" bson:"created_on"`
}

// FileRepository represents the file repository contract
type FileRepository interface {
	Create(ctx context.Context, file *File) error
	GetByID(ctx context.Context, id string) (*File, error)
	GetAll(ctx context.Context, fileType FileType, limit, offset int64) ([]*File, error)
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context, fileType FileType) (int64, error)
}

// FileUsecase represents the file usecase contract
type FileUsecase interface {
	Upload(ctx context.Context, file *multipart.FileHeader, fileType FileType) (*File, error)
	GetByID(ctx context.Context, id string) (*File, error)
	GetAll(ctx context.Context, fileType FileType, limit, offset int64) ([]*File, int64, error)
	Delete(ctx context.Context, id string) error
}
