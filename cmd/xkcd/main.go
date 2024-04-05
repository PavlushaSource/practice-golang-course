package main

import (
	"fmt"
	"github.com/PavlushaSource/yadro-practice-course/pkg/config"
	"github.com/PavlushaSource/yadro-practice-course/pkg/database"
	"github.com/PavlushaSource/yadro-practice-course/pkg/logger"
	"github.com/PavlushaSource/yadro-practice-course/pkg/words/spellcheck"
	"github.com/PavlushaSource/yadro-practice-course/pkg/words/stemmer"
	"github.com/PavlushaSource/yadro-practice-course/pkg/xkcd"
	"os"
	"sync"
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

	// init stemmer
	log.Debug("Init stemmer")
	st, err := stemmer.NewSnowballStemmer()
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	// init spellchecker
	log.Debug("Init spellchecker")
	checker := spellcheck.NewFuzzyChecker()
	err = checker.SaveModel(cfg.SpellcheckModel)
	stemmer.Check(err)

	// normalize transcripts
	comicsToJson := make(map[int]database.ComicUnit)
	log.Debug("Normalize transcripts")
	wg := sync.WaitGroup{}
	mapMutex := sync.Mutex{}
	for _, comic := range comics {
		wg.Add(1)
		go func(comic xkcd.ComicInfo) {
			defer wg.Done()
			var currentUnit database.ComicUnit
			currentUnit.Url = comic.Img
			keywords, err := st.NormalizeString(comic.Transcript, checker)
			if err != nil {
				log.Error("Cannot normalize transcript", "comicID", comic.Num, "error", err)
			}
			currentUnit.Keywords = keywords
			mapMutex.Lock()
			comicsToJson[comic.Num] = currentUnit
			mapMutex.Unlock()
		}(comic)
	}
	wg.Wait()

	// save to database
	log.Debug("Save to database")
	storage.SaveComics(comicsToJson, log)
}
