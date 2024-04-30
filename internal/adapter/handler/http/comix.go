package http

import (
	"context"
	"github.com/PavlushaSource/yadro-practice-course/internal/core/port"
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
