package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"

	"676f.dev/rfc7807"
	"go.uber.org/zap"
)

type quoteResponse struct {
	Quote string `json:"quote"`
}

type QuoteHandler struct {
	Logger *zap.Logger
	Quotes []string
}

func NewQuoteHandler(logger *zap.Logger, quotes []string) *QuoteHandler {
	return &QuoteHandler{Logger: logger, Quotes: quotes}
}

func (qh *QuoteHandler) getQuote() quoteResponse {
	return quoteResponse{Quote: qh.Quotes[rand.Intn(len(qh.Quotes))]}
}

func (qh *QuoteHandler) Handle(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(qh.getQuote())
	if err != nil {
		qh.Logger.Error("failed to encode quote", zap.Error(err))
		err = rfc7807.SimpleResponse(w, http.StatusInternalServerError, "failed to encode quote")
		if err != nil {
			qh.Logger.Error("failed to send error", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
