package xkcd

import (
	"encoding/json"
	"fmt"
	"github.com/PavlushaSource/yadro-practice-course/internal/core/entities"
	"net/http"
)

const (
	xkcd = "https://xkcd.com"
)

func NewClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{},
	}
}

func generateComicURL(id int, siteURL string) string {
	switch siteURL {
	case xkcd:
		return fmt.Sprintf("%s/%d/info.0.json", siteURL, id)
	default:
		return ""
	}
}

func GetComicByID(client *http.Client, ID int) (*entities.ComicInfo, error) {
	urlName := generateComicURL(ID, xkcd)
	resp, err := client.Get(urlName)
	if err != nil {
		return nil, fmt.Errorf("cannot get comic from url %s: %w", urlName, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cannot get comic from url %s: status code=%v", urlName, resp.StatusCode)
	}

	var comic entities.ComicInfo
	err = json.NewDecoder(resp.Body).Decode(&comic)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal comic from url %s: %w", urlName, err)
	}
	return &comic, nil
}
