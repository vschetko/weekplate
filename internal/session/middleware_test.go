package session

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/alexedwards/scs/v2"
)

var uuidRe = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

func TestMiddleware_FirstRequestSetsCookieAndUUID(t *testing.T) {
	old := manager
	manager = scs.New()
	defer func() { manager = old }()

	var capturedID string
	handler := Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedID = GetUserID(r)
		w.WriteHeader(http.StatusOK)
	}))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	handler.ServeHTTP(rec, req)

	if rec.Header().Get("Set-Cookie") == "" {
		t.Fatal("expected Set-Cookie header on first request")
	}
	if capturedID == "" {
		t.Fatal("expected non-empty user_id")
	}
	if !uuidRe.MatchString(capturedID) {
		t.Fatalf("user_id %q does not match UUID v4 format", capturedID)
	}
}

func TestMiddleware_SameUserIDOnSubsequentRequests(t *testing.T) {
	old := manager
	manager = scs.New()
	defer func() { manager = old }()

	var capturedID string
	handler := Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedID = GetUserID(r)
		w.WriteHeader(http.StatusOK)
	}))

	rec1 := httptest.NewRecorder()
	req1 := httptest.NewRequest(http.MethodGet, "/", nil)
	handler.ServeHTTP(rec1, req1)
	firstID := capturedID

	rec2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/", nil)
	for _, c := range rec1.Result().Cookies() {
		req2.AddCookie(c)
	}
	handler.ServeHTTP(rec2, req2)

	if capturedID != firstID {
		t.Fatalf("expected same user_id %q on second request, got %q", firstID, capturedID)
	}
}

func TestMiddleware_NilManager(t *testing.T) {
	old := manager
	manager = nil
	defer func() { manager = old }()

	var called bool
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		if id := GetUserID(r); id != "" {
			t.Errorf("expected empty user_id with nil manager, got %q", id)
		}
		w.WriteHeader(http.StatusOK)
	})

	Middleware(next).ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))

	if !called {
		t.Fatal("expected next handler to be called even with nil manager")
	}
}
