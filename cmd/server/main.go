package main

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/vschetko/weekplate/internal/db"
	"github.com/vschetko/weekplate/internal/session"
)

//go:embed static/index.html
var indexHTML []byte

func main() {
	_ = godotenv.Load() // load .env if present; ignore error when absent

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		log.Fatal("SESSION_SECRET is required")
	}

	var pool *pgxpool.Pool
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		var err error
		pool, err = pgxpool.New(context.Background(), dbURL)
		if err != nil {
			log.Printf("warn: failed to create db pool: %v", err)
		} else if err = pool.Ping(context.Background()); err != nil {
			log.Printf("warn: db ping failed: %v", err)
			pool.Close()
			pool = nil
		} else {
			log.Println("database connection established")
			defer pool.Close()
			db.MustMigrate(dbURL)
		}
	} else {
		log.Println("warn: DATABASE_URL not set — running without database")
	}

	if pool != nil {
		if _, err := session.New(pool, secret); err != nil {
			log.Printf("warn: session store init failed: %v", err)
		}
	}

	healthMux := http.NewServeMux()
	healthMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if pool != nil {
			ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
			defer cancel()
			if err := pool.Ping(ctx); err != nil {
				http.Error(w, fmt.Sprintf("db unreachable: %v", err), http.StatusServiceUnavailable)
				return
			}
		}
		fmt.Fprint(w, "ok")
	})

	appMux := http.NewServeMux()
	appMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if id := session.GetUserID(r); id != "" {
			log.Printf("request user_id: %s", id)
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(indexHTML)
	})

	// session.Middleware is a no-op when session.New was not called (manager == nil)
	rootHandler := session.Middleware(appMux)

	srv := &http.Server{
		Addr: ":" + port,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/health" {
				healthMux.ServeHTTP(w, r)
				return
			}
			rootHandler.ServeHTTP(w, r)
		}),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("server listening on :%s", port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-quit
	log.Println("shutting down…")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}
	log.Println("server stopped")
}
