package domain

import "errors"

var (
	// ErrNotFound is returned when a resource is not found
	ErrNotFound = errors.New("resource not found")

	// ErrAlreadyExists is returned when a resource already exists
	ErrAlreadyExists = errors.New("resource already exists")

	// ErrInvalidID is returned when an invalid ID is provided
	ErrInvalidID = errors.New("invalid id format")

	// ErrInvalidInput is returned when invalid input is provided
	ErrInvalidInput = errors.New("invalid input")

	// ErrInternalServer is returned when an internal server error occurs
	ErrInternalServer = errors.New("internal server error")

	// ErrEmailAlreadyExists is returned when email already registered
	ErrEmailAlreadyExists = errors.New("email already registered")

	// ErrPhoneAlreadyExists is returned when phone already registered
	ErrPhoneAlreadyExists = errors.New("phone already registered")

	// ErrInvalidFileType is returned when file type is not supported
	ErrInvalidFileType = errors.New("invalid file type")

	// ErrFileTooLarge is returned when file size exceeds limit
	ErrFileTooLarge = errors.New("file size exceeds limit")
)
