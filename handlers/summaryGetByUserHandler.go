package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/esrever001/toyserver/db"
	"github.com/julienschmidt/httprouter"
)

type SummaryResponse struct {
	EventsByDate map[string][]string
	EventsByType map[string]*SummaryByType
}

type DateDetails struct {
	Year     int
	Month    int
	Day      int
	YearDay  int
	WeekYear int
	Week     int
}

type SummaryByType struct {
	StartingTime   int64
	MostRecentTime int64
	TotalDays      int
	TotalWeeks     int
	ContinuesDays  int
	ContinuesWeeks int
	GapDays        int
	Status         string
	Description    string

	DaysDetails   []DateDetails
	EventsDetails []db.Events
}

type SummaryGetByUserHandler struct {
	Database *db.Database
}

func (handler SummaryGetByUserHandler) Path() string {
	return "/events/summary/:user"
}

func (handler SummaryGetByUserHandler) Method() HttpMethod {
	return GET
}

func (handler SummaryGetByUserHandler) getDateDetail(events []db.Events) []DateDetails {
	daysDetails := []DateDetails{}
	for _, event := range events {
		eventTime := time.Unix(event.Timestamp, 0)
		year, month, day := eventTime.Date()
		weekYear, week := eventTime.ISOWeek()
		eventDayDetail := DateDetails{
			Year:     year,
			Month:    int(month),
			Day:      day,
			YearDay:  eventTime.YearDay(),
			WeekYear: weekYear,
			Week:     week,
		}
		daysDetails = append(daysDetails, eventDayDetail)
	}
	return daysDetails
}

func (handler SummaryGetByUserHandler) getSummary(events []db.Events) *SummaryByType {
	firstElement := events[0]
	lastElement := events[len(events)-1]
	gapDays := int(time.Now().Sub(time.Unix(lastElement.Timestamp, 0)).Hours() / 24.0)

	daysDetails := handler.getDateDetail(events)
	daysMap := make(map[int]int)
	weeksMap := make(map[int]int)
	for _, daysDetail := range daysDetails {
		if daysDetail.Year == time.Now().Year() {
			daysMap[daysDetail.YearDay] = daysMap[daysDetail.YearDay] + 1
		}
		currentWeekYear, _ := time.Now().ISOWeek()
		if daysDetail.WeekYear == currentWeekYear {
			weeksMap[daysDetail.Week] = weeksMap[daysDetail.Week] + 1
		}
	}
	continuesDays := 0
	currentDay := time.Now().YearDay() - 1
	if daysMap[currentDay+1] > 0 {
		continuesDays = 1
	}
	for currentDay > 0 && daysMap[currentDay] > 0 {
		continuesDays = continuesDays + 1
		currentDay = currentDay - 1
	}
	continuesWeeks := 0
	_, currentWeek := time.Now().ISOWeek()
	if weeksMap[currentWeek] >= 2 {
		continuesWeeks = 1
	}
	currentWeek = currentWeek - 1
	for currentWeek > 0 && weeksMap[currentWeek] >= 2 {
		continuesWeeks = continuesWeeks + 1
		currentWeek = currentWeek - 1
	}

	gapDays = 0
	if continuesDays == 0 {
		currentDay := time.Now().YearDay() - 1
		for currentDay > 0 && daysMap[currentDay] == 0 {
			gapDays = gapDays + 1
			currentDay = currentDay - 1
		}
	}

	status := "DISCONTINUED"
	if continuesWeeks > 0 {
		status = "WEEK_CONTINUES"
	}
	if continuesDays > 0 {
		status = "DAY_CONTINUES"
	}

	return &SummaryByType{
		StartingTime:   firstElement.Timestamp,
		MostRecentTime: lastElement.Timestamp,
		TotalDays:      len(daysMap),
		TotalWeeks:     len(weeksMap),
		ContinuesDays:  continuesDays,
		ContinuesWeeks: continuesWeeks,
		GapDays:        gapDays,
		Status:         status,
		Description:    getDescription(firstElement.Type, status, gapDays, continuesDays, continuesWeeks, firstElement.Timestamp),
		DaysDetails:    daysDetails,
		EventsDetails:  events,
	}
}

func getDescription(eventType string, status string, gapDays int, continuesDays int, continuesWeeks int, startingTime int64) string {
	eventTime := time.Unix(startingTime, 0)
	date := eventTime.Format("2006-01-02")
	switch status {
	case "DISCONTINUED":
		return fmt.Sprintf("距离上次 %s 已经过去了 %d 天", eventType, gapDays)
	case "WEEK_CONTINUES":
		return fmt.Sprintf("%s 第一次开始 %s, 最近已经坚持连续 %d 周", date, eventType, continuesWeeks)
	case "DAY_CONTINUES":
		return fmt.Sprintf("%s 第一次开始 %s, 最近已经坚持连续 %d 天", date, eventType, continuesDays)
	}
	return "Error"
}

func (handler SummaryGetByUserHandler) getEventsByDate(events []db.Events) map[string][]string {
	eventsByDate := make(map[string][]string)
	for _, event := range events {
		eventTime := time.Unix(event.Timestamp, 0)
		date := eventTime.Format("2006-01-02")
		if eventsByDate[date] == nil {
			eventsByDate[date] = []string{event.Type}
		} else {
			exists := false
			for _, eventType := range eventsByDate[date] {
				if eventType == event.Type {
					exists = true
					break
				}
			}
			if !exists {
				eventsByDate[date] = append(eventsByDate[date], event.Type)
			}
		}
	}
	return eventsByDate
}

func (handler SummaryGetByUserHandler) Handle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var events []db.Events
	handler.Database.Database.Where("user = ?", ps.ByName("user")).Order("Timestamp").Find(&events)
	eventsByType := make(map[string][]db.Events)
	for _, event := range events {
		if eventsByType[event.Type] == nil {
			eventsByType[event.Type] = []db.Events{event}
		} else {
			eventsByType[event.Type] = append(eventsByType[event.Type], event)
		}
	}
	response := &SummaryResponse{
		EventsByDate: handler.getEventsByDate(events),
		EventsByType: make(map[string]*SummaryByType),
	}
	for k, v := range eventsByType {
		response.EventsByType[k] = handler.getSummary(v)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	fmt.Printf("Getting events summary for user %s\n", ps.ByName("user"))
}

