package controllers

import "github.com/grounded042/capacious/services"

type List struct {
	Events   EventsController
	Invitees InviteesController
	Auth     AuthController
}

func NewControllersList(coord services.Coordinator) List {
	return List{
		Events:   NewEventsController(coord),
		Invitees: NewInviteesController(coord),
		Auth:     NewAuthController(coord),
	}
}
