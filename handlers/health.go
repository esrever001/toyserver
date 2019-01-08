package handlers

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type HealthHandler struct{}

func (handler HealthHandler) Path() string {
	return "/health"
}

func (handler HealthHandler) Method() HttpMethod {
	return GET
}

func (handler HealthHandler) Handle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, "Welcome!\n")
	fmt.Printf("OK")
}
