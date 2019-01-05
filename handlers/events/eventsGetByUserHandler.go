package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/esrever001/toyserver/db"
	"github.com/julienschmidt/httprouter"
)

type EventsGetByUserHandler struct {
	Database *db.Database
}

func (handler EventsGetByUserHandler) Path() string {
	return "/events/get/:user"
}

func (handler EventsGetByUserHandler) Handle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var events []db.Events
	handler.Database.Database.Where("user = ?", ps.ByName("user")).Find(&events)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
	fmt.Printf("Getting events for user %s\n", ps.ByName("user"))
}
