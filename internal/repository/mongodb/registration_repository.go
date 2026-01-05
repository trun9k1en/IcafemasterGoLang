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

const registrationCollection = "registrations"

type registrationRepository struct {
	collection *mongo.Collection
}

// NewRegistrationRepository creates a new registration repository
func NewRegistrationRepository(db *mongo.Database) domain.RegistrationRepository {
	return &registrationRepository{
		collection: db.Collection(registrationCollection),
	}
}

// Create creates a new registration
func (r *registrationRepository) Create(ctx context.Context, registration *domain.Registration) error {
	registration.ID = primitive.NewObjectID()
	registration.CreatedOn = time.Now()
	registration.ModifiedOn = time.Now()

	_, err := r.collection.InsertOne(ctx, registration)
	return err
}

// GetByID gets a registration by ID
func (r *registrationRepository) GetByID(ctx context.Context, id string) (*domain.Registration, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrInvalidID
	}

	var registration domain.Registration
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&registration)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &registration, nil
}

// GetByEmail gets a registration by email
func (r *registrationRepository) GetByEmail(ctx context.Context, email string) (*domain.Registration, error) {
	var registration domain.Registration
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&registration)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &registration, nil
}

// GetAll gets all registrations with pagination
func (r *registrationRepository) GetAll(ctx context.Context, limit, offset int64) ([]*domain.Registration, error) {
	opts := options.Find().
		SetLimit(limit).
		SetSkip(offset).
		SetSort(bson.D{{Key: "created_on", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var registrations []*domain.Registration
	if err := cursor.All(ctx, &registrations); err != nil {
		return nil, err
	}

	return registrations, nil
}

// Update updates a registration
func (r *registrationRepository) Update(ctx context.Context, id string, registration *domain.Registration) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrInvalidID
	}

	registration.ModifiedOn = time.Now()

	update := bson.M{
		"$set": bson.M{
			"full_name":       registration.FullName,
			"phone_number":    registration.PhoneNumber,
			"email":           registration.Email,
			"address":         registration.Address,
			"workstation_num": registration.WorkstationNum,
			"modified_on":     registration.ModifiedOn,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// Delete deletes a registration
func (r *registrationRepository) Delete(ctx context.Context, id string) error {
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

// Count counts all registrations
func (r *registrationRepository) Count(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{})
}
