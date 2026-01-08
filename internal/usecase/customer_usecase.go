package usecase

import (
	"context"
	"time"

	"icafe-registration/internal/domain"
)

type customerUsecase struct {
	customerRepo   domain.CustomerRepository
	contextTimeout time.Duration
}

// NewCustomerUsecase creates a new customer usecase
func NewCustomerUsecase(repo domain.CustomerRepository, timeout time.Duration) domain.CustomerUsecase {
	return &customerUsecase{
		customerRepo:   repo,
		contextTimeout: timeout,
	}
}

// Create creates a new customer
func (u *customerUsecase) Create(ctx context.Context, req *domain.CreateCustomerRequest) (*domain.Customer, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// Check if phone already exists
	existingCustomer, err := u.customerRepo.GetByPhone(ctx, req.PhoneNumber)
	if err != nil && err != domain.ErrNotFound {
		return nil, err
	}
	if existingCustomer != nil {
		return nil, domain.ErrPhoneAlreadyExists
	}

	customer := &domain.Customer{
		FullName:         req.FullName,
		PhoneNumber:      req.PhoneNumber,
		Email:            req.Email,
		Address:          req.Address,
		Note:             req.Note,
		WorkstationRange: req.WorkstationRange,
		IsActive:         true,
		CreatedOn:        time.Now(),
		ModifiedOn:       time.Now(),
	}

	if err := u.customerRepo.Create(ctx, customer); err != nil {
		return nil, err
	}

	return customer, nil
}

// GetByID gets a customer by ID
func (u *customerUsecase) GetByID(ctx context.Context, id string) (*domain.Customer, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.customerRepo.GetByID(ctx, id)
}

// GetAll gets all customers with pagination
func (u *customerUsecase) GetAll(ctx context.Context, limit, offset int64) ([]*domain.Customer, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	customers, err := u.customerRepo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := u.customerRepo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return customers, total, nil
}

// Update updates a customer
func (u *customerUsecase) Update(ctx context.Context, id string, req *domain.UpdateCustomerRequest) (*domain.Customer, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// Get existing customer
	existing, err := u.customerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check if new phone already exists (if phone is being changed)
	if req.PhoneNumber != "" && req.PhoneNumber != existing.PhoneNumber {
		existingByPhone, err := u.customerRepo.GetByPhone(ctx, req.PhoneNumber)
		if err != nil && err != domain.ErrNotFound {
			return nil, err
		}
		if existingByPhone != nil {
			return nil, domain.ErrPhoneAlreadyExists
		}
		existing.PhoneNumber = req.PhoneNumber
	}

	// Update fields if provided
	if req.FullName != "" {
		existing.FullName = req.FullName
	}
	if req.Email != "" {
		existing.Email = req.Email
	}
	if req.Address != "" {
		existing.Address = req.Address
	}
	if req.Note != "" {
		existing.Note = req.Note
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}

	if err := u.customerRepo.Update(ctx, id, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

// Delete deletes a customer
func (u *customerUsecase) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.customerRepo.Delete(ctx, id)
}
