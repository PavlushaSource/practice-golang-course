package main

import (
	"flag"
	"fmt"
	"github.com/kljensen/snowball"
	"io"
	"os"
	"slices"
	"strings"
)

func deleteStopWords(currentString string) []string {
	f, err := os.Open("stopwords_en.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println("Error closing file:", err)
		}
	}(f)

	stopWords, err := io.ReadAll(f)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}
	stopWordString := strings.Split(string(stopWords), "\n")
	res := make([]string, 0, len(stopWordString))

	for _, word := range strings.Split(currentString, " ") {
		word = strings.ToLower(word)
		if !slices.Contains(stopWordString, word) {
			res = append(res, word)
		}
	}
	return res
}

func stemmingWords(currentString []string) []string {
	res := make([]string, 0, len(currentString))
	for i := range currentString {
		stemWord, err := snowball.Stem(currentString[i], "english", true)
		if err == nil {
			res = append(res, stemWord)
		}
	}
	return res
}

func main() {

	var stemmerFlag bool
	flag.BoolVar(&stemmerFlag, "s", false,
		"A flag for normalizing the sentence that will be passed to the input program. Enter what needs to be normalized after the flag to work correctly.")

	flag.Parse()

	if stemmerFlag {
		fmt.Println(stemmingWords(deleteStopWords(flag.Args()[0])))
	}
}
