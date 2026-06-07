package session_test

import (
	"context"
	"net/http"
	"net/http/httptest"
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

func TestSessionRoundTrip_Integration(t *testing.T) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		t.Skip("DATABASE_URL not set")
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		t.Fatalf("pool: %v", err)
	}
	defer pool.Close()

	if _, err := session.New(pool, "testsecret"); err != nil {
		t.Fatalf("session.New: %v", err)
	}

	var capturedID string
	handler := session.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedID = session.GetUserID(r)
		w.WriteHeader(http.StatusOK)
	}))

	// First request — mints a new UUID and persists it in pgxstore.
	rec1 := httptest.NewRecorder()
	req1 := httptest.NewRequest(http.MethodGet, "/", nil)
	handler.ServeHTTP(rec1, req1)
	firstID := capturedID
	if firstID == "" {
		t.Fatal("expected non-empty user_id on first request")
	}

	// Second request with same cookie — loads UUID from DB-backed store.
	rec2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/", nil)
	for _, c := range rec1.Result().Cookies() {
		req2.AddCookie(c)
	}
	handler.ServeHTTP(rec2, req2)
	if capturedID != firstID {
		t.Fatalf("expected same user_id %q, got %q", firstID, capturedID)
	}
}
