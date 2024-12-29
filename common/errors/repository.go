package errors

import "fmt"

// NotFoundError represents an error for a missing entity
type NotFoundError struct {
	Entity string
	Field  string
	Value  string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with %s: %s not found ", e.Entity, e.Field, e.Value)
}

// NewNotFoundError creates a new instance of NotFoundError
func NewNotFoundError(entity, field, value string) error {
	return &NotFoundError{
		Entity: entity,
		Field:  field,
		Value:  value,
	}
}

type InternalServerError struct {
}

func (e *InternalServerError) Error() string {
	return "internal server error"
}

func NewInternalServerError() error {
	return &InternalServerError{}
}
