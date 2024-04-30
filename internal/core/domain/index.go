package domain

type ComicIndex struct {
	Index map[string][]uint64 `json:"indexes"`
}
