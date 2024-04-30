package port

import "github.com/PavlushaSource/yadro-practice-course/internal/core/domain"

type IndexRepository interface {
	UpdateIndex(comix []domain.Comic) (domain.ComicIndex, error)
	Get() (domain.ComicIndex, error)
}
