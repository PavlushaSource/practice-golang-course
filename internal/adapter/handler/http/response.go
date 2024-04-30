package http

import (
	"encoding/json"
	"github.com/PavlushaSource/yadro-practice-course/internal/core/domain"
	"net/http"
)

type UpdateResponse struct {
	NewComicsCount   int `json:"NewComicsCount"`
	TotalComicsCount int `json:"TotalComicsCount"`
}

func newUpdateResponse(newComicsCount int, totalComicsCount int) UpdateResponse {
	return UpdateResponse{NewComicsCount: newComicsCount, TotalComicsCount: totalComicsCount}
}

type ComicResponse struct {
	ImgURL []SuggestImgURLResponse `json:"SuggestedComics"`
}

func newComicsResponse(comics []domain.Comic) ComicResponse {
	return ComicResponse{ImgURL: newSuggestImgURLsResponse(comics)}
}

type SuggestImgURLResponse struct {
	Img string `json:"ImgURL"`
}

func newSuggestImgURLResponse(img string) SuggestImgURLResponse {
	return SuggestImgURLResponse{Img: img}
}

func newSuggestImgURLsResponse(comics []domain.Comic) []SuggestImgURLResponse {
	res := make([]SuggestImgURLResponse, 0, len(comics))
	for _, comic := range comics {
		res = append(res, newSuggestImgURLResponse(comic.URL))
	}
	return res
}

func errorResponse(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(err)
}
