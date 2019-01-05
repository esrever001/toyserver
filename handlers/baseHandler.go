package handlers

import (
	"net/http"

	"github.com/esrever001/toyserver/db"
	eventsHandler "github.com/esrever001/toyserver/handlers/events"
	"github.com/julienschmidt/httprouter"
)

type BaseHandler interface {
	Path() string
	Handle(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
}

func CreateHandlers(db *db.Database) []BaseHandler {
	return []BaseHandler{
		eventsHandler.EventsAddHandler{Database: db},
		eventsHandler.EventsGetByUserHandler{Database: db},
	}
}
