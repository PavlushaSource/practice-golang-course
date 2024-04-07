package xkcd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
)

var errStatus = errors.New("response status error")

const (
	xkcd = "https://xkcd.com"
)

type ComicInfo struct {
	Num        int    `json:"num"`
	Transcript string `json:"transcript"`
	Img        string `json:"img"`
}

func NewClient() *http.Client {
	// not to load the connection too much
	return &http.Client{
		Transport: &http.Transport{MaxConnsPerHost: 100},
	}
}

type StatusError struct {
	status string
}

func (e *StatusError) Error() string {
	return fmt.Sprintf("status error:  response status code=%s", e.status)
}

func findNumberComics(client http.Client, urlName string) (int, error) {
	switch urlName {
	case xkcd:
		comic, err := getComicFromURL(&client, fmt.Sprintf("%s/info.0.json", urlName))
		if err != nil {
			return 0, fmt.Errorf("cannot find number of comics: %w", err)
		}
		return comic.Num, nil
	default:
		return 0, fmt.Errorf("cannot find number of comics from url: %s", urlName)
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

func GetComics(client *http.Client, urlName string, log *slog.Logger, numbersToGet ...int) map[int]ComicInfo {
	log = log.With("GetComicsFromSite", urlName)
	var lastNumberComic int
	var err error

	if len(numbersToGet) > 0 {
		lastNumberComic = numbersToGet[0]
	} else {
		lastNumberComic, err = findNumberComics(*client, urlName)
		if err != nil {
			log.Error(err.Error())
			return nil
		}
	}

	readComics := make(map[int]ComicInfo)
	mapMutex := sync.Mutex{}
	wg := sync.WaitGroup{}

	for i := 1; i <= lastNumberComic; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			comicURL := generateComicURL(i, urlName)
			comic, err := getComicFromURL(client, comicURL)
			switch {
			case errors.Is(err, errStatus):
				log.Debug(err.Error(), "comicID", i)
				return
			case err != nil:
				log.Error(err.Error(), "comicID", i)
				return
			}
			mapMutex.Lock()
			readComics[i] = *comic
			mapMutex.Unlock()
		}(i)
	}
	wg.Wait()
	return readComics
}

func getComicFromURL(client *http.Client, urlName string) (*ComicInfo, error) {
	resp, err := client.Get(urlName)
	if err != nil {
		return nil, fmt.Errorf("cannot get comic from url %s: %w", urlName, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read comic from url %s: %w", urlName, err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("cannot get comic from url %s: %w=%d", urlName, errStatus, resp.StatusCode)
	}

	var comic ComicInfo
	err = json.Unmarshal(body, &comic)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal comic from url %s: %w", urlName, err)
	}
	return &comic, nil
}
