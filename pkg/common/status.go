package common

import "errors"

var (
	ErrInternalServerError = errors.New("internal server error")
	ErrNotFound            = errors.New("user not found")
	ErrInvalidEmail        = errors.New("invalid email")
	ErrInvalidPassword     = errors.New("invalid password")
	ErrInvalidName         = errors.New("invalid name")
	ErrConflict            = errors.New("email already exists")
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrNameAlreadyExists   = errors.New("name already exists")
	ErrUserInactive        = errors.New("user is inactive")
	ErrUnauthorizedAccess  = errors.New("unauthorized access")
)
