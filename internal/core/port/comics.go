package port

import (
	"context"
	"github.com/PavlushaSource/yadro-practice-course/internal/core/domain"
)

type ComicsRepository interface {
	WriteComics(inputComixs []domain.Comic) error
	GetComicByID(ID uint64) (*domain.Comic, error)
	ListComics() ([]domain.Comic, error)
}

type ComicsService interface {
	GetComics() ([]domain.Comic, error)
	DownloadAll(ctx context.Context) ([]domain.Comic, error)
	GetRelevantComics(phrase string, length int) ([]domain.Comic, error)
}
