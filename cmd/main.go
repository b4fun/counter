package main

import (
	"fmt"
	"net/http"

	"github.com/b4fun/counter"
	"github.com/b4fun/counter/api"
	"github.com/b4fun/counter/inmem"
)

func main() {
	inmemCounter := inmem.New()

	apiServer := api.NewServer(api.NewServerOpt{
		Counters: map[string]counter.Counter{
			"inmem": inmemCounter,
		},
	})

	components := []startFunc{
		func(errs chan error) {
			errs <- http.ListenAndServe(":8081", apiServer)
		},
		startCounter(inmemCounter),
	}
	errs := make(chan error, len(components))
	for _, c := range components {
		go c(errs)
	}

	select {
	case e := <-errs:
		fmt.Printf("err: %v\n", e)
		return
	}
}

type startFunc func(chan error)

type startableCounter interface {
	Start() error
}

func startCounter(impl counter.Counter) startFunc {
	return func(errs chan error) {
		if impl, ok := impl.(startableCounter); ok {
			errs <- impl.Start()
		}
	}
}
