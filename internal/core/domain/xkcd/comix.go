package xkcd

import (
	"github.com/PavlushaSource/yadro-practice-course/internal/core/domain"
	"github.com/PavlushaSource/yadro-practice-course/internal/core/port"
	"strings"
)

type Comix struct {
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

func (c *Comix) Format(normalizeSvc port.NormalizeService) domain.Comix {
	var comix domain.Comix

	safeTitle, err := normalizeSvc.Normalize(c.SafeTitle)
	if err != nil {
		safeTitle = strings.Split(c.SafeTitle, " ")
	}

	transcript, err := normalizeSvc.Normalize(c.Transcript)
	if err != nil {
		transcript = strings.Split(c.Transcript, " ")
	}

	alt, err := normalizeSvc.Normalize(c.Alt)
	if err != nil {
		alt = strings.Split(c.Alt, " ")
	}

	comix.Keywords = append(comix.Keywords, safeTitle...)
	comix.Keywords = append(comix.Keywords, transcript...)
	comix.Keywords = append(comix.Keywords, alt...)

	comix.ID = uint64(c.Num)
	comix.URL = c.Img

	return comix
}
