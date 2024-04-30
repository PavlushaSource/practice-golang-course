package main

import (
	"context"
	"github.com/PavlushaSource/yadro-practice-course/internal/adapter/config"
	"github.com/PavlushaSource/yadro-practice-course/internal/adapter/handler/http"
	"github.com/PavlushaSource/yadro-practice-course/internal/adapter/logger"
	"github.com/PavlushaSource/yadro-practice-course/internal/adapter/storage/json/repository"
	"github.com/PavlushaSource/yadro-practice-course/internal/core/service"
	"github.com/PavlushaSource/yadro-practice-course/pkg/words/spellcheck"
	"github.com/PavlushaSource/yadro-practice-course/pkg/words/stemmer"
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

	// ctx with interruption
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Dependency injection

	// Normalize
	slog.Info("Initialize stemmer")
	st, err := stemmer.NewSnowballStemmer(cfg.Normalize.StopwordsPath)
	if err != nil {
		slog.Error("Initialize stemmer failed", "error", err)
		os.Exit(1)
	}
	ch := spellcheck.NewFuzzyChecker(
		cfg.Spellchecker.ModelPath,
		cfg.Spellchecker.AllWordsPath,
		[]string{cfg.Spellchecker.DictPathEn, cfg.Spellchecker.DictPathRus},
	)
	normalizeService := service.NewNormalizeService(st, ch)

	// Comix
	slog.Info("Initialize and try download comixs")
	comixService := service.NewComixService(indexDB, DB, normalizeService, cfg)
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
