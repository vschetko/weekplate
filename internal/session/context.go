package session

import "net/http"

// GetUserID returns the guest UUID stored in the active session.
// Returns "" if called outside session middleware.
func GetUserID(r *http.Request) string {
	if manager == nil {
		return ""
	}
	return manager.GetString(r.Context(), "user_id")
}
