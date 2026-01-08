package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Customer represents the customer entity
type Customer struct {
	ID               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FullName         string             `json:"full_name" bson:"full_name"`
	PhoneNumber      string             `json:"phone_number" bson:"phone_number"`
	Email            string             `json:"email,omitempty" bson:"email,omitempty"`
	Address          string             `json:"address,omitempty" bson:"address,omitempty"`
	Note             string             `json:"note,omitempty" bson:"note,omitempty"`
	WorkstationRange string             `json:"workstation_range" bson:"workstation_range"`
	IsActive         bool               `json:"is_active" bson:"is_active"`
	CreatedOn        time.Time          `json:"created_on" bson:"created_on"`
	ModifiedOn       time.Time          `json:"modified_on" bson:"modified_on"`
}

// CreateCustomerRequest represents the request body for creating customer
type CreateCustomerRequest struct {
	FullName         string `json:"full_name" validate:"required,min=2,max=100"`
	PhoneNumber      string `json:"phone_number" validate:"required,min=10,max=15"`
	Email            string `json:"email" validate:"omitempty,email"`
	Address          string `json:"address" validate:"omitempty,max=255"`
	Note             string `json:"note" validate:"omitempty,max=500"`
	WorkstationRange string `json:"workstation_range" validate:"required,oneof=1-10 10-20 20-50 50+"`
}

// UpdateCustomerRequest represents the request body for updating customer
type UpdateCustomerRequest struct {
	FullName         string `json:"full_name" validate:"omitempty,min=2,max=100"`
	PhoneNumber      string `json:"phone_number" validate:"omitempty,min=10,max=15"`
	Email            string `json:"email" validate:"omitempty,email"`
	Address          string `json:"address" validate:"omitempty,max=255"`
	Note             string `json:"note" validate:"omitempty,max=500"`
	WorkstationRange string `json:"workstation_range" validate:"omitempty,oneof=1-10 10-20 20-50 50+"`
	IsActive         *bool  `json:"is_active" validate:"omitempty"`
}

// CustomerRepository represents the customer repository contract
type CustomerRepository interface {
	Create(ctx context.Context, customer *Customer) error
	GetByID(ctx context.Context, id string) (*Customer, error)
	GetByPhone(ctx context.Context, phone string) (*Customer, error)
	GetByEmail(ctx context.Context, email string) (*Customer, error)
	GetAll(ctx context.Context, limit, offset int64) ([]*Customer, error)
	Update(ctx context.Context, id string, customer *Customer) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int64, error)
}

// CustomerUsecase represents the customer usecase contract
type CustomerUsecase interface {
	Create(ctx context.Context, req *CreateCustomerRequest) (*Customer, error)
	GetByID(ctx context.Context, id string) (*Customer, error)
	GetAll(ctx context.Context, limit, offset int64) ([]*Customer, int64, error)
	Update(ctx context.Context, id string, req *UpdateCustomerRequest) (*Customer, error)
	Delete(ctx context.Context, id string) error
}
