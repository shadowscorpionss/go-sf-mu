package postgres

import (
	"context"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := New(ctx, "postgres://postgres:postgres@localhost:5432/comm")
	if err != nil {
		t.Fatal(err)
	}
}
