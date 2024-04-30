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

type ComicJSONStorage struct {
	filePath string
}

const perm = 0644

func NewComicRepository(filePath string) *ComicJSONStorage {
	_, err := os.Stat(filePath)
	if errors.Is(err, os.ErrNotExist) {
		_, err = os.Create(filePath)
		if err != nil {
			log.Fatal("error creating json file", err)
		}
	}
	return &ComicJSONStorage{filePath: filePath}
}

// WriteComics Atomic write one or more comics to JSON DB
func (s *ComicJSONStorage) WriteComics(comics []domain.Comic) error {
	alreadyComics, err := s.ListComics()
	if err != nil {
		return fmt.Errorf("error read json file: %w", err)
	}

	file, err := os.CreateTemp(filepath.Dir(s.filePath), filepath.Base(s.filePath)+".tmp")
	if err != nil {
		return fmt.Errorf("error create temp json file: %w", err)
	}

	tmpName := file.Name()
	defer func() {
		if err != nil {
			file.Close()
			os.Remove(tmpName)
		}
	}()

	comics = append(comics, alreadyComics...)

	bytes, err := json.Marshal(comics)
	if err != nil {
		return fmt.Errorf("error marshal json: %w", err)
	}
	if _, err := file.Write(bytes); err != nil {
		return fmt.Errorf("error write json: %w", err)
	}
	if err := file.Chmod(perm); err != nil {
		return fmt.Errorf("error change permission: %w", err)
	}
	if err := file.Sync(); err != nil {
		return fmt.Errorf("error sync file: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("error close file: %w", err)
	}
	return os.Rename(tmpName, s.filePath)
}

// GetComicByID Read comic by ID from JSON DB
func (s *ComicJSONStorage) GetComicByID(ID uint64) (*domain.Comic, error) {
	comics, err := s.ListComics()
	if err != nil {
		return nil, fmt.Errorf("error read json file: %w", err)
	}

	for _, comic := range comics {
		if comic.ID == ID {
			return &comic, nil
		}
	}

	return nil, nil
}

// ListComics Read all comics from JSON DB
func (s *ComicJSONStorage) ListComics() ([]domain.Comic, error) {
	var comics []domain.Comic
	file, err := os.Open(s.filePath)
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	if err != nil {
		return nil, fmt.Errorf("error list comics: %w", err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&comics)

	if err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("error decode json file: %w", err)
	}
	return comics, nil
}
