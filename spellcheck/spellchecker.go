package spellcheck

import (
	"fmt"
	"github.com/sajari/fuzzy"
	"os"
	"strings"
)

type SpellChecker interface {
	LoadDataset(dataPath string) error
	SpellCheckString(input string) string
}

type fuzzyChecker struct {
	checker *fuzzy.Model
}

func (checker *fuzzyChecker) LoadDataset(dataPath string) error {
	data, err := os.ReadFile(dataPath)
	if err != nil {
		return fmt.Errorf("dataset file not found: %w", err)
	}
	checker.checker.Train(strings.Split(string(data), "\n"))
	return nil
}

func (checker *fuzzyChecker) SpellCheckString(input string) string {
	words := strings.Fields(input)

	res := make([]string, 0, len(words))
	for _, word := range words {
		spellCheckedWord := checker.checker.SpellCheck(word)
		if spellCheckedWord != "" {
			res = append(res, spellCheckedWord)
		} else {
			res = append(res, word)
		}
	}
	return strings.Join(res, " ")
}

func NewFuzzyChecker(threshold, depth int) SpellChecker {
	model := fuzzy.NewModel()
	model.SetDepth(depth)
	model.SetThreshold(threshold)

	return &fuzzyChecker{checker: model}
}
