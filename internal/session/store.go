package session

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

const createSessionsTable = `
CREATE TABLE IF NOT EXISTS sessions (
	token TEXT PRIMARY KEY,
	data BYTEA NOT NULL,
	expiry TIMESTAMPTZ NOT NULL
)`

// manager is the package-level session manager set by New.
var manager *scs.SessionManager

// New creates a pgxstore-backed scs session manager and ensures the sessions
// table exists. Returns nil, err if pool is nil.
func New(pool *pgxpool.Pool, secret string) (*scs.SessionManager, error) {
	if pool == nil {
		return nil, fmt.Errorf("db pool is nil")
	}

	if err := ensureSchema(context.Background(), pool); err != nil {
		return nil, fmt.Errorf("sessions schema: %w", err)
	}

	sm := scs.New()
	sm.Store = pgxstore.New(pool)
	sm.Lifetime = 30 * 24 * time.Hour
	sm.Cookie.HttpOnly = true
	sm.Cookie.SameSite = http.SameSiteLaxMode
	sm.Cookie.Secure = os.Getenv("APP_ENV") != "development"

	manager = sm
	return sm, nil
}

func ensureSchema(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, createSessionsTable)
	return err
}
