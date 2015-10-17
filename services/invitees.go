package services

import (
	"github.com/grounded042/capacious/entities"
	"github.com/grounded042/capacious/utils"
)

type InviteeGateway interface {
	// GetAllInvitees gets all of the invitees in
	// the db for a specified event.
	// TODO: this will need to pagination at some point
	GetAllInviteesForEvent(string) ([]entities.Invitee, error)
	// CreateInvitee creates an invitee from a supplied
	// invitee object
	CreateInvitee(*entities.Invitee) error
	// GetInviteeFromId fetches an invitee from the database
	// based on the supplied id
	GetInviteeFromId(string) (entities.Invitee, error)
	// UpdateInvitee updates an invitee in the database
	// with info from the passed in object
	UpdateInvitee(entities.Invitee) error
}

type InviteeService struct {
	da InviteeGateway
}

func NewInviteeService(newDa InviteeGateway) InviteeService {
	return InviteeService{
		da: newDa,
	}
}

func (is InviteeService) GetInviteesForEvent(eventId string) ([]entities.Invitee, utils.Error) {
	invitees, err := is.da.GetAllInviteesForEvent(eventId)

	if err != nil {
		return []entities.Invitee{}, utils.NewApiError(500, err.Error())
	}

	return invitees, nil
}

func (is InviteeService) CreateInviteeForEvent(invitee *entities.Invitee, event entities.Event) utils.Error {
	invitee.FkEventId = event.EventId

	err := is.da.CreateInvitee(invitee)

	if err != nil {
		return utils.NewApiError(500, err.Error())
	}

	return nil
}

func (is InviteeService) GetInviteeFromId(id string) (entities.Invitee, utils.Error) {
	invitee, err := is.da.GetInviteeFromId(id)

	if err != nil {
		return entities.Invitee{}, utils.NewApiError(500, err.Error())
	}

	return invitee, nil
}

func (is InviteeService) EditInvitee(updateMe entities.Invitee) utils.Error {
	// TODO: make sure to constrain the number of guests here
	err := is.da.UpdateInvitee(updateMe)

	if err != nil {
		return utils.NewApiError(500, err.Error())
	}

	return nil
}
