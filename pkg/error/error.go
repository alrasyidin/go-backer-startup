package customerror

import "errors"

var (
	ErrNotFound          = errors.New("data not found")
	ErrInvalidPassword   = errors.New("password is invalid")
	ErrEmailAlreadyTaken = errors.New("email has been registered")
)
