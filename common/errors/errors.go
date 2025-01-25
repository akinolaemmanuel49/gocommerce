package errors

import "fmt"

// NotFoundError represents a not found error
type NotFoundError struct {
	Entity string
	Field  string
	Value  string
}

// Error method for the NotFoundError struct. This method defines how the
// NotFoundError error type should be represented as a string when it is converted to an error
// message.
func (e *NotFoundError) Error() string {
	if e.Entity == "" && e.Field == "" && e.Value == "" {
		return "resource not found"
	}
	return fmt.Sprintf("%s with %s: '%s' not found ", e.Entity, e.Field, e.Value)
}

// NewNotFoundError creates a new instance of NotFoundError
func NewNotFoundError(entity, field, value string) error {
	return &NotFoundError{
		Entity: entity,
		Field:  field,
		Value:  value,
	}
}

// InternalServerError represents an internal server error.
type InternalServerError struct {
}

// Error method for the InternalServerError struct. This method defines how the
// InternalServerError error type should be represented as a string when it is converted to an error
// message.
func (e *InternalServerError) Error() string {
	return "internal server error"
}

// NewInternalServerError creates a new instance of InternalServerError
func NewInternalServerError() error {
	return &InternalServerError{}
}

// MethodNotAllowedError represents method not allowed error
type MethodNotAllowedError struct {
	Method string
}

// Error method for the MethodNotAllowedError struct. This method defines how the
// MethodNotAllowedError error type should be represented as a string when it is converted to an error
// message.
func (e *MethodNotAllowedError) Error() string {
	return fmt.Sprintf("%s method not allowed", e.Method)
}

// NewMethodNotAllowedError creates a new instance of MethodNotAllowedError
func NewMethodNotAllowedError(method string) error {
	return &MethodNotAllowedError{
		Method: method,
	}
}

// ConflictError represents conflict error
type ConflictError struct {
	Entity string
	Field  string
	Value  string
}

// Error method for the ConflictError struct. This method defines how the
// ConflictError error type should be represented as a string when it is converted to an error
// message.
func (e *ConflictError) Error() string {
	return fmt.Sprintf(" %s with %s: %s already exists", e.Entity, e.Field, e.Value)
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

// Error method for the ValidationError struct. This method defines how the
// ValidationError error type should be represented as a string when it is converted to an error
// message.
func (e *ValidationError) Error() string {
	if e.Field == "" {
		return fmt.Sprint(e.Err)
	}
	return fmt.Sprintf("Validation failed for %s: %s", e.Field, e.Err)
}

// NewValidationError creates a new instance of ValidationError
func NewValidationError(field, err string) error {
	return &ValidationError{
		Field: field,
		Err:   err,
	}
}

// BadRequestError represents bad request error
type BadRequestError struct {
	Err string
}

// Error method for the BadRequestError struct. This method defines how the
// BadRequestError error type should be represented as a string when it is converted to an error
// message.
func (e *BadRequestError) Error() string {
	return fmt.Sprintf("Bad Request: %s", e.Err)
}

// NewBadRequestError creates a new instance of BadRequestError
func NewBadRequestError(err string) error {
	return &BadRequestError{
		Err: err,
	}
}

// AuthorizationError represents an error for unauthorized access
type AuthorizationError struct {
	Err string
}

// Error method for the AuthorizationError struct. This method defines how the
// AuthorizationError error type should be represented as a string when it is converted to an error
// message.
func (e *AuthorizationError) Error() string {
	if e.Err != "" {
		return "Unauthorized"
	}
	return fmt.Sprintf("Unauthorized: %s", e.Err)
}

// NewAuthorizationError creates a new instance of AuthorizationError
func NewAuthorizationError(err string) error {
	return &AuthorizationError{
		Err: err,
	}
}

// ForbiddenError represents an error for forbidden
type ForbiddenError struct {
	Err string
}

// Error method for the ForbiddenError struct. This method defines how the
// ForbiddenError error type should be represented as a string when it is converted to an error
// message.
func (e *ForbiddenError) Error() string {
	if e.Err != "" {
		return "Forbidden"
	}
	return fmt.Sprintf("Forbidden: %s", e.Err)
}

// NewForbiddenError creates a new instance of ForbiddenError
func NewForbiddenError(err string) error {
	return &ForbiddenError{
		Err: err,
	}
}
