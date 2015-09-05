package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/grounded042/capacious/services"
	"github.com/zenazn/goji/web"
)

type EventsController struct {
	evSvc services.EventsService
}

func NewEventsController(newEvSvc services.EventsService) EventsController {
	return EventsController{
		evSvc: newEvSvc,
	}
}

func (ec EventsController) GetEvents(c web.C, w http.ResponseWriter, r *http.Request) {

	if events, err := ec.evSvc.GetEvents(); err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
	} else {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(events)
	}
}
