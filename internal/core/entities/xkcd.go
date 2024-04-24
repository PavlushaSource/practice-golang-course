package entities

import (
	"fmt"
	"github.com/PavlushaSource/yadro-practice-course/pkg/words/stemmer"
	"sync"
)

type ComicInfo struct {
	Month      string `json:"month"`
	Num        int    `json:"num"`
	Link       string `json:"link"`
	Year       string `json:"year"`
	News       string `json:"news"`
	SafeTitle  string `json:"safe_title"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
	Img        string `json:"img"`
	Title      string `json:"title"`
	Day        string `json:"day"`
}

type ComicToJSON struct {
	URL      string   `json:"url"`
	Keywords []string `json:"keywords"`
}

var myStemmer stemmer.Stemmer
var onceStemmer sync.Once

func (c *ComicInfo) ToJSON() (*ComicToJSON, error) {
	comicToJSON := &ComicToJSON{
		URL: c.Img,
	}
	var err error

	onceStemmer.Do(func() {
		myStemmer, err = stemmer.NewSnowballStemmer()
		if err != nil {
			panic(err)
		}
	})

	if err != nil {
		return nil, fmt.Errorf("error create stemmer: %w", err)
	}

	safeTitle, err := myStemmer.NormalizeString(c.SafeTitle)
	if err != nil {
		return nil, fmt.Errorf("error normalize safe title: %w", err)
	}

	transcript, err := myStemmer.NormalizeString(c.Transcript)
	if err != nil {
		return nil, fmt.Errorf("error normalize transcript: %w", err)
	}

	alt, err := myStemmer.NormalizeString(c.Alt)
	if err != nil {
		return nil, fmt.Errorf("error normalize save title: %w", err)
	}

	comicToJSON.Keywords = append(comicToJSON.Keywords, safeTitle...)
	comicToJSON.Keywords = append(comicToJSON.Keywords, transcript...)
	comicToJSON.Keywords = append(comicToJSON.Keywords, alt...)

	return comicToJSON, nil
}
