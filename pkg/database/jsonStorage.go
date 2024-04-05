package database

import (
	"fmt"
	"os"
)

type ComicUnit struct {
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
