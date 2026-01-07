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

const customerCollection = "customers"

type customerRepository struct {
	collection *mongo.Collection
}

// NewCustomerRepository creates a new customer repository
func NewCustomerRepository(db *mongo.Database) domain.CustomerRepository {
	collection := db.Collection(customerCollection)

	// Create unique index on phone_number
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "phone_number", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	collection.Indexes().CreateOne(context.Background(), indexModel)

	return &customerRepository{
		collection: collection,
	}
}

// Create creates a new customer
func (r *customerRepository) Create(ctx context.Context, customer *domain.Customer) error {
	customer.ID = primitive.NewObjectID()
	customer.IsActive = true
	customer.CreatedOn = time.Now()
	customer.ModifiedOn = time.Now()

	_, err := r.collection.InsertOne(ctx, customer)
	return err
}

// GetByID gets a customer by ID
func (r *customerRepository) GetByID(ctx context.Context, id string) (*domain.Customer, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrInvalidID
	}

	var customer domain.Customer
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&customer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &customer, nil
}

// GetByPhone gets a customer by phone number
func (r *customerRepository) GetByPhone(ctx context.Context, phone string) (*domain.Customer, error) {
	var customer domain.Customer
	err := r.collection.FindOne(ctx, bson.M{"phone_number": phone}).Decode(&customer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &customer, nil
}

// GetAll gets all customers with pagination
func (r *customerRepository) GetAll(ctx context.Context, limit, offset int64) ([]*domain.Customer, error) {
	opts := options.Find().
		SetLimit(limit).
		SetSkip(offset).
		SetSort(bson.D{{Key: "created_on", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var customers []*domain.Customer
	if err := cursor.All(ctx, &customers); err != nil {
		return nil, err
	}

	return customers, nil
}

// Update updates a customer
func (r *customerRepository) Update(ctx context.Context, id string, customer *domain.Customer) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrInvalidID
	}

	customer.ModifiedOn = time.Now()

	update := bson.M{
		"$set": bson.M{
			"full_name":    customer.FullName,
			"phone_number": customer.PhoneNumber,
			"email":        customer.Email,
			"address":      customer.Address,
			"note":         customer.Note,
			"is_active":    customer.IsActive,
			"modified_on":  customer.ModifiedOn,
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

// Delete deletes a customer
func (r *customerRepository) Delete(ctx context.Context, id string) error {
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

// Count counts all customers
func (r *customerRepository) Count(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{})
}
