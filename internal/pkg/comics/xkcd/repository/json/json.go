package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PavlushaSource/yadro-practice-course/internal/core/entities"
	"io"
	"maps"
	"os"
	"path/filepath"
)

type JSONStorage struct {
	filePath string
}

const perm = 0644

func NewJSONComicsStorage(filePath string) (*JSONStorage, error) {
	_, err := os.Stat(filePath)
	if errors.Is(err, os.ErrNotExist) {
		_, err = os.Create(filePath)
		if err != nil {
			return nil, fmt.Errorf("error creating json file: %w", err)
		}
	}
	return &JSONStorage{filePath: filePath}, nil
}

func (s *JSONStorage) Read() (map[int]entities.ComicToJSON, error) {
	var jsonComics map[int]entities.ComicToJSON
	file, err := os.Open(s.filePath)
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	if err != nil {
		return nil, fmt.Errorf("error read json file: %w", err)
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&jsonComics)

	if err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("error decode json file: %w", err)
	}
	return jsonComics, nil
}

// Atomic write file
func (s *JSONStorage) Write(comics map[int]entities.ComicToJSON) error {
	alreadyComics, err := s.Read()
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
	maps.Copy(comics, alreadyComics)

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
