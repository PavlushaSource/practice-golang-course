package port

import "github.com/PavlushaSource/yadro-practice-course/internal/core/domain"

type IndexRepository interface {
	UpdateIndex(comix []domain.Comix) (domain.ComixIndex, error)
	Get() (domain.ComixIndex, error)
}
