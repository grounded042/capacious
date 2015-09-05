package controllers

import "github.com/grounded042/capacious/services"

type List struct {
	Events EventsController
}

func NewControllersList(newSvcsList services.List) List {
	return List{
		Events: NewEventsController(newSvcsList.Events),
	}
}
