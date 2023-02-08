package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/michaelpeterswa/letterkennyapi/internal/handlers"
	"github.com/michaelpeterswa/letterkennyapi/internal/logging"
	"github.com/michaelpeterswa/letterkennyapi/internal/quotes"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func main() {
	logger, err := logging.InitZap()
	if err != nil {
		log.Panicf("could not acquire zap logger: %s", err.Error())
	}
	logger.Info("letterkennyapi init...")

	hostname, err := os.Hostname()
	if err != nil {
		logger.Fatal("could not get hostname", zap.Error(err))
	}

	// metrics and healthcheck
	go func() {
		internalRouter := mux.NewRouter()
		internalRouter.HandleFunc("/healthcheck", handlers.HealthcheckHandler)
		internalRouter.Handle("/metrics", promhttp.Handler())
		err = http.ListenAndServe(":8081", internalRouter)
		if err != nil {
			logger.Fatal("could not start metrics http server", zap.Error(err))
		}
	}()

	// main router
	mainRouter := mux.NewRouter()
	apiRouter := mainRouter.PathPrefix("/api").Subrouter()
	v1Router := apiRouter.PathPrefix("/v1").Subrouter()
	v1Router.HandleFunc("/quotes", handlers.NewQuoteHandler(quotes.LetterkennyQuotes).Handle)
	mainRouter.HandleFunc("/", handlers.NewHomeHandler(hostname).Handle)
	err = http.ListenAndServe(":8080", mainRouter)
	if err != nil {
		logger.Fatal("could not start http server", zap.Error(err))
	}
}
