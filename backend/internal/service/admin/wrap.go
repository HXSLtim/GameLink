package admin

import (
	"context"
	"errors"

	"gamelink/internal/apierr"
	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// WrapOrder wraps an order with error handling
func (s *AdminService) WrapOrder(ctx context.Context, orderFunc func() (*model.Order, error), operation string) (*model.Order, error) {
	order, err := orderFunc()
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apierr.NotFound("order not found").WithDetails(err.Error())
		}
		return nil, apierr.InternalError(operation + " failed").WithDetails(err.Error())
	}
	return order, nil
}

// WrapError wraps an error from repository to apierr type
func WrapError(err error, operation string) error {
	if err == nil {
		return nil
	}

	// If already an apierr, return as-is
	if _, ok := err.(*apierr.APIError); ok {
		return err
	}

	// Check for specific repository errors
	if errors.Is(err, repository.ErrNotFound) {
		return apierr.NotFound("record not found").WithDetails(err.Error())
	}

	// Check for validation errors
	if errors.Is(err, ErrValidation) {
		return apierr.BadRequest("validation failed").WithDetails(err.Error())
	}

	if errors.Is(err, ErrOrderInvalidTransition) {
		return apierr.BadRequest("invalid order status transition").WithDetails(err.Error())
	}

	// Default to internal error
	return apierr.InternalError(operation + " failed").WithDetails(err.Error())
}

// WrapVoid wraps a void function with error handling
func WrapVoid(err error, operation string) error {
	return WrapError(err, operation)
}
