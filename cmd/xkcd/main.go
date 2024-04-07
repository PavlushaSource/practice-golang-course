package main

import (
	"fmt"
	"github.com/PavlushaSource/yadro-practice-course/internal/flags"
	"github.com/PavlushaSource/yadro-practice-course/pkg/config"
	"github.com/PavlushaSource/yadro-practice-course/pkg/database"
	"github.com/PavlushaSource/yadro-practice-course/pkg/logger"
	"github.com/PavlushaSource/yadro-practice-course/pkg/words/stemmer"
	"github.com/PavlushaSource/yadro-practice-course/pkg/xkcd"
	"os"
	"sync"
)

type UsersFlags struct {
	NumberComics    int
	OutputToConsole bool
}

func main() {

	// parse flags
	userFlags, err := flags.GetFlagsFromCommandlineInput()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

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
	var comics map[int]xkcd.ComicInfo
	if userFlags.NumberComics > 0 {
		log.Debug(fmt.Sprintf("Read only %d comics", userFlags.NumberComics))
		comics = xkcd.GetComics(client, cfg.SourceUrl, log, userFlags.NumberComics)
	} else {
		log.Debug(fmt.Sprintf("Read all comics from %s", cfg.SourceUrl))
		comics = xkcd.GetComics(client, cfg.SourceUrl, log)
	}
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

	// normalize transcripts
	comicsToJson := make(map[int]database.ComicJsonUnit)
	log.Debug("Normalize transcripts")
	wg := sync.WaitGroup{}
	mapMutex := sync.Mutex{}
	for _, comic := range comics {
		wg.Add(1)
		go func(comic xkcd.ComicInfo) {
			defer wg.Done()
			var currentUnit database.ComicJsonUnit
			currentUnit.Url = comic.Img
			keywords, err := st.NormalizeString(comic.Transcript)
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
	storage.SaveComics(comicsToJson, log, userFlags.OutputToConsole)

}
