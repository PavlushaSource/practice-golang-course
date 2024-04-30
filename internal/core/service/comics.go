package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PavlushaSource/yadro-practice-course/internal/adapter/config"
	"github.com/PavlushaSource/yadro-practice-course/internal/core/domain"
	xkcdDomain "github.com/PavlushaSource/yadro-practice-course/internal/core/domain/xkcd"
	"github.com/PavlushaSource/yadro-practice-course/internal/core/port"
	"github.com/PavlushaSource/yadro-practice-course/internal/core/util/xkcd"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"slices"
	"sort"
)

type ComicsService struct {
	indexRepo    port.IndexRepository
	comicsRepo   port.ComicsRepository
	normalizeSrv port.NormalizeService
	batchSize    int
	workers      int
	siteURL      string
}

func (s *ComicsService) GetComics() ([]domain.Comic, error) {
	comics, err := s.comicsRepo.ListComics()
	if err != nil {
		return nil, fmt.Errorf("error get comics: %w", err)
	}
	return comics, nil
}

func (s *ComicsService) GetRelevantComics(phrase string, length int) ([]domain.Comic, error) {
	// correct and normalize user request
	keywords, err := s.normalizeSrv.CorrectAndNormalize(phrase)
	if err != nil {
		return nil, fmt.Errorf("error normalize phrase: %w", err)
	}

	// get index map
	index, err := s.indexRepo.Get()
	if err != nil {
		return nil, fmt.Errorf("error get index: %w", err)
	}

	countID := make(map[uint64]int)
	for _, keyword := range keywords {
		for _, id := range index.Index[keyword] {
			countID[id]++

		}
	}

	neededID := make([]uint64, 0)
	for ID := range countID {
		neededID = append(neededID, ID)
	}

	sort.SliceStable(
		neededID, func(i, j int) bool {
			return countID[neededID[uint64(i)]] > countID[neededID[uint64(j)]]
		},
	)

	resultSlice := make([]uint64, 0, length)
	for _, ID := range neededID {
		resultSlice = append(resultSlice, ID)
		if len(resultSlice) == length {
			break
		}
	}

	// get comics by needed ID
	comicsSlice := make([]domain.Comic, 0, len(resultSlice))
	for _, id := range resultSlice {
		comic, _ := s.comicsRepo.GetComicByID(id)
		comicsSlice = append(comicsSlice, *comic)
	}
	return comicsSlice, nil
}

func NewComicsService(
	indexRepo port.IndexRepository,
	comixRepo port.ComicsRepository,
	normalizeSrv port.NormalizeService,
	cfg *config.Config,
) *ComicsService {
	return &ComicsService{
		indexRepo:    indexRepo,
		comicsRepo:   comixRepo,
		normalizeSrv: normalizeSrv,
		batchSize:    cfg.ComicsSource.BatchSize,
		workers:      cfg.ComicsSource.Parallel,
		siteURL:      cfg.ComicsSource.URL,
	}
}

func (s *ComicsService) DownloadAll(ctx context.Context) ([]domain.Comic, error) {
	wg, ctx := errgroup.WithContext(ctx)

	neededComicsID := make(chan uint64)
	batches := make(chan []domain.Comic)
	// create client
	client := xkcd.NewClient()

	// Generate ID (jobs)
	go func() {
		err := generateID(ctx, neededComicsID, s.comicsRepo)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// start worker pool
	for w := 1; w <= s.workers; w++ {
		wg.Go(
			func() error {
				return downloadWorker(ctx, neededComicsID, batches, client, s.batchSize, s.siteURL, s.normalizeSrv)
			},
		)
	}

	comics := make([]domain.Comic, 0)

	go func() {
		errWg := wg.Wait()
		if errWg != nil && !errors.Is(errWg, domain.ErrStatusNotOK) {
			log.Fatal(errWg)
		}
		close(batches)
	}()

	// write to DB our batches
	for batch := range batches {
		comics = append(comics, batch...)

		err := s.comicsRepo.WriteComics(batch)
		if err != nil {
			log.Fatal(err)
		}
	}

	_, err := s.indexRepo.UpdateIndex(comics)
	if err != nil {
		log.Fatal(err)
	}

	return comics, nil
}

func downloadComicByID(client *http.Client, ID uint64, siteURL string) (*xkcdDomain.Comic, error) {
	url := fmt.Sprintf("%s/%d/info.0.json", siteURL, ID)
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("cannot get comic from url %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cannot get comic from url %s: %w", url, domain.ErrStatusNotOK)
	}

	var comic xkcdDomain.Comic
	err = json.NewDecoder(resp.Body).Decode(&comic)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal comic from url %s: %w", url, err)
	}
	return &comic, nil
}

func generateID(ctx context.Context, jobs chan<- uint64, storage port.ComicsRepository) error {
	defer close(jobs)
	comicAlreadyExist, err := storage.ListComics()
	if err != nil {
		return fmt.Errorf("error getting already downloaded IDs: %w", err)
	}
	alreadyExistID := make([]uint64, 0)
	for _, v := range comicAlreadyExist {
		alreadyExistID = append(alreadyExistID, v.ID)
	}

	// Skip 404 joke page
	alreadyExistID = append(alreadyExistID, 404)

	// Write ID to channel before interruption or not status request equal 200
	for i := uint64(1); ; i++ {
		select {
		case <-ctx.Done():
			return nil
		default:
			if !slices.Contains(alreadyExistID, i) {
				jobs <- i
			}
		}
	}
}

func downloadWorker(
	ctx context.Context,
	ID <-chan uint64,
	batches chan<- []domain.Comic,
	client *http.Client,
	batchSize int,
	siteURL string,
	normalizeSrv port.NormalizeService,
) error {
	batch := make([]domain.Comic, 0, batchSize)
	defer func() {
		batches <- batch
	}()
	for {
		select {
		case currID, ok := <-ID:
			if !ok {
				return nil
			}
			xkcdComic, err := downloadComicByID(client, currID, siteURL)
			if err != nil {
				return fmt.Errorf("error downloading comic: %w", err)
			}

			comic := xkcdComic.Format(normalizeSrv)
			batch = append(batch, comic)

			if len(batch) == batchSize {
				batches <- batch
				batch = make([]domain.Comic, 0, batchSize)
			}
		case <-ctx.Done():
			return nil
		}
	}
}
