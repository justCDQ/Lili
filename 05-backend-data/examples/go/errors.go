package examplesgo

import (
	"errors"
	"fmt"
)

var (
	ErrOrderNotFound = errors.New("order not found")
	ErrInvalidOrder  = errors.New("invalid order")
)

type FieldError struct {
	Field string
	Value string
	Rule  string
}

func (e *FieldError) Error() string {
	return fmt.Sprintf("field %s with value %q: %s", e.Field, e.Value, e.Rule)
}

func ValidateOrder(id string, quantity int) error {
	var errs []error
	if id == "" {
		errs = append(errs, &FieldError{Field: "id", Value: id, Rule: "must not be empty"})
	}
	if quantity < 1 {
		errs = append(errs, &FieldError{Field: "quantity", Value: fmt.Sprint(quantity), Rule: "must be positive"})
	}
	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("%w: %w", ErrInvalidOrder, errors.Join(errs...))
}

func LoadOrder(id string) error {
	if id == "missing" {
		return fmt.Errorf("load order %q: %w", id, ErrOrderNotFound)
	}
	return nil
}
