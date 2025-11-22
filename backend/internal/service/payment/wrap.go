package payment

import (
	"errors"
	
	"gamelink/internal/apierr"
	"gamelink/internal/repository"
)

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
	
	// Default to internal error
	return apierr.InternalError(operation + " failed").WithDetails(err.Error())
}