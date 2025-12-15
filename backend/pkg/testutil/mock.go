// Package testutil provides utilities for testing.
package testutil

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

// MockHelper provides helper methods for working with mocks.
type MockHelper struct{}

// AssertMockCalls asserts that all expected mock calls were made.
func AssertMockCalls(t *testing.T, m *mock.Mock) {
	t.Helper()

	m.AssertExpectations(t)
}

// AssertMockNotCalled asserts that a method was not called.
func AssertMockNotCalled(t *testing.T, m *mock.Mock, methodName string) {
	t.Helper()

	for _, call := range m.Calls {
		if call.Method == methodName {
			t.Errorf("method %s should not have been called", methodName)
		}
	}
}

// ResetMocks resets all mock calls.
func ResetMocks(mocks ...*mock.Mock) {
	for _, m := range mocks {
		m.Calls = []mock.Call{}
	}
}

// MockOption is a function that configures mock behavior.
type MockOption func(*mock.Mock)

// WithReturn sets up a mock return value.
func WithReturn(method string, args ...interface{}) MockOption {
	return func(m *mock.Mock) {
		m.On(method, mock.Anything).Return(args...)
	}
}

// SetupMocks configures multiple mocks with options.
func SetupMocks(t *testing.T, opts ...MockOption) {
	t.Helper()

	for _, opt := range opts {
		m := &mock.Mock{}
		opt(m)
	}
}
