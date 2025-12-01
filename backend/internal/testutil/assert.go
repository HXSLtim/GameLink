// Package testutil provides utilities for testing.
package testutil

import (
	"testing"
	"time"

	"gamelink/internal/model"
	"github.com/stretchr/testify/assert"
)

// AssertUserEqual asserts that two users are equal.
func AssertUserEqual(t *testing.T, expected, actual *model.User) {
	t.Helper()

	if expected == nil && actual == nil {
		return
	}

	assert.NotNil(t, expected)
	assert.NotNil(t, actual)

	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.Email, actual.Email)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Role, actual.Role)
	assert.WithinDuration(t, expected.CreatedAt, actual.CreatedAt, time.Second)
	assert.WithinDuration(t, expected.UpdatedAt, actual.UpdatedAt, time.Second)
}

// AssertPlayerEqual asserts that two players are equal.
func AssertPlayerEqual(t *testing.T, expected, actual *model.Player) {
	t.Helper()

	if expected == nil && actual == nil {
		return
	}

	assert.NotNil(t, expected)
	assert.NotNil(t, actual)

	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.UserID, actual.UserID)
	assert.Equal(t, expected.Nickname, actual.Nickname)
	assert.Equal(t, expected.VerificationStatus, actual.VerificationStatus)
	assert.WithinDuration(t, expected.CreatedAt, actual.CreatedAt, time.Second)
	assert.WithinDuration(t, expected.UpdatedAt, actual.UpdatedAt, time.Second)
}

// AssertOrderEqual asserts that two orders are equal.
func AssertOrderEqual(t *testing.T, expected, actual *model.Order) {
	t.Helper()

	if expected == nil && actual == nil {
		return
	}

	assert.NotNil(t, expected)
	assert.NotNil(t, actual)

	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.UserID, actual.UserID)
	assert.Equal(t, expected.Status, actual.Status)
	assert.Equal(t, expected.TotalPriceCents, actual.TotalPriceCents)
	assert.WithinDuration(t, expected.CreatedAt, actual.CreatedAt, time.Second)
	assert.WithinDuration(t, expected.UpdatedAt, actual.UpdatedAt, time.Second)
}

// AssertPaymentEqual asserts that two payments are equal.
func AssertPaymentEqual(t *testing.T, expected, actual *model.Payment) {
	t.Helper()

	if expected == nil && actual == nil {
		return
	}

	assert.NotNil(t, expected)
	assert.NotNil(t, actual)

	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.OrderID, actual.OrderID)
	assert.Equal(t, expected.UserID, actual.UserID)
	assert.Equal(t, expected.AmountCents, actual.AmountCents)
	assert.Equal(t, expected.Status, actual.Status)
	assert.WithinDuration(t, expected.CreatedAt, actual.CreatedAt, time.Second)
	assert.WithinDuration(t, expected.UpdatedAt, actual.UpdatedAt, time.Second)
}

// AssertError asserts that an error occurred and matches the expected message.
func AssertError(t *testing.T, err error, contains string) {
	t.Helper()

	assert.Error(t, err)
	if contains != "" {
		assert.Contains(t, err.Error(), contains)
	}
}

// AssertNoError asserts that no error occurred.
func AssertNoError(t *testing.T, err error) {
	t.Helper()

	assert.NoError(t, err)
}
