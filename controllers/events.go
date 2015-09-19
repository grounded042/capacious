package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/grounded042/capacious/entities"
	"github.com/grounded042/capacious/services"
	"github.com/zenazn/goji/web"
)

type EventsController struct {
	sl services.List
}

func NewEventsController(newSL services.List) EventsController {
	return EventsController{
		sl: newSL,
	}
}

func (ec EventsController) GetEvents(c web.C, w http.ResponseWriter, r *http.Request) {

	if events, err := ec.sl.Events.GetEvents(); err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
	} else {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(events)
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

	if err := ec.sl.Events.CreateEvent(&event); err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
	} else {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(event)
	}
}
