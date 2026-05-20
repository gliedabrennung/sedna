package apperr

import "errors"

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidUsername    = errors.New("username must be between 3 and 24 characters")
	ErrInvalidPassword    = errors.New("password must be between 8 and 72 characters")
)
