package domain

import "errors"

var (
	ErrInvalidContact = errors.New("invalid contact: name and phone are required")
	ErrContactExists  = errors.New("contact already exists")
	ErrContactNotFound = errors.New("contact not found")
)
