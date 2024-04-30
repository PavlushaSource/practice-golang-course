package domain

type ComixIndex struct {
	Index map[string][]uint64 `json:"indexes"`
}
