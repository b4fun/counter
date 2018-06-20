package api

import (
	"net/http"

	"github.com/b4fun/counter"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

type NewServerOpt struct {
	Counters map[string]counter.Counter
}

// NewServer creates a api server.
func NewServer(opt NewServerOpt) http.Handler {
	mux := http.NewServeMux()
	for namespace, impl := range opt.Counters {
		BindCounterHandler(mux, namespace, impl)
	}

	httpLogger := negroni.NewLogger()
	httpLogger.ALogger = logrus.StandardLogger()

	n := negroni.New()
	n.Use(httpLogger)
	n.UseHandler(mux)

	return n
}
