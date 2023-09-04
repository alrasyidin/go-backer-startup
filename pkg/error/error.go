package customerror

import "errors"

var (
	ErrNotFound          = errors.New("data not found")
	ErrInvalidPassword   = errors.New("password is invalid")
	ErrEmailAlreadyTaken = errors.New("email has been registered")
	ErrNotOwnedCampaign  = errors.New("not an owner of the campaign")
)
