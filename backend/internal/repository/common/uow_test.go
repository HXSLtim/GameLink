package common

import (
	"context"
	"testing"

	sqlite "github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestUnitOfWorkWithTx(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	uow := NewUnitOfWork(db)
	err = uow.WithTx(context.Background(), func(r *Repos) error {
		if r.Users == nil || r.Games == nil || r.Players == nil || r.Orders == nil || r.Payments == nil {
			t.Fatal("repository dependencies should be initialised")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("WithTx returned error: %v", err)
	}
}
