package errors

import "fmt"

type ConflictError struct {
	Entity string
	Field  string
	Value  string
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("%s with %s: %s already exists", e.Entity, e.Field, e.Value)
}

// NewConflictError creates a new instance of ConflictError
func NewConflictError(entity, field, value string) error {
	return &ConflictError{
		Entity: entity,
		Field:  field,
		Value:  value,
	}
}
