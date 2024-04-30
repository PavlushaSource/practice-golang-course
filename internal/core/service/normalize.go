package service

import (
	"github.com/PavlushaSource/yadro-practice-course/pkg/words/spellcheck"
	"github.com/PavlushaSource/yadro-practice-course/pkg/words/stemmer"
)

type NormalizeService struct {
	stemmer stemmer.Stemmer
	checker spellcheck.SpellChecker
}

func NewNormalizeService(st stemmer.Stemmer, ch spellcheck.SpellChecker) *NormalizeService {
	return &NormalizeService{
		stemmer: st,
		checker: ch,
	}
}

func (n *NormalizeService) Normalize(phrase string) ([]string, error) {
	return n.stemmer.NormalizeString(phrase)
}

func (n *NormalizeService) CorrectAndNormalize(phrase string) ([]string, error) {
	return n.stemmer.NormalizeString(phrase, n.checker)
}
