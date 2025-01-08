package domain

import "errors"

var (
	ErrInvalidContactName = errors.New("invalid contact: Name is required")
	ErrInvalidContactNumber = errors.New("invalid contact: Phone is required")
	ErrContactExists  = errors.New("contact already exists")
	ErrContactNotFound = errors.New("contact not found")
)
