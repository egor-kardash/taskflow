package domain

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound   = errors.New("user not found")
	ErrPhoneTaken = errors.New("phone number already taken")
)

type ValidationError struct {
	Fields map[string]string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("user validation failed: %d field(s)", len(e.Fields))
}
