package usecase

import (
	"context"
	"time"

	"icafe-registration/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepo       domain.UserRepository
	contextTimeout time.Duration
}

// NewUserUsecase creates a new user usecase
func NewUserUsecase(repo domain.UserRepository, timeout time.Duration) domain.UserUsecase {
	return &userUsecase{
		userRepo:       repo,
		contextTimeout: timeout,
	}
}

// Create creates a new user (admin only)
func (u *userUsecase) Create(ctx context.Context, req *domain.CreateUserRequest) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// Check if username already exists
	existingUser, err := u.userRepo.GetByUsername(ctx, req.Username)
	if err != nil && err != domain.ErrNotFound {
		return nil, err
	}
	if existingUser != nil {
		return nil, domain.ErrAlreadyExists
	}

	// Check if phone already exists
	existingUser, err = u.userRepo.GetByPhone(ctx, req.Phone)
	if err != nil && err != domain.ErrNotFound {
		return nil, err
	}
	if existingUser != nil {
		return nil, domain.ErrPhoneAlreadyExists
	}

	// Check if email already exists (if provided)
	if req.Email != "" {
		existingUser, err = u.userRepo.GetByEmail(ctx, req.Email)
		if err != nil && err != domain.ErrNotFound {
			return nil, err
		}
		if existingUser != nil {
			return nil, domain.ErrEmailAlreadyExists
		}
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: string(hashedPassword),
		FullName: req.FullName,
		Role:     req.Role,
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetByID gets a user by ID
func (u *userUsecase) GetByID(ctx context.Context, id string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.userRepo.GetByID(ctx, id)
}

// GetAll gets all users with pagination
func (u *userUsecase) GetAll(ctx context.Context, limit, offset int64) ([]*domain.User, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	users, err := u.userRepo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := u.userRepo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// Update updates a user
func (u *userUsecase) Update(ctx context.Context, id string, req *domain.UpdateUserRequest) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// Get existing user
	existing, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check if new email already exists
	if req.Email != "" && req.Email != existing.Email {
		existingByEmail, err := u.userRepo.GetByEmail(ctx, req.Email)
		if err != nil && err != domain.ErrNotFound {
			return nil, err
		}
		if existingByEmail != nil {
			return nil, domain.ErrEmailAlreadyExists
		}
		existing.Email = req.Email
	}

	// Check if new phone already exists
	if req.Phone != "" && req.Phone != existing.Phone {
		existingByPhone, err := u.userRepo.GetByPhone(ctx, req.Phone)
		if err != nil && err != domain.ErrNotFound {
			return nil, err
		}
		if existingByPhone != nil {
			return nil, domain.ErrPhoneAlreadyExists
		}
		existing.Phone = req.Phone
	}

	// Update fields if provided
	if req.FullName != "" {
		existing.FullName = req.FullName
	}
	if req.Role != "" {
		existing.Role = req.Role
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}
	// Update custom permissions (admin can assign extra permissions beyond role)
	if req.CustomPermissions != nil {
		existing.CustomPermissions = req.CustomPermissions
	}

	if err := u.userRepo.Update(ctx, id, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

// UpdateRole updates user role and custom permissions (admin only)
func (u *userUsecase) UpdateRole(ctx context.Context, id string, req *domain.UpdateUserRoleRequest) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// Get existing user
	existing, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update role if provided
	if req.Role != "" {
		existing.Role = req.Role
		existing.Permissions = domain.GetPermissionsForRole(req.Role)
	}

	// Update custom permissions
	if req.CustomPermissions != nil {
		existing.CustomPermissions = req.CustomPermissions
	}

	if err := u.userRepo.Update(ctx, id, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

// ChangePassword changes user password
func (u *userUsecase) ChangePassword(ctx context.Context, id string, req *domain.ChangePasswordRequest) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// Get existing user
	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return domain.ErrInvalidCredentials
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	return u.userRepo.Update(ctx, id, user)
}

// Delete deletes a user
func (u *userUsecase) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.userRepo.Delete(ctx, id)
}
