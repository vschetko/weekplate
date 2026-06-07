package session

import (
	"crypto/rand"
	"fmt"
	"net/http"
)

// Middleware mints a guest user_id UUID on first visit and stores it in the
// scs session. Must be used after session.New has been called (manager != nil).
func Middleware(next http.Handler) http.Handler {
	if manager == nil {
		return next
	}
	return manager.LoadAndSave(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if manager.GetString(r.Context(), "user_id") == "" {
			manager.Put(r.Context(), "user_id", newUUID())
		}
		next.ServeHTTP(w, r)
	}))
}

func newUUID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant RFC 4122
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
