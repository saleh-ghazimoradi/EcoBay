package customErr

import "errors"

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrValidate        = errors.New("email, password(8 characters long), and phone are required")
	ErrEmailExists     = errors.New("email already exists")
	ErrInvalidInput    = errors.New("please provide a valid input")
	ErrInvalidPassword = errors.New("password does not match")
	ErrValidateCode    = errors.New("code must be 6 characters long")
	ErrInternal        = errors.New("internal server error")
	ErrUnauthorized    = errors.New("unauthorized")
)

var (
	ErrsCreate  = errors.New("failed to create")
	ErrNotFound = errors.New("not found")
	ErrUpdate   = errors.New("failed to update")
	ErrDelete   = errors.New("failed to delete")
)
