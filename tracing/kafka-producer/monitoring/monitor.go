package monitoring

import (
	"context"
	"time"

	"github.com/bygui86/go-kafka/tracing/kafka-producer/logging"
)

func New() *Server {
	logging.Log.Info("Create new monitoring server")

	cfg := loadConfig()
	server := &Server{
		config: cfg,
	}
	server.newRouter()
	server.newHTTPServer()
	return server
}

func (s *Server) Start() {
	logging.Log.Info("Start monitoring server")

	if s.httpServer != nil && !s.running {
		go func() {
			err := s.httpServer.ListenAndServe()
			if err != nil {
				logging.SugaredLog.Errorf("Monitoring server start failed: %s", err.Error())
			}
		}()
		s.running = true
		logging.SugaredLog.Infof("Monitoring server listen on port %d", s.config.restPort)
		return
	}

	logging.Log.Error("Monitoring server start failed: HTTP server not initialized or HTTP server already running")
}

func (s *Server) Shutdown(timeout int) {
	logging.SugaredLog.Warnf("Shutdown monitoring server, timeout %d", timeout)

	if s.httpServer != nil && s.running {
		// create a deadline to wait for.
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
		defer cancel()
		// does not block if no connections, otherwise wait until the timeout deadline
		err := s.httpServer.Shutdown(ctx)
		if err != nil {
			logging.SugaredLog.Errorf("Monitoring server shutdown failed: %s", err.Error())
		}
		s.running = false
		return
	}

	logging.Log.Error("Monitoring server shutdown failed: HTTP server not initialized or HTTP server not running")
}
