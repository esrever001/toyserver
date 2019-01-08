package handlers

import (
	"net/http"

	"github.com/esrever001/toyserver/db"
	"github.com/julienschmidt/httprouter"
)

type HttpMethod int

const (
	GET HttpMethod = 1 + iota
	POST
)

type BaseHandler interface {
	Path() string
	Method() HttpMethod
	Handle(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
}

func CreateHandlers(db *db.Database) []BaseHandler {
	return []BaseHandler{
		HealthHandler{},
		EventsAddHandler{Database: db},
		EventsGetByUserHandler{Database: db},
	}
}
