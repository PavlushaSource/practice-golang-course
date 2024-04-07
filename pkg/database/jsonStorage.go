package database

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
)

type JSONComicUnit struct {
	URL      string   `json:"url"`
	Keywords []string `json:"keywords"`
}

type JSONStorage struct {
	storage *os.File
}

func NewJSONStorage(jsonPath string) (*JSONStorage, func() error, error) {
	file, err := os.Create(jsonPath)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating json file: %w", err)
	}

	return &JSONStorage{
		storage: file,
	}, file.Close, nil
}

func (s *JSONStorage) SaveComics(comics map[int]JSONComicUnit, log *slog.Logger, outputToConsole bool) {
	bytes, err := json.MarshalIndent(comics, "", "  ")
	if err != nil {
		log.Error("cannot marshal json", "error", err)
		return
	}
	_, err = s.storage.Write(bytes)
	if err != nil {
		log.Error("cannot write json", "error", err)
		return
	}

	if outputToConsole {
		for i := 1; i <= len(comics); i++ {
			if _, exist := comics[i]; exist {
				log.Info("comic saved", "comicID", i, "keywords", comics[i].Keywords, "imgUrl", comics[i].URL)
			}
		}
	}
}
