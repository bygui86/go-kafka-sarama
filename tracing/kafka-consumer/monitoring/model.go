package monitoring

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Config struct {
	restHost string
	restPort int
}

type Server struct {
	config     *Config
	router     *mux.Router
	httpServer *http.Server
	running    bool
}
