package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/PavlushaSource/yadro-practice-course/internal/adapter/config"
	"github.com/PavlushaSource/yadro-practice-course/internal/core/port"
	"log/slog"
	"net"
	"net/http"
	"time"
)

type Router struct {
	server *http.Server
}

func NewRouter(
	cfg *config.Config,
	comixHandler *ComixHandler,
	ctx context.Context,
) (*Router, error) {
	mux := http.NewServeMux()

	//Get from Index file
	//mux.HandleFunc("GET /pics", comixHandler.SuggestRelevantURLIndex(ctx))

	// Uncomment for compare simple find relevant URL
	//mux.HandleFunc("GET /pics", comixHandler.GetRelevantURL)

	//mux.HandleFunc("POST /update", comixHandler.Update(ctx))

	mux.HandleFunc("/hello", comixHandler.SayHello())
	mux.HandleFunc("GET /pics", comixHandler.GetSuggestRelevantURL())

	//TODO add logger middleware

	server := http.Server{
		Addr:    net.JoinHostPort(cfg.HTTP.Host, cfg.HTTP.Port),
		Handler: mux,
	}

	return &Router{server: &server}, nil
}

func (r *Router) Serve(updateInterval time.Duration, ctx context.Context, comixSvc port.ComixService) error {
	ticker := time.NewTicker(updateInterval)

	go func() {
		slog.Info("listen server on", "addr", r.server.Addr)
		err := r.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "error", err)
		}
	}()
	for {
		select {
		case <-ctx.Done():
			shutdownCtx := context.Background()
			shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
			defer cancel()

			if err := r.server.Shutdown(shutdownCtx); err != nil {
				return fmt.Errorf("server shutdown failed: %w", err)
			}
			return nil
		case <-ticker.C:
			slog.Info("update comix DB")
			_, err := comixSvc.DownloadAll(ctx)
			if err != nil {
				return fmt.Errorf("update comix download failed: %w", err)
			}
		}
	}
}
