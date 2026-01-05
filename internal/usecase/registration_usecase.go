package usecase

import (
	"context"
	"time"

	"icafe-registration/internal/domain"
)

type registrationUsecase struct {
	registrationRepo domain.RegistrationRepository
	contextTimeout   time.Duration
}

// NewRegistrationUsecase creates a new registration usecase
func NewRegistrationUsecase(repo domain.RegistrationRepository, timeout time.Duration) domain.RegistrationUsecase {
	return &registrationUsecase{
		registrationRepo: repo,
		contextTimeout:   timeout,
	}
}

// Create creates a new registration
func (u *registrationUsecase) Create(ctx context.Context, req *domain.CreateRegistrationRequest) (*domain.Registration, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// Check if email already exists
	existingReg, err := u.registrationRepo.GetByEmail(ctx, req.Email)
	if err != nil && err != domain.ErrNotFound {
		return nil, err
	}
	if existingReg != nil {
		return nil, domain.ErrEmailAlreadyExists
	}

	registration := &domain.Registration{
		FullName:       req.FullName,
		PhoneNumber:    req.PhoneNumber,
		Email:          req.Email,
		Address:        req.Address,
		WorkstationNum: req.WorkstationNum,
	}

	if err := u.registrationRepo.Create(ctx, registration); err != nil {
		return nil, err
	}

	return registration, nil
}

// GetByID gets a registration by ID
func (u *registrationUsecase) GetByID(ctx context.Context, id string) (*domain.Registration, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.registrationRepo.GetByID(ctx, id)
}

// GetAll gets all registrations with pagination
func (u *registrationUsecase) GetAll(ctx context.Context, limit, offset int64) ([]*domain.Registration, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	registrations, err := u.registrationRepo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := u.registrationRepo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return registrations, total, nil
}

// Update updates a registration
func (u *registrationUsecase) Update(ctx context.Context, id string, req *domain.UpdateRegistrationRequest) (*domain.Registration, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// Get existing registration
	existing, err := u.registrationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check if new email already exists (if email is being changed)
	if req.Email != "" && req.Email != existing.Email {
		existingByEmail, err := u.registrationRepo.GetByEmail(ctx, req.Email)
		if err != nil && err != domain.ErrNotFound {
			return nil, err
		}
		if existingByEmail != nil {
			return nil, domain.ErrEmailAlreadyExists
		}
		existing.Email = req.Email
	}

	// Update fields if provided
	if req.FullName != "" {
		existing.FullName = req.FullName
	}
	if req.PhoneNumber != "" {
		existing.PhoneNumber = req.PhoneNumber
	}
	if req.Address != "" {
		existing.Address = req.Address
	}
	if req.WorkstationNum > 0 {
		existing.WorkstationNum = req.WorkstationNum
	}

	if err := u.registrationRepo.Update(ctx, id, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

// Delete deletes a registration
func (u *registrationUsecase) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.registrationRepo.Delete(ctx, id)
}
