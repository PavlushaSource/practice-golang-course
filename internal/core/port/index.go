package port

import "github.com/PavlushaSource/yadro-practice-course/internal/core/domain"

type IndexRepository interface {
	CreateIndex(comix []domain.Comix) (domain.ComixIndex, error)
	UpdateIndex(comix []domain.Comix) (domain.ComixIndex, error)
	Get() (domain.ComixIndex, error)
}
