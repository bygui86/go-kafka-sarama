package monitoring

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/bygui86/go-kafka/tracing/kafka-producer/commons"
	"github.com/bygui86/go-kafka/tracing/kafka-producer/logging"
)

func (s *Server) newRouter() {
	logging.SugaredLog.Debugf("Setup new monitoring router")

	s.router = mux.NewRouter().StrictSlash(true)

	s.router.Handle("/metrics", promhttp.Handler())
}

func (s *Server) newHTTPServer() {
	logging.SugaredLog.Debugf("Setup new monitoring HTTP server on port %d...", s.config.restPort)

	if s.config != nil {
		s.httpServer = &http.Server{
			Addr:    fmt.Sprintf(commons.HttpServerHostFormat, s.config.restHost, s.config.restPort),
			Handler: s.router,
			// Good practice to set timeouts to avoid Slowloris attacks.
			WriteTimeout: commons.HttpServerWriteTimeoutDefault,
			ReadTimeout:  commons.HttpServerReadTimeoutDefault,
			IdleTimeout:  commons.HttpServerIdelTimeoutDefault,
		}
		return
	}

	logging.Log.Error("Monitoring HTTP server creation failed: monitoring configurations not loaded")
}
