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
