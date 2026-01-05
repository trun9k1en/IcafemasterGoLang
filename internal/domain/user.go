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
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username    string             `json:"username" bson:"username"`
	Email       string             `json:"email" bson:"email"`
	Password    string             `json:"-" bson:"password"` // Never expose password in JSON
	FullName    string             `json:"full_name" bson:"full_name"`
	Role        Role               `json:"role" bson:"role"`
	Permissions []Permission       `json:"permissions" bson:"permissions"`
	IsActive    bool               `json:"is_active" bson:"is_active"`
	CreatedOn   time.Time          `json:"created_on" bson:"created_on"`
	ModifiedOn  time.Time          `json:"modified_on" bson:"modified_on"`
	LastLogin   *time.Time         `json:"last_login,omitempty" bson:"last_login,omitempty"`
}

// CreateUserRequest represents request to create user
type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=100"`
	FullName string `json:"full_name" validate:"required,min=2,max=100"`
	Role     Role   `json:"role" validate:"required,oneof=admin manager staff customer"`
}

// UpdateUserRequest represents request to update user
type UpdateUserRequest struct {
	Email    string `json:"email" validate:"omitempty,email"`
	FullName string `json:"full_name" validate:"omitempty,min=2,max=100"`
	Role     Role   `json:"role" validate:"omitempty,oneof=admin manager staff customer"`
	IsActive *bool  `json:"is_active" validate:"omitempty"`
}

// ChangePasswordRequest represents request to change password
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6,max=100"`
}

// UserRepository represents the user repository contract
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
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
	ChangePassword(ctx context.Context, id string, req *ChangePasswordRequest) error
	Delete(ctx context.Context, id string) error
}

// HasPermission checks if user has a specific permission
func (u *User) HasPermission(permission Permission) bool {
	for _, p := range u.Permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// GetPermissionsForRole returns permissions for a given role
func GetPermissionsForRole(role Role) []Permission {
	if permissions, ok := RolePermissions[role]; ok {
		return permissions
	}
	return []Permission{}
}
