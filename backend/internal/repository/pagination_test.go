package repository

import "testing"

func TestNormalizePage(t *testing.T) {
	t.Run("zero becomes default", func(t *testing.T) {
		if got := NormalizePage(0); got != 1 {
			t.Fatalf("NormalizePage(0)=%d, want 1", got)
		}
	})
	t.Run("negative becomes default", func(t *testing.T) {
		if got := NormalizePage(-5); got != 1 {
			t.Fatalf("NormalizePage(-5)=%d, want 1", got)
		}
	})
	t.Run("positive unchanged", func(t *testing.T) {
		if got := NormalizePage(3); got != 3 {
			t.Fatalf("NormalizePage(3)=%d, want 3", got)
		}
	})
}

func TestNormalizePageSize(t *testing.T) {
	t.Run("zero becomes default", func(t *testing.T) {
		if got := NormalizePageSize(0); got != 20 {
			t.Fatalf("NormalizePageSize(0)=%d, want 20", got)
		}
	})
	t.Run("negative becomes default", func(t *testing.T) {
		if got := NormalizePageSize(-10); got != 20 {
			t.Fatalf("NormalizePageSize(-10)=%d, want 20", got)
		}
	})
	t.Run("cap at max", func(t *testing.T) {
		if got := NormalizePageSize(1000); got != 100 {
			t.Fatalf("NormalizePageSize(1000)=%d, want 100", got)
		}
	})
	t.Run("positive unchanged", func(t *testing.T) {
		if got := NormalizePageSize(30); got != 30 {
			t.Fatalf("NormalizePageSize(30)=%d, want 30", got)
		}
	})
}
