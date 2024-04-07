package stemmer

import (
	"encoding/json"
	"fmt"
	"github.com/PavlushaSource/yadro-practice-course/pkg/words/spellcheck"
	"github.com/kljensen/snowball"
	"github.com/pemistahl/lingua-go"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
)

type ISOCode639_1 string

type Stemmer interface {
	NormalizeString(input string, spellcheck ...spellcheck.SpellChecker) ([]string, error)
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
		// try convert to number
		if _, err := strconv.Atoi(word); err == nil {
			res = append(res, word)
			continue
		}

		if usedLanguage, exist := detector.DetectLanguageOf(word); exist {
			stemWord, err := snowball.Stem(word, strings.ToLower(usedLanguage.String()), true)
			if err != nil {
				return nil, fmt.Errorf("stemming error: %w", err)
			}
			res = append(res, stemWord)
		} else {
			// try convert to number and resolve special cases (2007-07-01)
			specialWords := strings.Split(word, "-")
			for _, specialWord := range specialWords {
				specialWord = deleteAllPunctuationWithBuilder(specialWord)
				if _, err := strconv.Atoi(specialWord); err == nil {
					res = append(res, specialWord)
					continue
				}
				// uncomment if you want to see the words that stemmer skips
				//if specialWord != "" {
				//	fmt.Println("language for word: ", word, "is not detected. Skip this word.")
				//}
			}
		}
	}
	return slices.Clip(res), nil
}

func removeDuplicateStrings(s []string) []string {
	stringWithoutDuplicate := make([]string, 0)
	exist := make(map[string]struct{})
	for _, currS := range s {
		if _, ok := exist[currS]; !ok {
			stringWithoutDuplicate = append(stringWithoutDuplicate, currS)
			exist[currS] = struct{}{}
		}
	}
	return stringWithoutDuplicate
}

func (stemmer *snowballStemmer) NormalizeString(input string, spellchecker ...spellcheck.SpellChecker) ([]string, error) {
	if len(spellchecker) > 0 {
		input = spellchecker[0].SpellCheckString(input)
	}
	wordsWithoutStopWords := stemmer.deleteStopWords(deleteAllPunctuationWithBuilder(input))
	resString, err := stemmer.normalizeWords(wordsWithoutStopWords)
	if err != nil {
		return nil, fmt.Errorf("error normalize string: %w", err)
	}
	resString = removeDuplicateStrings(resString)
	return resString, nil
}

func NewSnowballStemmer(stopWordsPath ...string) (Stemmer, error) {
	availableLanguages := []ISOCode639_1{"en", "ru"} // not need more for now
	var currentStopWordsPath string

	if len(stopWordsPath) > 0 {
		currentStopWordsPath = stopWordsPath[0]
	} else {
		currentStopWordsPath = "internal/resources/words/stopwords/stopwords-iso.json"
	}

	jsonFile, err := os.Open(currentStopWordsPath)

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

	necessaryStopWords := make(map[ISOCode639_1][]string, len(availableLanguages))
	for _, langISO639 := range availableLanguages {
		necessaryStopWords[langISO639] = stopWordMapping[langISO639]
	}

	return &snowballStemmer{stopWords: necessaryStopWords, availableLanguages: availableLanguages}, nil
}
