package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/grounded042/capacious/entities"
	"github.com/grounded042/capacious/services"
	"github.com/grounded042/capacious/utils"
	"github.com/zenazn/goji/web"
)

type EventsStub interface {
	GetEvents(string) ([]entities.Event, utils.Error)
	GetEventInfo(eventId string) (entities.Event, utils.Error)
	GetEventStats(eventID string, userID string) (services.EventStats, utils.Error)
	CreateEvent(*entities.Event, string) utils.Error
	GetMenuItemsForEvent(eventID string) ([]entities.MenuItem, utils.Error)
	GetListOfSeatingRequestChoices(eventID string) ([]entities.SeatingRequestChoice, utils.Error)
}

type EventsController struct {
	es EventsStub
}

func NewEventsController(newEs EventsStub) EventsController {
	return EventsController{
		es: newEs,
	}
}

func (ec EventsController) GetEvents(c web.C, w http.ResponseWriter, r *http.Request) {
	userID, ok := checkForAndHandleUserIDInContext(c, w, "You need a valid user id to get your list of events!")

	if !ok {
		return
	}

	if events, err := ec.es.GetEvents(userID); err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
	} else {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(events)
	}
}

// GetEventInfo gets the info for a specific event. This function does not
// require auth as it is used to get event info for responses to invitations
func (ec EventsController) GetEventInfo(c web.C, w http.ResponseWriter, r *http.Request) {

	if event, err := ec.es.GetEventInfo(c.URLParams["id"]); err != nil {
		w.WriteHeader(utils.GetCodeForError(err))
		fmt.Println(err)
	} else {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(event)
	}
}

// GetEventStats gets the stats for the specified event.
func (ec EventsController) GetEventStats(c web.C, w http.ResponseWriter, r *http.Request) {
	userID, ok := checkForAndHandleUserIDInContext(c, w, "You need a valid user id to get your list of events!")

	if !ok {
		return
	}

	if stats, err := ec.es.GetEventStats(c.URLParams["id"], userID); err != nil {
		w.WriteHeader(utils.GetCodeForError(err))
		fmt.Println(err)
	} else {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(stats)
	}
}

func (ec EventsController) CreateEvent(c web.C, w http.ResponseWriter, r *http.Request) {
	userID, ok := checkForAndHandleUserIDInContext(c, w, "You need a valid user id to create an event!")

	if !ok {
		return
	}

	// decode the body into an event
	var event entities.Event

	rBody, ioErr := ioutil.ReadAll(r.Body)

	if ioErr != nil {
		w.WriteHeader(500)
		fmt.Println(ioErr)
		return
	}

	if err := json.Unmarshal(rBody, &event); err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
		return
	}

	// create the event
	if err := ec.es.CreateEvent(&event, userID); err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
	} else {
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(event)
	}
}

// GetMenuItemsForEvent renders the menu items for an event using w.
func (ec EventsController) GetMenuItemsForEvent(c web.C, w http.ResponseWriter, r *http.Request) {
	if items, err := ec.es.GetMenuItemsForEvent(c.URLParams["id"]); err != nil {
		w.WriteHeader(utils.GetCodeForError(err))
		fmt.Println(err)
	} else {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(items)
	}
}

func (ec EventsController) GetListOfSeatingRequestChoices(c web.C, w http.ResponseWriter, r *http.Request) {
	if choices, err := ec.es.GetListOfSeatingRequestChoices(c.URLParams["id"]); err != nil {
		w.WriteHeader(utils.GetCodeForError(err))
		fmt.Println(err)
	} else {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(choices)
	}
}
