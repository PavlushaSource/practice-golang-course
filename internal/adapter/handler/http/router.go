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
	comicHandler *ComicHandler,
	ctx context.Context,
) (*Router, error) {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /update", comicHandler.UpdateHandler(ctx))
	mux.HandleFunc("/hello", comicHandler.SayHello())
	mux.HandleFunc("GET /pics", comicHandler.GetSuggestRelevantURL())

	//TODO add logger middleware

	server := http.Server{
		Addr:    net.JoinHostPort(cfg.HTTP.Host, cfg.HTTP.Port),
		Handler: mux,
	}

	return &Router{server: &server}, nil
}

func (r *Router) Serve(updateInterval time.Duration, ctx context.Context, comicSvc port.ComicsService) error {
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
			slog.Info("update comics DB")
			_, err := comicSvc.DownloadAll(ctx)
			if err != nil {
				return fmt.Errorf("update comics download failed: %w", err)
			}
		}
	}
}
