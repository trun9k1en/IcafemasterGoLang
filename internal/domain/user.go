package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Role represents user role
type Role string

const (
	RoleAdmin    Role = "admin"
	RoleManager  Role = "manager"
	RoleSale     Role = "sale"
	RoleStaff    Role = "staff"
	RoleCustomer Role = "customer"
)

// Permission represents user permission
type Permission string

const (
	PermissionReadRegistration   Permission = "registration:read"
	PermissionWriteRegistration  Permission = "registration:write"
	PermissionDeleteRegistration Permission = "registration:delete"
	PermissionReadFile           Permission = "file:read"
	PermissionWriteFile          Permission = "file:write"
	PermissionDeleteFile         Permission = "file:delete"
	PermissionManageUser         Permission = "user:manage"
)

// RolePermissions defines permissions for each role
var RolePermissions = map[Role][]Permission{
	RoleAdmin: {
		PermissionReadRegistration,
		PermissionWriteRegistration,
		PermissionDeleteRegistration,
		PermissionReadFile,
		PermissionWriteFile,
		PermissionDeleteFile,
		PermissionManageUser,
	},
	RoleManager: {
		PermissionReadRegistration,
		PermissionWriteRegistration,
		PermissionDeleteRegistration,
		PermissionReadFile,
		PermissionWriteFile,
		PermissionDeleteFile,
	},
	RoleSale: {
		PermissionReadRegistration,
		PermissionWriteRegistration,
		PermissionReadFile,
		PermissionWriteFile,
	},
	RoleStaff: {
		PermissionReadRegistration,
		PermissionWriteRegistration,
		PermissionReadFile,
		PermissionWriteFile,
	},
	RoleCustomer: {
		PermissionReadRegistration,
		PermissionReadFile,
	},
}

// User represents the user entity
type User struct {
	ID                primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username          string             `json:"username" bson:"username"`
	Email             string             `json:"email,omitempty" bson:"email,omitempty"`
	Phone             string             `json:"phone" bson:"phone"`
	Password          string             `json:"-" bson:"password"` // Never expose password in JSON
	FullName          string             `json:"full_name" bson:"full_name"`
	Role              Role               `json:"role" bson:"role"`
	Permissions       []Permission       `json:"permissions" bson:"permissions"`
	CustomPermissions []Permission       `json:"custom_permissions,omitempty" bson:"custom_permissions,omitempty"` // Admin-assigned custom permissions
	IsActive          bool               `json:"is_active" bson:"is_active"`
	CreatedOn         time.Time          `json:"created_on" bson:"created_on"`
	ModifiedOn        time.Time          `json:"modified_on" bson:"modified_on"`
	LastLogin         *time.Time         `json:"last_login,omitempty" bson:"last_login,omitempty"`
}

// RegisterRequest represents request to register a new user (public)
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6,max=100"`
	Phone    string `json:"phone" validate:"required,min=10,max=15"`
	FullName string `json:"full_name" validate:"required,min=2,max=100"`
}

// CreateUserRequest represents request to create user (admin only)
type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"omitempty,email"`
	Phone    string `json:"phone" validate:"required,min=10,max=15"`
	Password string `json:"password" validate:"required,min=6,max=100"`
	FullName string `json:"full_name" validate:"required,min=2,max=100"`
	Role     Role   `json:"role" validate:"required,oneof=admin manager sale staff customer"`
}

// UpdateUserRequest represents request to update user
type UpdateUserRequest struct {
	Email             string       `json:"email" validate:"omitempty,email"`
	Phone             string       `json:"phone" validate:"omitempty,min=10,max=15"`
	FullName          string       `json:"full_name" validate:"omitempty,min=2,max=100"`
	Role              Role         `json:"role" validate:"omitempty,oneof=admin manager sale staff customer"`
	IsActive          *bool        `json:"is_active" validate:"omitempty"`
	CustomPermissions []Permission `json:"custom_permissions" validate:"omitempty"` // Admin can assign custom permissions
}

// ChangePasswordRequest represents request to change password
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6,max=100"`
}

// UpdateUserRoleRequest represents request to update user role and permissions (admin only)
type UpdateUserRoleRequest struct {
	Role              Role         `json:"role" validate:"omitempty,oneof=admin manager sale staff customer"`
	CustomPermissions []Permission `json:"custom_permissions" validate:"omitempty"`
}

// UserRepository represents the user repository contract
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByPhone(ctx context.Context, phone string) (*User, error)
	GetAll(ctx context.Context, limit, offset int64) ([]*User, error)
	Update(ctx context.Context, id string, user *User) error
	UpdateLastLogin(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int64, error)
}

// UserUsecase represents the user usecase contract
type UserUsecase interface {
	Create(ctx context.Context, req *CreateUserRequest) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	GetAll(ctx context.Context, limit, offset int64) ([]*User, int64, error)
	Update(ctx context.Context, id string, req *UpdateUserRequest) (*User, error)
	UpdateRole(ctx context.Context, id string, req *UpdateUserRoleRequest) (*User, error)
	ChangePassword(ctx context.Context, id string, req *ChangePasswordRequest) error
	Delete(ctx context.Context, id string) error
}

// HasPermission checks if user has a specific permission (from role or custom)
func (u *User) HasPermission(permission Permission) bool {
	// Check role-based permissions
	for _, p := range u.Permissions {
		if p == permission {
			return true
		}
	}
	// Check custom permissions
	for _, p := range u.CustomPermissions {
		if p == permission {
			return true
		}
	}
	return false
}

// GetAllPermissions returns all permissions (role-based + custom)
func (u *User) GetAllPermissions() []Permission {
	permMap := make(map[Permission]bool)
	for _, p := range u.Permissions {
		permMap[p] = true
	}
	for _, p := range u.CustomPermissions {
		permMap[p] = true
	}
	result := make([]Permission, 0, len(permMap))
	for p := range permMap {
		result = append(result, p)
	}
	return result
}

// GetPermissionsForRole returns permissions for a given role
func GetPermissionsForRole(role Role) []Permission {
	if permissions, ok := RolePermissions[role]; ok {
		return permissions
	}
	return []Permission{}
}
