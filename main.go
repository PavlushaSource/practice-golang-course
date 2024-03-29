package main

import (
	"fmt"
	"github.com/PavlushaSource/yadro-practice-course/spellcheck"
	"time"
)

func main() {
	input, err := getStringFromCommandlineInput()
	check(err)
	start := time.Now()

	st, err := NewSnowballStemmer()
	check(err)
	checker := spellcheck.NewFuzzyChecker()
	err = checker.SaveModel("spellcheck/savedModel")
	check(err)

	normalizedWithSpellchecker, err := st.NormalizeString(input, checker)
	check(err)
	normalized, err := st.NormalizeString(input)
	check(err)

	end := time.Now()
	fmt.Printf("duration: %s\n", end.Sub(start))
	fmt.Printf("result with spellcheck - %s\n", normalizedWithSpellchecker)
	fmt.Printf("result without spellcheck - %s\n", normalized)
}
