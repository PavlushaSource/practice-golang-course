package http

import (
	"context"
	"encoding/json"
	"github.com/PavlushaSource/yadro-practice-course/internal/core/port"
	"log/slog"
	"net/http"
)

type ComicHandler struct {
	svc port.ComicsService
}

func NewComicHandler(svc port.ComicsService) *ComicHandler {
	return &ComicHandler{svc: svc}
}

func (ch *ComicHandler) UpdateHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		oldComics, err := ch.svc.GetComics()
		lastCount := len(oldComics)
		if err != nil {
			slog.Error("GetComics failed", "error", err)
			errorResponse(w, http.StatusInternalServerError, err)
			return
		}

		newComics, err := ch.svc.DownloadAll(ctx)

		if err != nil {
			slog.Error("DownloadAll failed", "error", err)
			errorResponse(w, http.StatusInternalServerError, err)
			return
		}
		updatedCount := len(newComics)
		rsp := newUpdateResponse(updatedCount, lastCount+updatedCount)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(rsp)
	}
}

func (ch *ComicHandler) GetSuggestRelevantURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		req := newSuggestRelevantURLRequest(r)

		comics, err := ch.svc.GetRelevantComics(req.search, req.limit)
		if err != nil {
			slog.Error("GetRelevantComics failed", "error", err)
			errorResponse(w, http.StatusInternalServerError, err)
			return
		}

		rsp := newComicsResponse(comics)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(rsp)
	}

}

func (ch *ComicHandler) SayHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello\n"))
	}
}
