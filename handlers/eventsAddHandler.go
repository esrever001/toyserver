package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/esrever001/toyserver/db"
	"github.com/julienschmidt/httprouter"
)

type EventsAddRequest struct {
	User  string
	Type  string
	Time  *time.Time
	Notes string
	Image string
}

type EventsAddHandler struct {
	Database *db.Database
}

func (handler EventsAddHandler) Path() string {
	return "/events/add"
}

func (handler EventsAddHandler) Method() HttpMethod {
	return POST
}

func (handler EventsAddHandler) getEvent(request EventsAddRequest) db.Events {
	timeFromRequest := time.Now().Unix()
	if request.Time != nil {
		timeFromRequest = (*request.Time).Unix()
	}
	return db.Events{
		User:      request.User,
		Type:      request.Type,
		Timestamp: timeFromRequest,
		Notes:     request.Notes,
		Image:     request.Image,
	}
}

func (handler EventsAddHandler) Handle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	decoder := json.NewDecoder(r.Body)

	var requestBody EventsAddRequest
	err := decoder.Decode(&requestBody)
	if err != nil {
		panic(err)
	}
	event := handler.getEvent(requestBody)

	handler.Database.Database.Create(&event)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
	fmt.Printf("Adding event for user %s\n", ps.ByName("user"))
}
