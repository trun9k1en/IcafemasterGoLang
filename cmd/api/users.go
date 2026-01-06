package main

import (
	"context"
	"log"

	"icafe-registration/internal/domain"
)

// DefaultUsers list of default users to create on startup
var DefaultUsers = []DefaultUser{
	{
		Username: "admin",
		Email:    "admin@icafe.local",
		Phone:    "0900000001",
		Password: "admin123",
		FullName: "Administrator",
		Role:     domain.RoleAdmin,
	},
	{
		Username: "sale",
		Email:    "sale@icafe.local",
		Phone:    "0900000002",
		Password: "sale123",
		FullName: "Sale User",
		Role:     domain.RoleSale,
	},
}

// createDefaultUsers creates default users if they don't exist
func (a *App) createDefaultUsers() {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	for _, user := range DefaultUsers {
		a.createUserIfNotExists(ctx, user)
	}
}

// createUserIfNotExists creates a user if it doesn't already exist
func (a *App) createUserIfNotExists(ctx context.Context, defaultUser DefaultUser) {
	_, err := a.Repos.User.GetByUsername(ctx, defaultUser.Username)
	if err == nil {
		log.Printf("User '%s' already exists", defaultUser.Username)
		return
	}

	req := &domain.CreateUserRequest{
		Username: defaultUser.Username,
		Email:    defaultUser.Email,
		Phone:    defaultUser.Phone,
		Password: defaultUser.Password,
		FullName: defaultUser.FullName,
		Role:     defaultUser.Role,
	}

	_, err = a.Usecases.User.Create(ctx, req)
	if err != nil {
		log.Printf("Failed to create user '%s': %v", defaultUser.Username, err)
		return
	}

	log.Printf("Created default user: %s (password: %s)", defaultUser.Username, defaultUser.Password)
}
