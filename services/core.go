package services

import (
	"github.com/grounded042/capacious/dal"
	"github.com/grounded042/capacious/entities"
	"github.com/grounded042/capacious/utils"
)

// the Coordinator coordinates interactions between different services.
// it applies any overarching business logic relates to more than one service.
// this allows each service to only care about it's objects and eliminates
// the possibility of circular dependencies in the services

type Coordinator struct {
	events   eventsService
	invitees inviteeService
}

func NewCoordinator(newDa dal.DataHandler) Coordinator {
	return Coordinator{
		events:   newEventsService(newDa),
		invitees: newInviteeService(newDa),
	}
}

// events coordination

func (c Coordinator) GetEvents() ([]entities.Event, utils.Error) {
	return c.events.GetEvents()
}

func (c Coordinator) GetEventInfo(eventId string) (entities.Event, utils.Error) {
	return c.events.GetEventInfo(eventId)
}

func (c Coordinator) CreateEvent(event *entities.Event) utils.Error {
	return c.events.CreateEvent(event)
}

// end events coordination

// invitee coordination

func (c Coordinator) GetInviteesForEvent(eventId string) ([]entities.Invitee, utils.Error) {
	return c.invitees.GetInviteesForEvent(eventId)
}

func (c Coordinator) CreateInviteeForEvent(invitee *entities.Invitee, event entities.Event) utils.Error {
	return c.invitees.CreateInviteeForEvent(invitee, event)
}

func (c Coordinator) GetInviteeFromId(id string) (entities.Invitee, utils.Error) {
	return c.invitees.GetInviteeFromId(id)
}

func (c Coordinator) EditInvitee(updateMe entities.Invitee) utils.Error {
	return c.invitees.EditInvitee(updateMe)
}

func (c Coordinator) CreateInviteeFriend(updateMe *entities.InviteeFriend) utils.Error {
	// TODO: make sure to constrain the number of friends here
	return c.invitees.CreateInviteeFriend(updateMe)
}

func (c Coordinator) EditInviteeFriend(updateMe entities.InviteeFriend) utils.Error {
	// TODO: make sure to constrain the number of friends here
	return c.invitees.EditInviteeFriend(updateMe)
}

//end invitee coordination
