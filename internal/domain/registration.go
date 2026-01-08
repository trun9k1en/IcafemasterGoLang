package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Registration represents the registration entity
type Registration struct {
	ID               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FullName         string             `json:"full_name" bson:"full_name" validate:"required,min=2,max=100"`
	PhoneNumber      string             `json:"phone_number" bson:"phone_number" validate:"required,min=10,max=15"`
	Email            string             `json:"email" bson:"email" validate:"required,email"`
	Address          string             `json:"address" bson:"address" validate:"required,min=5,max=255"`
	WorkstationRange string             `json:"workstation_range" bson:"workstation_range"`
	CreatedOn        time.Time          `json:"created_on" bson:"created_on"`
	ModifiedOn       time.Time          `json:"modified_on" bson:"modified_on"`
}

type CreateRegistrationRequest struct {
	FullName         string `json:"full_name" validate:"required,min=2,max=100"`
	PhoneNumber      string `json:"phone_number" validate:"required,min=10,max=15"`
	Email            string `json:"email" validate:"required,email"`
	Address          string `json:"address" validate:"required,min=5,max=255"`
	WorkstationRange string `json:"workstation_range" validate:"required,oneof='1-10' '10-20' '20-50' '50+'"`
}

// UpdateRegistrationRequest represents the request body for updating registration
type UpdateRegistrationRequest struct {
	FullName         string `json:"full_name" validate:"omitempty,min=2,max=100"`
	PhoneNumber      string `json:"phone_number" validate:"omitempty,min=10,max=15"`
	Email            string `json:"email" validate:"omitempty,email"`
	Address          string `json:"address" validate:"omitempty,min=5,max=255"`
	WorkstationRange string `json:"workstation_range" validate:"required,oneof='1-10' '10-20' '20-50' '50+'"`
}

// RegistrationRepository represents the registration repository contract
type RegistrationRepository interface {
	Create(ctx context.Context, registration *Registration) error
	GetByID(ctx context.Context, id string) (*Registration, error)
	GetByEmail(ctx context.Context, email string) (*Registration, error)
	GetAll(ctx context.Context, limit, offset int64) ([]*Registration, error)
	Update(ctx context.Context, id string, registration *Registration) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int64, error)
}

// RegistrationUsecase represents the registration usecase contract
type RegistrationUsecase interface {
	Create(ctx context.Context, req *CreateRegistrationRequest) (*Registration, error)
	GetByID(ctx context.Context, id string) (*Registration, error)
	GetAll(ctx context.Context, limit, offset int64) ([]*Registration, int64, error)
	Update(ctx context.Context, id string, req *UpdateRegistrationRequest) (*Registration, error)
	Delete(ctx context.Context, id string) error
}
