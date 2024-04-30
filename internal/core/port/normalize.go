package port

type NormalizeService interface {
	Normalize(phrase string) ([]string, error)
	CorrectAndNormalize(phrase string) ([]string, error)
}
