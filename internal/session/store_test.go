package session_test

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vschetko/weekplate/internal/session"
)

func TestNew_Integration(t *testing.T) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		t.Skip("DATABASE_URL not set")
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		t.Fatalf("pool: %v", err)
	}
	defer pool.Close()

	sm, err := session.New(pool, "testsecret")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if sm == nil {
		t.Fatal("expected non-nil session manager")
	}

	// Verify sessions table was created.
	var exists bool
	err = pool.QueryRow(context.Background(),
		`SELECT EXISTS (
			SELECT 1 FROM information_schema.tables
			WHERE table_name = 'sessions'
		)`).Scan(&exists)
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if !exists {
		t.Fatal("sessions table not created")
	}
}

func TestNew_NilPool(t *testing.T) {
	sm, err := session.New(nil, "secret")
	if sm != nil {
		t.Error("expected nil session manager for nil pool")
	}
	if err == nil {
		t.Error("expected error for nil pool")
	}
}
