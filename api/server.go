package api

import (
	"net/http"

	"github.com/b4fun/counter"
)

type NewServerOpt struct {
	Counters map[string]counter.Counter
}

// NewServer creates a api server.
func NewServer(opt NewServerOpt) http.Handler {
	server := http.NewServeMux()

	for namespace, impl := range opt.Counters {
		BindCounterHandler(server, namespace, impl)
	}

	return server
}
