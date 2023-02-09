package handlers

import (
	"encoding/json"
	"net/http"

	"676f.dev/rfc7807"
	"github.com/michaelpeterswa/letterkennyapi/internal/quotes"
	"go.uber.org/zap"
)

type quoteResponse struct {
	Quote string `json:"quote"`
}

type QuoteHandler struct {
	Logger *zap.Logger
}

func NewQuoteHandler(logger *zap.Logger, quotes []string) *QuoteHandler {
	return &QuoteHandler{Logger: logger}
}

func (qh *QuoteHandler) getQuoteResponse() quoteResponse {
	return quoteResponse{Quote: quotes.GetRandomQuote()}
}

func (qh *QuoteHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	err := json.NewEncoder(w).Encode(qh.getQuoteResponse())
	if err != nil {
		qh.Logger.Error("failed to encode quote", zap.Error(err))
		err = rfc7807.SimpleResponse(w, http.StatusInternalServerError, "failed to encode quote")
		if err != nil {
			qh.Logger.Error("failed to send error", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
