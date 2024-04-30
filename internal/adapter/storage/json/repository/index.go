package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PavlushaSource/yadro-practice-course/internal/core/domain"
	"io"
	"log"
	"os"
	"path/filepath"
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
	index := domain.ComixIndex{Index: make(map[string][]uint64)}
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

func (r *IndexJSONStorage) UpdateIndex(inputComixs []domain.Comix) (domain.ComixIndex, error) {
	index, err := r.Get()
	if err != nil {
		return index, fmt.Errorf("error get index: %w", err)
	}

	for _, comix := range inputComixs {
		for _, keywords := range comix.Keywords {
			index.Index[keywords] = append(index.Index[keywords], comix.ID)
		}
	}

	// write to file index
	file, err := os.CreateTemp(filepath.Dir(r.filePath), filepath.Base(r.filePath)+".tmp")
	if err != nil {
		return index, fmt.Errorf("error create temp json file: %w", err)
	}

	tmpName := file.Name()
	defer func() {
		if err != nil {
			file.Close()
			os.Remove(tmpName)
		}
	}()

	bytes, err := json.Marshal(index)
	if err != nil {
		return index, fmt.Errorf("error marshal json: %w", err)
	}
	if _, err := file.Write(bytes); err != nil {
		return index, fmt.Errorf("error write json: %w", err)
	}
	if err := file.Chmod(perm); err != nil {
		return index, fmt.Errorf("error change permission: %w", err)
	}
	if err := file.Sync(); err != nil {
		return index, fmt.Errorf("error sync file: %w", err)
	}
	if err := file.Close(); err != nil {
		return index, fmt.Errorf("error close file: %w", err)
	}
	return index, os.Rename(tmpName, r.filePath)
}
