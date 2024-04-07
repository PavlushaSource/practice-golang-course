package database

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
)

type ComicJsonUnit struct {
	Url      string   `json:"url"`
	Keywords []string `json:"keywords"`
}

type JsonStorage struct {
	storage *os.File
}

func NewJsonStorage(jsonPath string) (*JsonStorage, func() error, error) {
	file, err := os.Create(jsonPath)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating json file: %w", err)
	}

	return &JsonStorage{
		storage: file,
	}, file.Close, nil
}

func (s *JsonStorage) SaveComics(comics map[int]ComicJsonUnit, log *slog.Logger, outputToConsole bool) {
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
				log.Info("comic saved", "comicID", i, "keywords", comics[i].Keywords, "imgUrl", comics[i].Url)
			}
		}
	}
}
