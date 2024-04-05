package main

import (
	"fmt"
	"github.com/PavlushaSource/yadro-practice-course/pkg/config"
	"github.com/PavlushaSource/yadro-practice-course/pkg/database"
	"github.com/PavlushaSource/yadro-practice-course/pkg/logger"
	"github.com/PavlushaSource/yadro-practice-course/pkg/xkcd"
	"os"
)

func main() {
	// read config file
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// init logger
	log, err := logger.SetupLogger(cfg.Env)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// init storage
	log.Debug("Init database", "db_file", cfg.DbFile)

	storage, cancel, err := database.NewJsonStorage(cfg.DbFile)
	defer cancel()
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	//init client
	log.Debug("Init client")
	client := xkcd.NewClient()

	// read site
	log.Debug("Reading all comics from site")
	comics := xkcd.GetComics(client, cfg.SourceUrl, log)
	if comics == nil {
		log.Error("Cannot read all comics from site")
		os.Exit(1)
	}

	for k := range comics {
		fmt.Println(k)
	}
	_ = storage

}
