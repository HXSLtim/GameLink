package admin

import (
	"errors"

	"gamelink/internal/repository"
	"gamelink/pkg/apierr"
)

// wrapRepositoryError wraps repository errors into apierr types
func wrapRepositoryError(operation string, err error) error {
	if err == nil {
		return nil
	}

	// Check for specific repository errors
	if errors.Is(err, repository.ErrNotFound) {
		return apierr.NotFound("record not found").WithDetails(err.Error())
	}

	// Check for validation errors
	if errors.Is(err, ErrValidation) {
		return apierr.BadRequest("validation failed").WithDetails(err.Error())
	}

	// Default to internal error with context
	return apierr.InternalError(operation + " failed").WithDetails(err.Error())
}

// wrapError wraps any error into apierr type based on error kind
func wrapError(operation string, err error) error {
	if err == nil {
		return nil
	}

	// If already an apierr, return as-is
	if _, ok := err.(*apierr.APIError); ok {
		return err
	}

	// Check for specific errors
	if errors.Is(err, repository.ErrNotFound) {
		return apierr.NotFound("record not found").WithDetails(err.Error())
	}

	if errors.Is(err, ErrValidation) {
		return apierr.BadRequest("validation failed").WithDetails(err.Error())
	}

	if errors.Is(err, ErrOrderInvalidTransition) {
		return apierr.BadRequest("invalid order status transition").WithDetails(err.Error())
	}

	// Default to internal error
	return apierr.InternalError(operation + " failed").WithDetails(err.Error())
}

// isNotFound checks if error is a not found error
func isNotFound(err error) bool {
	if err == nil {
		return false
	}

	// Check apierr first
	if apierr.IsNotFound(err) {
		return true
	}

	// Check underlying errors
	return errors.Is(err, repository.ErrNotFound) || errors.Is(err, ErrUserNotFound)
}

// isValidationError checks if error is a validation error
func isValidationError(err error) bool {
	if err == nil {
		return false
	}

	// Check apierr first
	if apierr.IsValidationError(err) {
		return true
	}

	// Check underlying errors
	return errors.Is(err, ErrValidation)
}

// isBadRequest checks if error is a bad request error
func isBadRequest(err error) bool {
	if err == nil {
		return false
	}

	// Check apierr first
	if apierr.IsBadRequest(err) {
		return true
	}

	// Check underlying errors
	return errors.Is(err, ErrValidation) || errors.Is(err, ErrOrderInvalidTransition)
}
