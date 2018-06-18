package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/b4fun/counter"
)

func abortWithError(w http.ResponseWriter, err error) {
	w.WriteHeader(500)
	w.Write([]byte(fmt.Sprintf("err: %v", err)))

}

// HTTPRouter represents a http router.
type HTTPRouter interface {
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}

// BindCounterHandler binds a counter to api under given namespace.
func BindCounterHandler(
	router HTTPRouter,
	namespace string,
	impl counter.Counter,
) {
	const counterIdField = "counter_id"

	router.HandleFunc(
		fmt.Sprintf("/v1/%s/get", namespace),
		func(w http.ResponseWriter, req *http.Request) {
			if req.Method != "GET" {
				w.WriteHeader(405)
				return
			}

			id := req.URL.Query().Get(counterIdField)
			if id == "" {
				w.WriteHeader(400)
				return
			}

			counter, err := impl.Get(req.Context(), id)
			if err != nil {
				abortWithError(w, err)
				return
			}

			output := json.NewEncoder(w)
			if err := output.Encode(counter); err != nil {
				abortWithError(w, err)
				return
			}

			w.Header().Set("content-type", "application/json")
		},
	)

	router.HandleFunc(
		fmt.Sprintf("/v1/%s/incr", namespace),
		func(w http.ResponseWriter, req *http.Request) {
			if req.Method != "POST" {
				w.WriteHeader(405)
				return
			}

			id := req.URL.Query().Get(counterIdField)
			if id == "" {
				w.WriteHeader(400)
				return
			}

			if err := impl.Incr(req.Context(), id); err != nil {
				abortWithError(w, err)
				return
			}

			w.WriteHeader(204)
		},
	)
}
