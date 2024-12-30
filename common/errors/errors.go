package errors

import "fmt"

// NotFoundError represents an error for a missing entity
type NotFoundError struct {
	Entity string
	Field  string
	Value  string
}

func (e *NotFoundError) Error() string {
	if e.Entity == "" && e.Field == "" && e.Value == "" {
		return "resource not found"
	}
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

type MethodNotAllowedError struct {
	Method string
}

func (e *MethodNotAllowedError) Error() string {
	return fmt.Sprintf("%s method not allowed", e.Method)
}

func NewMethodNotAllowedError(method string) error {
	return &MethodNotAllowedError{
		Method: method,
	}
}

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

// ValidationError represents an error for invalid input
type ValidationError struct {
	Field string
	Err   string
}

func (e *ValidationError) Error() string {
	if e.Field == "" {
		return fmt.Sprintf("validation failed for %s", e.Err)
	}
	return fmt.Sprintf("validation failed for %s: %s", e.Field, e.Err)
}

// NewValidationError creates a new instance of ValidationError
func NewValidationError(field, err string) error {
	return &ValidationError{
		Field: field,
		Err:   err,
	}
}
