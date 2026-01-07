package usecase

import (
	"context"
	"fmt"
	"time"

	"icafe-registration/internal/domain"
)

type registrationUsecase struct {
	registrationRepo domain.RegistrationRepository
	customerRepo     domain.CustomerRepository
	contextTimeout   time.Duration
}

func NewRegistrationUsecase(
	repo domain.RegistrationRepository,
	custRepo domain.CustomerRepository,
	timeout time.Duration,
) domain.RegistrationUsecase {
	return &registrationUsecase{
		registrationRepo: repo,
		customerRepo:     custRepo,
		contextTimeout:   timeout,
	}
}

func (u *registrationUsecase) Create(ctx context.Context, req *domain.CreateRegistrationRequest) (*domain.Registration, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// 1. Kiểm tra Email trong bảng Customers
	existingEmail, err := u.customerRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingEmail != nil {
		return nil, fmt.Errorf("EMAIL_ALREADY_EXISTS")
	}

	// 2. Kiểm tra Số điện thoại trong bảng Customers
	existingPhone, err := u.customerRepo.GetByPhone(ctx, req.PhoneNumber)
	if err == nil && existingPhone != nil {
		return nil, fmt.Errorf("PHONE_ALREADY_EXISTS")
	}

	// 3. Nếu không trùng, tiến hành tạo mới Customer
	customer := &domain.Customer{
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		Address:     req.Address,
		Note:        fmt.Sprintf("Web - Máy: %d", req.WorkstationNum),
		IsActive:    true,
	}

	if err := u.customerRepo.Create(ctx, customer); err != nil {
		return nil, err
	}

	// 4. Lưu log registration
	registration := &domain.Registration{
		FullName:       req.FullName,
		PhoneNumber:    req.PhoneNumber,
		Email:          req.Email,
		Address:        req.Address,
		WorkstationNum: req.WorkstationNum,
	}
	_ = u.registrationRepo.Create(ctx, registration)

	return registration, nil
}

// Các hàm bổ trợ khác giữ nguyên như bản bạn đã viết...
func (u *registrationUsecase) GetByID(ctx context.Context, id string) (*domain.Registration, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.registrationRepo.GetByID(ctx, id)
}

func (u *registrationUsecase) GetAll(ctx context.Context, limit, offset int64) ([]*domain.Registration, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	regs, _ := u.registrationRepo.GetAll(ctx, limit, offset)
	total, _ := u.registrationRepo.Count(ctx)
	return regs, total, nil
}

func (u *registrationUsecase) Update(ctx context.Context, id string, req *domain.UpdateRegistrationRequest) (*domain.Registration, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	existing, err := u.registrationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.FullName != "" {
		existing.FullName = req.FullName
	}
	if req.PhoneNumber != "" {
		existing.PhoneNumber = req.PhoneNumber
	}

	err = u.registrationRepo.Update(ctx, id, existing)
	return existing, err
}

func (u *registrationUsecase) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.registrationRepo.Delete(ctx, id)
}
