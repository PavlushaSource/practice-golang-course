package main

import (
	"encoding/json"
	"fmt"
	"github.com/PavlushaSource/yadro-practice-course/spellcheck"
	"github.com/kljensen/snowball"
	"github.com/pemistahl/lingua-go"
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

func (stemmer *snowballStemmer) deleteStopWords(currentString string) []string {
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

func (stemmer *snowballStemmer) normalizeWords(words []string) ([]string, error) {
	mappingIso639 := map[ISOCode639_1]lingua.Language{
		"en": lingua.English,
		"ru": lingua.Russian,
		"es": lingua.Spanish,
		"fr": lingua.French,
		"sw": lingua.Swedish,
		"hu": lingua.Hungarian,
	}

	defineLanguages := make([]lingua.Language, 0)
	for _, langIso639 := range stemmer.availableLanguages {
		l, exist := mappingIso639[langIso639]
		if !exist {
			return nil, fmt.Errorf("language %v not supported for stemming", langIso639)
		}
		defineLanguages = append(defineLanguages, l)
	}

	detector := lingua.NewLanguageDetectorBuilder().FromLanguages(defineLanguages...).Build()
	res := make([]string, 0, len(words))
	for _, word := range words {
		if usedLanguage, exist := detector.DetectLanguageOf(word); exist {
			stemWord, err := snowball.Stem(word, strings.ToLower(usedLanguage.String()), true)
			if err != nil {
				return nil, fmt.Errorf("stemming error: %w", err)
			}
			res = append(res, stemWord)
		} else {
			return nil, fmt.Errorf("language for word: %s is not detected", word)
		}
	}
	return slices.Clip(res), nil
}

func (stemmer *snowballStemmer) NormalizeString(input string) (string, error) {
	wordsWithoutStopWords := stemmer.deleteStopWords(deleteAllPunctuation(input))
	resString, err := stemmer.normalizeWords(wordsWithoutStopWords)
	if err != nil {
		return "", fmt.Errorf("error normalize string: %w", err)
	}
	return strings.Join(resString, " "), nil
}

func (stemmer *snowballStemmer) NormalizeStringWithSpellcheck(input string, spellchecker spellcheck.SpellChecker) (string, error) {
	wordsSpellChecked := spellchecker.SpellCheckString(deleteAllPunctuation(input))
	wordsWithoutStopWords := stemmer.deleteStopWords(wordsSpellChecked)
	resString, err := stemmer.normalizeWords(wordsWithoutStopWords)
	if err != nil {
		return "", fmt.Errorf("error normalize string: %w", err)
	}
	return strings.Join(resString, " "), nil
}

func NewSnowballStemmer(stopWordsPath string, langForStopWords []ISOCode639_1) (Stemmer, error) {
	availableLanguages := []ISOCode639_1{"en", "ru"} // not need more for now

	err := func() error {
		for _, l := range langForStopWords {
			if !slices.Contains(availableLanguages, l) {
				return fmt.Errorf("not found necessary stopwords for language %v", l)
			}
		}
		return nil
	}()
	if err != nil {
		return nil, fmt.Errorf("error find stopwords data: %w", err)
	}

	jsonFile, err := os.Open(stopWordsPath)

	if err != nil {
		return nil, fmt.Errorf("stopwords file not found: %w", err)
	}
	defer jsonFile.Close()

	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("error reading stopwords file: %w", err)
	}

	stopWordMapping := make(map[ISOCode639_1][]string)

	err = json.Unmarshal(jsonData, &stopWordMapping)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling stopwords: %w", err)
	}

	necessaryStopWords := make(map[ISOCode639_1][]string, len(langForStopWords))
	for _, langISO639 := range langForStopWords {
		necessaryStopWords[langISO639] = stopWordMapping[langISO639]
	}

	return &snowballStemmer{stopWords: necessaryStopWords, availableLanguages: availableLanguages}, nil
}

func deleteAllPunctuation(input string) string {

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
