package spellcheck

import (
	"errors"
	"fmt"
	"github.com/sajari/fuzzy"
	"io"
	"os"
	"strings"
)

type SpellChecker interface {
	LoadDataset(dataPath string) error
	SaveModel(filename string) error
	SpellCheckString(input string) string
}

type fuzzyChecker struct {
	checker  *fuzzy.Model
	allWords map[string]struct{}
}

// LoadDataset load dataset from file and save in your spell checker
func (checker *fuzzyChecker) LoadDataset(dataPath string) error {
	data, err := os.ReadFile(dataPath)
	if err != nil {
		return fmt.Errorf("dataset file not found: %w", err)
	}
	checker.checker.Train(strings.Split(string(data), "\n"))
	return nil
}

func (checker *fuzzyChecker) SaveModel(filename string) error {
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		println("saving model...")
		err = checker.checker.Save(filename)
		if err != nil {
			return err
		}
	}
	return nil
}

func (checker *fuzzyChecker) SpellCheckString(input string) string {
	words := strings.Fields(input)
	res := make([]string, 0, len(words))
	for _, word := range words {
		if _, exist := checker.allWords[word]; exist {
			res = append(res, word)
			continue
		}

		spellCheckedWord := checker.checker.SpellCheck(word)
		if spellCheckedWord != "" {
			res = append(res, spellCheckedWord)
		} else {
			res = append(res, word)
		}
	}
	return strings.Join(res, " ")
}

func loadAllWords(path string) map[string]struct{} {
	resMap := make(map[string]struct{}, 100000)
	f, _ := os.Open(path)

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(f)

	allWords, _ := io.ReadAll(f)
	for _, word := range strings.Split(string(allWords), "\n") {
		resMap[word] = struct{}{}
	}
	return resMap
}

func NewFuzzyChecker(args ...int) SpellChecker {
	loadedModel, errLoad := fuzzy.Load("spellcheck/savedModel")

	if errLoad != nil {
		model := fuzzy.NewModel()
		if len(args) == 0 {
			model.SetDepth(2)
			model.SetThreshold(1)
		}
		if len(args) >= 1 {
			model.SetDepth(args[0])
		}
		if len(args) >= 2 {
			model.SetThreshold(args[1])
		}
		fzCheck := fuzzyChecker{checker: model, allWords: loadAllWords("spellcheck/all-words.txt")}
		_ = fzCheck.LoadDataset("spellcheck/10000-english.txt")
		_ = fzCheck.LoadDataset("spellcheck/10000-russian.txt")

		return &fzCheck
	}

	fzChecker := fuzzyChecker{checker: loadedModel, allWords: loadAllWords("spellcheck/all-words.txt")}
	return &fzChecker
}
