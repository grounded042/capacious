package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/grounded042/capacious/entities"
	"github.com/grounded042/capacious/utils"
	"github.com/zenazn/goji/web"
)

type EventsStub interface {
	GetEvents() ([]entities.Event, utils.Error)
	GetEventInfo(eventId string) (entities.Event, utils.Error)
	CreateEvent(*entities.Event) utils.Error
	GetMenuItemsForEvent(eventID string) ([]entities.MenuItem, utils.Error)
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

	if events, err := ec.es.GetEvents(); err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
	} else {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(events)
	}
}

func (ec EventsController) GetEventInfo(c web.C, w http.ResponseWriter, r *http.Request) {

	if event, err := ec.es.GetEventInfo(c.URLParams["id"]); err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
	} else {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(event)
	}
}

func (ec EventsController) CreateEvent(c web.C, w http.ResponseWriter, r *http.Request) {
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

	if err := ec.es.CreateEvent(&event); err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
	} else {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(event)
	}
}

// GetMenuItemsForEvent renders the menu items for an event using w.
func (ec EventsController) GetMenuItemsForEvent(c web.C, w http.ResponseWriter, r *http.Request) {
	if items, err := ec.es.GetMenuItemsForEvent(c.URLParams["id"]); err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
	} else {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(items)
	}
}
