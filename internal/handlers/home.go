package handlers

import (
	"encoding/json"
	"net/http"
)

type home struct {
	Hostname string `json:"hostname"`
}

type HomeHandler struct {
	Hostname string
}

func NewHomeHandler(hostname string) *HomeHandler {
	return &HomeHandler{Hostname: hostname}
}

func (hh *HomeHandler) toHome() home {
	return home{Hostname: hh.Hostname}
}

func (hh *HomeHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(hh.toHome())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
