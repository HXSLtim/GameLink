package logging

import "testing"

func TestInitReturnsLogger(t *testing.T) {
	logger := Init("debug")
	if logger == nil {
		t.Fatal("expected logger instance")
	}
}
