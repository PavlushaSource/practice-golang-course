package stemmer

import (
	"slices"
	"strings"
)

func deleteAllPunctuation(input string) string {
	// punctuations without \' and -
	punctuations := []rune{
		'!', '?', '.', ',', ';', ':', '\'', '"', '@', '&', '#', '$', '%', '^', '*', '(', ')', '[', ']',
		'{', '}', '<', '>', '/', '|', '\\', '`', '~', '=', '+', '_',
	}

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
