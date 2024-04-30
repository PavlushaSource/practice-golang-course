package domain

type Comic struct {
	ID       uint64   `json:"num"`
	URL      string   `json:"url"`
	Keywords []string `json:"keywords"`
}
