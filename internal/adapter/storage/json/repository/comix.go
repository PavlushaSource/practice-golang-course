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

type ComixJSONStorage struct {
	filePath string
}

const perm = 0644

func NewComixRepository(filePath string) *ComixJSONStorage {
	_, err := os.Stat(filePath)
	if errors.Is(err, os.ErrNotExist) {
		_, err = os.Create(filePath)
		if err != nil {
			log.Fatal("error creating json file", err)
		}
	}
	return &ComixJSONStorage{filePath: filePath}
}

// WriteComixs Atomic write one or more comix to JSON DB
func (s *ComixJSONStorage) WriteComixs(comix []domain.Comix) error {
	alreadyComixs, err := s.ListComixs()
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

	comix = append(comix, alreadyComixs...)

	bytes, err := json.Marshal(comix)
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

// GetComixByID Read comix by ID from JSON DB
func (s *ComixJSONStorage) GetComixByID(ID uint64) (*domain.Comix, error) {
	comixs, err := s.ListComixs()
	if err != nil {
		return nil, fmt.Errorf("error read json file: %w", err)
	}

	for _, comix := range comixs {
		if comix.ID == ID {
			return &comix, nil
		}
	}

	return nil, nil
}

// ListComixs Read all comixs from JSON DB
func (s *ComixJSONStorage) ListComixs() ([]domain.Comix, error) {
	var comixs []domain.Comix
	file, err := os.Open(s.filePath)
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	if err != nil {
		return nil, fmt.Errorf("error list comixs: %w", err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&comixs)

	if err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("error decode json file: %w", err)
	}
	return comixs, nil
}
