package entity

import (
	"fmt"
	"strings"
)

type NotFoundError struct {
	Entity string
	Value  any
}

func NewNotFoundError(entity string, value any) *NotFoundError {
	return &NotFoundError{Entity: entity, Value: value}
}

func (n *NotFoundError) Error() string {
	return fmt.Sprintf("%s with value %v not found", n.Entity, n.Value)
}

type ValidationError struct {
	Field   string
	Message string
}

func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{Field: field, Message: message}
}

func (v *ValidationError) Error() string {
	return fmt.Sprintf("validation error: field '%s' - %s", v.Field, v.Message)
}

type ValidationErrors struct {
	Errors []*ValidationError
}

func WrapValidationError(errors ...*ValidationError) *ValidationErrors {
	return &ValidationErrors{Errors: errors}
}

func (v *ValidationErrors) Error() string {
	var errMessages []string
	for _, err := range v.Errors {
		errMessages = append(errMessages, err.Error())
	}
	return "validation errors: " + strings.Join(errMessages, "; ")
}

type Err struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
