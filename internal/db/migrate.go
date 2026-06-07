package db

import (
	"embed"
	"errors"
	"log"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations
var migrationsFS embed.FS

// MustMigrate applies all pending migrations. Panics on any error except
// ErrNoChange (already up-to-date), so a failed migration blocks server
// startup rather than serving with a broken schema.
func MustMigrate(dbURL string) {
	src, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		panic("db: failed to load migration source: " + err.Error())
	}

	// golang-migrate's pgx/v5 driver registers under the "pgx5" scheme.
	// Rewrite postgres:// / postgresql:// before handing off.
	m, err := migrate.NewWithSourceInstance("iofs", src, toPgx5URL(dbURL))
	if err != nil {
		panic("db: failed to create migrator: " + err.Error())
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		panic("db: migration failed: " + err.Error())
	}

	log.Println("database migrations applied")
}

// toPgx5URL rewrites a postgres:// or postgresql:// URL to pgx5://.
func toPgx5URL(dbURL string) string {
	switch {
	case strings.HasPrefix(dbURL, "postgresql://"):
		return "pgx5" + dbURL[len("postgresql"):]
	case strings.HasPrefix(dbURL, "postgres://"):
		return "pgx5" + dbURL[len("postgres"):]
	default:
		return dbURL
	}
}
