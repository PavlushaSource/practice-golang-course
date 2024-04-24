package main

import (
	"context"
	"fmt"
	"github.com/PavlushaSource/yadro-practice-course/internal/config"
	"github.com/PavlushaSource/yadro-practice-course/internal/core/entities"
	"github.com/PavlushaSource/yadro-practice-course/internal/flags"
	"github.com/PavlushaSource/yadro-practice-course/internal/logger"
	"github.com/PavlushaSource/yadro-practice-course/internal/pkg/comics/xkcd"
	jsonStorage "github.com/PavlushaSource/yadro-practice-course/internal/pkg/comics/xkcd/repository/json"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"slices"
)

func getAlreadyDownloadIDs(db DB) ([]int, error) {
	alreadyDownloadedComics, err := db.Read()
	if err != nil {
		return nil, err
	}
	IDs := make([]int, 0)
	for k := range alreadyDownloadedComics {
		IDs = append(IDs, k)
	}
	return IDs, nil
}

type DB interface {
	Read() (map[int]entities.ComicToJSON, error)
	Write(comics map[int]entities.ComicToJSON) error
}

func main() {

	// parse flags
	userFlags := flags.GetFlagsFromCommandlineInput()

	// read config file
	cfg, err := config.LoadConfig(userFlags)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// init logger
	log, err := logger.SetupLogger(cfg.Env)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// init storage
	log.Debug("Init database", "db_file", cfg.DBFile)

	var storage DB
	storage, err = jsonStorage.NewJSONComicsStorage(cfg.DBFile)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	//init client
	log.Debug("Init client")
	client := xkcd.NewClient()

	// read site

	log.Debug("Start read site")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	wg, ctx := errgroup.WithContext(ctx)
	defer cancel()

	neededComicsID := make(chan int)
	batches := make(chan map[int]entities.ComicToJSON)

	// function to get needed IDs
	go func(ctx context.Context, jobs chan<- int) {
		defer close(jobs)
		var alreadyDownloadedIDs []int
		alreadyDownloadedIDs, err = getAlreadyDownloadIDs(storage)
		skipID := append(alreadyDownloadedIDs, 404)
		if err != nil {
			log.Error("Error getting already downloaded IDs", "error", err)
			os.Exit(1)
		}

		for i := 1; ; i++ {
			select {
			case <-ctx.Done():
				return
			default:
				if !slices.Contains(skipID, i) {
					jobs <- i
				}
			}
		}
	}(ctx, neededComicsID)

	// start worker pool
	for w := 1; w <= cfg.Parallel; w++ {
		wg.Go(func() error {
			return getBatchesFromSite(neededComicsID, batches, client, log, cfg.BatchSize, ctx)
		})
	}

	go func() {
		errWg := wg.Wait()
		if errWg != nil {
			log.Debug("Error in worker pool", "error", errWg)
		}
		close(batches)
		log.Debug("Worker pool closed")
	}()

	// write to DB our batches
	for batch := range batches {
		//convert batch to JSON format
		err := storage.Write(batch)
		if err != nil {
			log.Error("Error writing batch to DB", "error", err)
			os.Exit(1)
		}
	}
	log.Debug("Write is finished")

}

func getBatchesFromSite(IDs <-chan int, batches chan<- map[int]entities.ComicToJSON, client *http.Client, log *slog.Logger, size int, ctx context.Context) error {
	batch := make(map[int]entities.ComicToJSON)
	defer func() {
		batches <- batch
	}()
	for {
		select {
		case ID, ok := <-IDs:
			if !ok {
				return nil
			}
			comicInfo, err := xkcd.GetComicByID(client, ID)
			if err != nil {
				return err
			}

			comicJSON, _ := comicInfo.ToJSON()
			batch[comicInfo.Num] = *comicJSON

			if len(batch) == size {
				batches <- batch
				batch = map[int]entities.ComicToJSON{}
			}
		case <-ctx.Done():
			return nil
		}
	}
}
