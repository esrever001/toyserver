package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/esrever001/toyserver/db"
	"github.com/julienschmidt/httprouter"
)

type EventsAddHandler struct {
	Database *db.Database
}

func (handler EventsAddHandler) Path() string {
	return "/events/add/:user/:type"
}

func (handler EventsAddHandler) Handle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	event := db.Events{
		User: ps.ByName("user"),
		Type: ps.ByName("type"),
		Time: time.Now(),
	}
	handler.Database.Database.Create(&event)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
	fmt.Printf("Adding event for user %s\n", ps.ByName("user"))
}
