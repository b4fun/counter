package api

import (
	"fmt"
	"net/http"
)

func HandleConsulHealthCheck(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "ok, consul")
}
