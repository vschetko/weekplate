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
)

//go:embed static/index.html
var indexHTML []byte

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
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
		}
	} else {
		log.Println("warn: DATABASE_URL not set — running without database")
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
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

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(indexHTML)
	})

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
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
