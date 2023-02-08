package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"

	"676f.dev/rfc7807"
)

type quoteResponse struct {
	Quote string `json:"quote"`
}

type QuoteHandler struct {
	Quotes []string
}

func NewQuoteHandler(quotes []string) *QuoteHandler {
	return &QuoteHandler{Quotes: quotes}
}

func (qh *QuoteHandler) getQuote() quoteResponse {
	return quoteResponse{Quote: qh.Quotes[rand.Intn(len(qh.Quotes))]}
}

func (qh *QuoteHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(qh.getQuote())
	if err != nil {
		rfc7807.SimpleResponse(w, http.StatusInternalServerError, "failed to encode quote")
	}
}
