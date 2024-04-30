package http

import (
	"context"
	"github.com/PavlushaSource/yadro-practice-course/internal/core/port"
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

func (ch *ComixHandler) GetRelevantURLIndex(ctx context.Context) {
	return
}

func (ch *ComixHandler) GetRelevantURL(ctx context.Context) {
	return
}

func (ch *ComixHandler) SayHello(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	}
}
