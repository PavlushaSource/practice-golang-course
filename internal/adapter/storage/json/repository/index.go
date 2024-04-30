package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PavlushaSource/yadro-practice-course/internal/core/domain"
	"io"
	"log"
	"os"
)

type IndexJSONStorage struct {
	filePath string
}

func NewIndexRepository(filePath string) *IndexJSONStorage {
	_, err := os.Stat(filePath)
	if errors.Is(err, os.ErrNotExist) {
		_, err = os.Create(filePath)
		if err != nil {
			log.Fatal("error creating json file", err)
		}
	}
	return &IndexJSONStorage{filePath: filePath}
}

func (r *IndexJSONStorage) Get() (domain.ComixIndex, error) {
	index := domain.ComixIndex{Index: make(map[string][]int)}
	file, err := os.Open(r.filePath)
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	if err != nil {
		return index, fmt.Errorf("can't open index file: %w", err)
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&index)

	if err != nil && !errors.Is(err, io.EOF) {
		return index, fmt.Errorf("error decode json file: %w", err)
	}

	return index, nil
}

func (r *IndexJSONStorage) UpdateIndex(comix []domain.Comix) (domain.ComixIndex, error) {
	index := domain.ComixIndex{Index: make(map[string][]int)}
	return index, nil
}
