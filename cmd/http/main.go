package main

import (
	"github.com/PavlushaSource/yadro-practice-course/internal/adapter/config"
	"github.com/PavlushaSource/yadro-practice-course/internal/adapter/logger"
	"github.com/PavlushaSource/yadro-practice-course/internal/adapter/storage/json/repository"
	"log"
	"os"
)

func main() {

	// Read cfg
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	// Set Log
	Log := logger.SetupLogger(cfg.App.Env)

	// Init DB

	DB, err := repository.NewComixRepository(cfg.JSONFlat.DBFilepath)
	if err != nil {
		Log.Error("Failed to init db", "error", err)

		os.Exit(1)
	}

	// Init index DB
	//TODO add index DB initializer

	// Run server
	//TODO add handlers

	_ = DB
}
