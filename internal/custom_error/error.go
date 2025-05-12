package custom_error

import (
	"errors"
)

var (
	ErrNotFound          = errors.New("entity not found")
	ErrIncorrectPassword = errors.New("incorrect password")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrUserAlreadyExists = errors.New("user already exists")
)
