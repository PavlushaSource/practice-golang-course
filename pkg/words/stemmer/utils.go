package stemmer

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
)

func GetStringFromCommandlineInput() (string, error) {
	var input string
	flag.StringVar(&input, "s", "", "A flag for normalizing the "+
		"sentence that will be passed to the input program. Enter string after the flag to work correctly.")
	flag.Parse()
	if otherInput := flag.Args(); len(otherInput) > 0 || input == "" {
		return "", fmt.Errorf("Stemmer work with one argument - string.\nExample: ./myApp -s \"current string\"")
	}
	return input, nil
}

func Check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func deleteAllPunctuationWithBuilder(input string) string {
	// punctuations without \' and -
	punctuations := []rune{'!', '?', '.', ',', ';', ':', '\'', '"', '@', '&', '#', '$', '%', '^', '*', '(', ')', '[', ']',
		'{', '}', '<', '>', '/', '|', '\\', '`', '~', '=', '+', '_'}

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
