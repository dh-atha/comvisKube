package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloudnativeapp/internal/dashboard"
	"cloudnativeapp/internal/identity"
	"cloudnativeapp/internal/store"
)

type App struct {
	startedAt time.Time
	store     *store.EventLog
	server    *http.Server
	identity  identity.Provider
	dashboard *dashboard.Page
}

func New() (*App, error) {
	eventLog := store.New(store.Config{Path: envOrDefault("LOG_PATH", "./logs/events.jsonl")})
	page, err := dashboard.New()
	if err != nil {
		return nil, err
	}

	app := &App{
		startedAt: time.Now().UTC(),
		store:     eventLog,
		identity:  identity.New(),
		dashboard: page,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.handleUI)
	mux.HandleFunc("/api/info", app.handleInfo)
	mux.HandleFunc("/pod", app.handlePod)
	mux.HandleFunc("/work", app.handleWork)
	mux.HandleFunc("/panic", app.handlePanic)
	mux.HandleFunc("/state", app.handleState)
	mux.HandleFunc("/healthz", app.handleHealth)

	app.server = &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	if err := app.store.Write(store.Event{Type: "startup", Note: "application initialized"}); err != nil {
		log.Printf("startup log error: %v", err)
	}

	return app, nil
}

func (a *App) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	shutdownDone := make(chan struct{})

	go func() {
		<-ctx.Done()
		defer close(shutdownDone)

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err := a.store.Write(store.Event{Type: "shutdown", Note: "termination signal received"}); err != nil {
			log.Printf("shutdown log error: %v", err)
		}

		if err := a.server.Shutdown(shutdownCtx); err != nil {
			log.Printf("graceful shutdown error: %v", err)
		}

		if err := a.store.Write(store.Event{Type: "shutdown_complete", Note: "server stopped cleanly"}); err != nil {
			log.Printf("shutdown completion log error: %v", err)
		}
	}()

	log.Printf("listening on %s", a.server.Addr)
	if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	if ctx.Err() != nil {
		<-shutdownDone
	}

	return nil
}

func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
