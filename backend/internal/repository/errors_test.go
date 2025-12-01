package repository_test

import (
	"errors"
	"testing"

	"gamelink/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestErrNotFound(t *testing.T) {
	t.Run("error should be not found", func(t *testing.T) {
		err := repository.ErrNotFound
		assert.Error(t, err)
		assert.Equal(t, "record not found", err.Error())
	})

	t.Run("error should be comparable", func(t *testing.T) {
		err1 := repository.ErrNotFound
		err2 := repository.ErrNotFound
		assert.True(t, errors.Is(err1, err2))
	})

	t.Run("error should work with errors.Is", func(t *testing.T) {
		err := repository.ErrNotFound
		assert.True(t, errors.Is(err, repository.ErrNotFound))
	})

	t.Run("error should not match other errors", func(t *testing.T) {
		err := repository.ErrNotFound
		otherErr := errors.New("other error")
		assert.False(t, errors.Is(err, otherErr))
		assert.False(t, errors.Is(otherErr, err))
	})
}

func TestErrNotFoundUsage(t *testing.T) {
	// Simulate a function that returns ErrNotFound
	simulateNotFound := func() error {
		return repository.ErrNotFound
	}

	t.Run("function returns ErrNotFound", func(t *testing.T) {
		err := simulateNotFound()
		assert.Error(t, err)
		assert.Equal(t, repository.ErrNotFound, err)
		assert.True(t, errors.Is(err, repository.ErrNotFound))
	})

	// Simulate error handling
	simulateErrorHandling := func() string {
		err := simulateNotFound()
		if errors.Is(err, repository.ErrNotFound) {
			return "not found"
		}
		return "other error"
	}

	t.Run("error handling works correctly", func(t *testing.T) {
		result := simulateErrorHandling()
		assert.Equal(t, "not found", result)
	})
}
