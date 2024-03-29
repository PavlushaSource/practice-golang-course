package main

import (
	"testing"
)

func BenchmarkBuilder(b *testing.B) {
	currentString := "Haha, you don't know which faster?!"
	for i := 0; i < 10; i++ {
		currentString += currentString
	}
	for i := 0; i < b.N; i++ {
		deleteAllPunctuationWithBuilder(currentString)
	}
}

func BenchmarkFields(b *testing.B) {
	currentString := "Haha, you don't know which faster?!"
	for i := 0; i < 10; i++ {
		currentString += currentString
	}
	for i := 0; i < b.N; i++ {
		deleteAllPunctuationWithField(currentString)
	}
}
