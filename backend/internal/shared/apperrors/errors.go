package apperrors

import "errors"

var (
	ErrNotFound            = errors.New("not found")
	ErrAlreadyExists       = errors.New("already exists")
	ErrInvalidReference    = errors.New("invalid reference")
	ErrConstraintViolation = errors.New("constraint violation")
)
