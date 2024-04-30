package xkcd

import "net/http"

func NewClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{},
	}
}
