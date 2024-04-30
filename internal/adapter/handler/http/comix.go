package http

import (
	"context"
	"encoding/json"
	"github.com/PavlushaSource/yadro-practice-course/internal/core/port"
	"log/slog"
	"net/http"
)

type ComixHandler struct {
	svc port.ComixService
}

func NewComixHandler(svc port.ComixService) *ComixHandler {
	return &ComixHandler{svc: svc}
}

func (ch *ComixHandler) Update(ctx context.Context) {
	return
}

func (ch *ComixHandler) GetSuggestRelevantURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		req := newSuggestRelevantURLRequest(r)

		comics, err := ch.svc.GetRelevantComics(req.search, req.limit)
		if err != nil {
			slog.Error("GetRelevantComics failed", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err)
			return
		}

		rsp := newComicsResponse(comics)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(rsp)
	}

}

func (ch *ComixHandler) SayHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello\n"))
	}
}
