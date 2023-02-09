package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"676f.dev/utilities/tools/shortid"
	"github.com/gorilla/mux"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/env"
	"github.com/michaelpeterswa/letterkennyapi/internal/handlers"
	"github.com/michaelpeterswa/letterkennyapi/internal/logging"
	"github.com/michaelpeterswa/letterkennyapi/internal/quotes"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

var k = koanf.New(".")

func main() {
	logger, err := logging.InitZap()
	if err != nil {
		log.Panicf("could not acquire zap logger: %s", err.Error())
	}

	k.Load(env.Provider("LKAPI_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "LKAPI_")), "_", ".", -1)
	}), nil)

	sid := shortid.NewShortID(shortid.Base58CharacterSet)
	instanceID, err := sid.Generate(6)
	if err != nil {
		logger.Fatal("could not generate instance id", zap.Error(err))
	}

	hostname, err := os.Hostname()
	if err != nil {
		logger.Fatal("could not get hostname", zap.Error(err))
	}

	logger.Info("starting letterkennyapi", zap.String("instance_id", instanceID), zap.String("hostname", hostname))

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

	// -=-=-=-=-=-=-
	homeTemplate, err := template.ParseFiles(k.String("templates.home"))
	if err != nil {
		logger.Fatal("could not parse home template", zap.Error(err))
	}
	// -=-=-=-=-=-=-

	// main router
	mainRouter := mux.NewRouter()
	apiRouter := mainRouter.PathPrefix("/api").Subrouter()
	v1Router := apiRouter.PathPrefix("/v1").Subrouter()
	v1Router.HandleFunc("/quotes", handlers.NewQuoteHandler(logger, quotes.LetterkennyQuotes).Handle)
	mainRouter.HandleFunc("/", handlers.NewHomeHandler(logger, instanceID, k.String("title"), k.String("productionurl"), homeTemplate).Handle)
	err = http.ListenAndServe(":8080", mainRouter)
	if err != nil {
		logger.Fatal("could not start http server", zap.Error(err))
	}
}
