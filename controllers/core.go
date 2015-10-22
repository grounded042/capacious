package controllers

import "github.com/grounded042/capacious/services"

type List struct {
	Events   EventsController
	Invitees InviteesController
}

func NewControllersList(coord services.Coordinator) List {
	return List{
		Events:   NewEventsController(coord),
		Invitees: NewInviteesController(coord),
	}
}
