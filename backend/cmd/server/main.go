// Command server runs the calculator HTTP API.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yzbaris/sezzle-calculator/backend/internal/api"
)

// newHandler returns the API router. If staticDir is set (as in the Docker image),
// it also serves the built frontend from that folder, so everything runs on one port.
func newHandler(staticDir string) http.Handler {
	router := api.NewRouter()
	if staticDir == "" {
		return router
	}

	mux := http.NewServeMux()
	mux.Handle("/api/", router)
	mux.Handle("/healthz", router)
	mux.Handle("/", http.FileServer(http.Dir(staticDir)))
	return mux
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           newHandler(os.Getenv("STATIC_DIR")),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// ctx is cancelled when the process gets Ctrl+C (SIGINT) or SIGTERM (docker stop).
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Print("shutting down")

	// Give in-flight requests up to 5 seconds to finish.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}
