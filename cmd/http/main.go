package main

import (
	"context"
	"github.com/PavlushaSource/yadro-practice-course/internal/adapter/config"
	"github.com/PavlushaSource/yadro-practice-course/internal/adapter/handler/http"
	"github.com/PavlushaSource/yadro-practice-course/internal/adapter/logger"
	"github.com/PavlushaSource/yadro-practice-course/internal/adapter/storage/json/repository"
	"github.com/PavlushaSource/yadro-practice-course/internal/core/service"
	"log"
	"log/slog"
	"os"
	"os/signal"
)

func main() {

	// Read cfg
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	// Set Log
	slog.SetDefault(logger.SetupLogger(cfg.App.Env))

	// Init DB

	DB := repository.NewComixRepository(cfg.JSONFlat.DBFilepath)
	indexDB := repository.NewIndexRepository(cfg.JSONFlat.IndexFilepath)

	// Dependency injection
	// Comix
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	slog.Info("Initialize and download comixs")
	comixService := service.NewComixService(indexDB, DB, cfg)
	comixHandler := http.NewComixHandler(comixService)

	_, err = comixService.DownloadAll(ctx)
	if err != nil {
		slog.Error("Download failed", "error", err)
		os.Exit(1)
	}

	// Run server
	srv, err := http.NewRouter(cfg, comixHandler, ctx)
	if err != nil {
		slog.Error("Initialize server failed", "error", err)
		os.Exit(1)
	}
	err = srv.Serve(cfg.HTTP.UpdateInterval, ctx, comixService)
	if err != nil {
		slog.Error("Run server failed", "error", err)
		os.Exit(1)
	}
}
