package http

import "github.com/PavlushaSource/yadro-practice-course/internal/core/domain"

type ComicResponse struct {
	ImgURL []SuggestImgURLResponse `json:"URL"`
}

func newComicsResponse(comics []domain.Comix) ComicResponse {
	return ComicResponse{ImgURL: newSuggestImgURLsResponse(comics)}
}

type SuggestImgURLResponse struct {
	Img string `json:"img"`
}

func newSuggestImgURLResponse(img string) SuggestImgURLResponse {
	return SuggestImgURLResponse{Img: img}
}

func newSuggestImgURLsResponse(comics []domain.Comix) []SuggestImgURLResponse {
	res := make([]SuggestImgURLResponse, 0, len(comics))
	for _, comic := range comics {
		res = append(res, newSuggestImgURLResponse(comic.URL))
	}
	return res
}
