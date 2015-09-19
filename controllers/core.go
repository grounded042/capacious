package controllers

import "github.com/grounded042/capacious/services"

type List struct {
	Events   EventsController
	Invitees InviteesController
}

func NewControllersList(newSvcsList services.List) List {
	return List{
		Events:   NewEventsController(newSvcsList),
		Invitees: NewInviteesController(newSvcsList),
	}
}
