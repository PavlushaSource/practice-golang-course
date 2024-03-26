package main

import (
	"encoding/json"
	"fmt"
	"github.com/PavlushaSource/yadro-practice-course/spellcheck"
	"github.com/kljensen/snowball"
	"io"
	"os"
	"slices"
	"strings"
)

type ISOCode639_1 string

type Stemmer interface {
	NormalizeString(input string) (string, error)
	NormalizeStringWithSpellcheck(input string, spellchecker spellcheck.SpellChecker) (string, error)
}

type snowballStemmer struct {
	availableLanguages []ISOCode639_1
	stopWords          map[ISOCode639_1][]string
}

func (stemmer *snowballStemmer) DeleteStopWords(currentString string) []string {
	containsMap := func(word string) bool {
		for _, v := range stemmer.stopWords {
			if slices.Contains(v, word) {
				return true
			}
		}
		return false
	}

	res := make([]string, 0, len(currentString))
	for _, word := range strings.Fields(currentString) {
		word = strings.ToLower(word)
		if !containsMap(word) {
			res = append(res, word)
		}
	}
	return slices.Clip(res)
}

func defineLang(word string) ISOCode639_1 {
	return "en"
}

func (stemmer *snowballStemmer) normalizeWords(words []string) ([]string, error) {
	res := make([]string, 0, len(words))
	for _, word := range words {
		usedLanguage := defineLang(word) // TODO: add support auto define language
		switch usedLanguage {
		case "en":
			stemWord, err := snowball.Stem(word, "english", true)
			if err != nil {
				return nil, fmt.Errorf("stemming error: %w", err)
			}
			res = append(res, stemWord)

		default:
			return nil, fmt.Errorf("%s language is not supported now", usedLanguage)
		}
	}
	return slices.Clip(res), nil
}

func (stemmer *snowballStemmer) NormalizeString(input string) (string, error) {
	wordsWithoutStopWords := stemmer.DeleteStopWords(DeleteAllPunctuation(input))
	resString, err := stemmer.normalizeWords(wordsWithoutStopWords)
	if err != nil {
		return "", fmt.Errorf("error normalize string: %w", err)
	}
	return strings.Join(resString, " "), nil
}

func (stemmer *snowballStemmer) NormalizeStringWithSpellcheck(input string, spellchecker spellcheck.SpellChecker) (string, error) {
	wordsSpellChecked := spellchecker.SpellCheckString(DeleteAllPunctuation(input))
	wordsWithoutStopWords := stemmer.DeleteStopWords(wordsSpellChecked)
	resString, err := stemmer.normalizeWords(wordsWithoutStopWords)
	if err != nil {
		return "", fmt.Errorf("error normalize string: %w", err)
	}
	return strings.Join(resString, " "), nil
}

func isSubset[T comparable](gS, lS []T) error {
	for _, l := range lS {
		if !slices.Contains(gS, l) {
			return fmt.Errorf("language %s not supported for stemming", l)
		}
	}
	return nil
}

func NewSnowballStemmer(stopWordsPath string, necessaryLanguages []ISOCode639_1) (Stemmer, error) {
	availableLanguages := []ISOCode639_1{"en"}
	err := isSubset(availableLanguages, necessaryLanguages)
	if err != nil {
		return nil, fmt.Errorf("error necessary languages: %w", err)
	}

	jsonFile, err := os.Open(stopWordsPath)
	defer jsonFile.Close()

	if err != nil {
		return nil, fmt.Errorf("stopwords file not found: %w", err)
	}

	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("error reading stopwords file: %w", err)
	}

	stopWordMapping := make(map[ISOCode639_1][]string)

	err = json.Unmarshal(jsonData, &stopWordMapping)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling stopwords: %w", err)
	}

	necessaryStopWords := make(map[ISOCode639_1][]string, len(necessaryLanguages))
	for _, langISO639 := range necessaryLanguages {
		necessaryStopWords[langISO639] = stopWordMapping[langISO639]
	}

	return &snowballStemmer{stopWords: necessaryStopWords, availableLanguages: availableLanguages}, nil
}

func DeleteAllPunctuation(input string) string {

	// punctuations without \' and -
	punctuations := []rune{'!', '?', '.', ',', ';', ':', '\'', '"', '@', '&', '#', '$', '%', '^', '*', '(', ')', '[', ']',
		'{', '}', '<', '>', '/', '|', '\\', '`', '~', '='}

	b := strings.Builder{}
	for _, r := range input {
		if !slices.Contains(punctuations, r) {
			b.WriteRune(r)
		} else {
			b.WriteRune(' ')
		}
	}
	return b.String()
}
