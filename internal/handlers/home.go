package handlers

import (
	"html/template"
	"net/http"

	"676f.dev/rfc7807"
	"go.uber.org/zap"
)

type HomeHandler struct {
	Logger        *zap.Logger
	InstanceID    string
	Title         string
	ProductionURL string
	HomeTemplate  *template.Template
}

type HomeData struct {
	InstanceID string
	URL        string
	Title      string
}

func NewHomeHandler(logger *zap.Logger, instanceID string, title string, produrl string, homeTemplate *template.Template) *HomeHandler {
	return &HomeHandler{InstanceID: instanceID, HomeTemplate: homeTemplate, Title: title, ProductionURL: produrl, Logger: logger}
}

func (hh *HomeHandler) Handle(w http.ResponseWriter, r *http.Request) {
	homeData := HomeData{InstanceID: hh.InstanceID, URL: hh.ProductionURL, Title: hh.Title}
	err := hh.HomeTemplate.Execute(w, homeData)
	if err != nil {
		err := rfc7807.SimpleResponse(w, http.StatusInternalServerError, "could not render home template")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			hh.Logger.Error("could not render home template", zap.Error(err))
		}
	}
}
