package repository

import "errors"

var (
	ErrsCreate  = errors.New("failed to create")
	ErrNotFound = errors.New("not found")
	ErrUpdate   = errors.New("failed to update")
)
