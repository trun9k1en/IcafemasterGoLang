package mongodb

import (
	"context"
	"time"

	"icafe-registration/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const fileCollection = "files"

type fileRepository struct {
	collection *mongo.Collection
}

// NewFileRepository creates a new file repository
func NewFileRepository(db *mongo.Database) domain.FileRepository {
	return &fileRepository{
		collection: db.Collection(fileCollection),
	}
}

// Create creates a new file record
func (r *fileRepository) Create(ctx context.Context, file *domain.File) error {
	file.ID = primitive.NewObjectID()
	file.CreatedOn = time.Now()

	_, err := r.collection.InsertOne(ctx, file)
	return err
}

// GetByID gets a file by ID
func (r *fileRepository) GetByID(ctx context.Context, id string) (*domain.File, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrInvalidID
	}

	var file domain.File
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&file)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &file, nil
}

// GetAll gets all files with pagination and optional type filter
func (r *fileRepository) GetAll(ctx context.Context, fileType domain.FileType, limit, offset int64) ([]*domain.File, error) {
	filter := bson.M{}
	if fileType != "" {
		filter["file_type"] = fileType
	}

	opts := options.Find().
		SetLimit(limit).
		SetSkip(offset).
		SetSort(bson.D{{Key: "created_on", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var files []*domain.File
	if err := cursor.All(ctx, &files); err != nil {
		return nil, err
	}

	return files, nil
}

// Delete deletes a file record
func (r *fileRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrInvalidID
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// Count counts files with optional type filter
func (r *fileRepository) Count(ctx context.Context, fileType domain.FileType) (int64, error) {
	filter := bson.M{}
	if fileType != "" {
		filter["file_type"] = fileType
	}
	return r.collection.CountDocuments(ctx, filter)
}
