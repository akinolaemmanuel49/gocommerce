package errors

import "fmt"

// ValidationError represents an error for invalid input
type ValidationError struct {
	Field string
	Err   string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed for %s: %s", e.Field, e.Err)
}

// NewValidationError creates a new instance of ValidationError
func NewValidationError(field, err string) error {
	return &ValidationError{
		Field: field,
		Err:   err,
	}
}
