package http

import (
	"fmt"
	"net/http"
	"strconv"
)

type SuggestRelevantURLRequest struct {
	search string
	limit  int
}

func newSuggestRelevantURLRequest(r *http.Request) SuggestRelevantURLRequest {
	req := SuggestRelevantURLRequest{limit: 10}
	var err error

	search := r.URL.Query().Get("search")
	if search == "" {
		return req
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		fmt.Printf("cannot convert")
		limit = 10
	}
	fmt.Printf("length: %d\n", limit)
	return SuggestRelevantURLRequest{
		search: search,
		limit:  limit,
	}
}
