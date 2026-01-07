package main

import (
	"icafe-registration/internal/config"
	httpDelivery "icafe-registration/internal/delivery/http"
	"icafe-registration/internal/domain"
)

// =============================================================================
// Application Model
// =============================================================================

// App holds all application dependencies
type App struct {
	Config   *config.Config
	Database *DatabaseDeps
	Repos    *RepositoryDeps
	Usecases *UsecaseDeps
	Router   *httpDelivery.Router
}

// =============================================================================
// Dependency Models
// =============================================================================

// DatabaseDeps holds database connections
type DatabaseDeps struct {
	MongoDB *config.MongoDB
}

// RepositoryDeps holds all repositories
type RepositoryDeps struct {
	Registration domain.RegistrationRepository
	File         domain.FileRepository
	User         domain.UserRepository
	Customer     domain.CustomerRepository
}

// UsecaseDeps holds all usecases
type UsecaseDeps struct {
	Registration domain.RegistrationUsecase
	File         domain.FileUsecase
	Auth         domain.AuthUsecase
	User         domain.UserUsecase
	Customer     domain.CustomerUsecase
}

// =============================================================================
// User Models
// =============================================================================

// DefaultUser represents a default user configuration
type DefaultUser struct {
	Username string
	Email    string
	Phone    string
	Password string
	FullName string
	Role     domain.Role
}
