package port

import (
	"context"
	"github.com/PavlushaSource/yadro-practice-course/internal/core/domain"
)

type ComixRepository interface {
	WriteComixs(comix []domain.Comix) error
	GetComixByID(ID uint64) (*domain.Comix, error)
	ListComixs() ([]domain.Comix, error)
}

type ComixService interface {
	DownloadAll(ctx context.Context) ([]domain.Comix, error)
	GetRelevantComixs(phrase string) ([]domain.Comix, error)
	GetRelevantComixsIndex(phrase string) ([]domain.Comix, error)
}
