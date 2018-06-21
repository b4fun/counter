package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/b4fun/counter"
	"github.com/b4fun/counter/api"
	"github.com/b4fun/counter/inmem"
	credis "github.com/b4fun/counter/redis"
	"github.com/gomodule/redigo/redis"
)

func main() {
	var inmemCounter counter.Counter
	{
		inmemCounter = inmem.New()
	}

	var redisCounter counter.Counter
	{
		redisPool := &redis.Pool{
			MaxIdle:     3,
			IdleTimeout: 10 * time.Second,
			Dial: func() (redis.Conn, error) {
				// FIXME support get host from env var
				return redis.Dial("tcp", "counter-redis:6379")
			},
		}
		redisCounter = credis.New(redisPool)
	}

	apiServer := api.NewServer(api.NewServerOpt{
		Counters: map[string]counter.Counter{
			"inmem": inmemCounter,
			"redis": redisCounter,
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
