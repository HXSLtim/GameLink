package model

import "testing"

func TestRatingValid(t *testing.T) {
	if !Rating(1).Valid() || !Rating(5).Valid() {
		t.Fatalf("expected bounds 1 and 5 to be valid")
	}
	if Rating(0).Valid() || Rating(6).Valid() {
		t.Fatalf("expected 0 and 6 to be invalid")
	}
}

func TestMustRating(t *testing.T) {
	_ = MustRating(3)
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic for out-of-range value")
		}
	}()
	_ = MustRating(9)
}
