package main

import (
	"fmt"
	"github.com/PavlushaSource/yadro-practice-course/pkg/words/spellcheck"
	"github.com/PavlushaSource/yadro-practice-course/pkg/words/stemmer"
	"time"
)

func main() {
	input, err := stemmer.GetStringFromCommandlineInput()
	stemmer.Check(err)
	start := time.Now()

	st, err := stemmer.NewSnowballStemmer()
	stemmer.Check(err)
	checker := spellcheck.NewFuzzyChecker()
	err = checker.SaveModel("internal/resources/spellchecker/savedModel")
	stemmer.Check(err)

	normalizedWithSpellchecker, err := st.NormalizeString(input, checker)
	stemmer.Check(err)
	normalized, err := st.NormalizeString(input)
	stemmer.Check(err)

	end := time.Now()
	fmt.Printf("duration: %s\n", end.Sub(start))
	fmt.Printf("result with spellcheck - %s\n", normalizedWithSpellchecker)
	fmt.Printf("result without spellcheck - %s\n", normalized)
}
