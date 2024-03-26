package main

import (
	"flag"
	"fmt"
	"github.com/PavlushaSource/yadro-practice-course/spellcheck"
	"os"
)

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	var stemmerFlag bool
	flag.BoolVar(&stemmerFlag, "s", false,
		"A flag for normalizing the sentence that will be passed to the input program. Enter string after the flag to work correctly.")

	flag.Parse()

	if stemmerFlag {
		args := flag.Args()
		if len(args) != 1 {
			fmt.Println("Stemmer work with one argument - string.\nExample: ./myApp -s \"current string\"")
			return
		}
		langStopWords := []ISOCode639_1{"en", "ru"}
		stopWordsDatapath := "stopwords-iso.json"
		st, err := NewSnowballStemmer(stopWordsDatapath, langStopWords)
		check(err)
		checker := spellcheck.NewFuzzyChecker(1, 2)
		err = checker.LoadDataset("spellcheck/all-words.txt")
		check(err)
		normalizedWithSpellchecker, err := st.NormalizeStringWithSpellcheck(args[0], checker)
		check(err)
		normalized, err := st.NormalizeString(args[0])
		check(err)

		fmt.Printf("result with spellcheck - %s\n", normalizedWithSpellchecker)
		fmt.Printf("result without spellcheck - %s\n", normalized)
	}
}
