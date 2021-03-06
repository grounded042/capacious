package services

import (
	"github.com/grounded042/capacious/entities"
	"github.com/grounded042/capacious/utils"
)

type inviteeGateway interface {
	// GetAllInvitees gets all of the invitees in
	// the db for a specified event.
	// TODO: this will need to pagination at some point
	GetAllInviteesForEvent(string, int, int) ([]entities.Invitee, error)
	// CreateInvitee creates an invitee from a supplied
	// invitee object
	CreateInvitee(*entities.Invitee) error
	// GetInviteeFromId fetches an invitee from the database
	// based on the supplied id
	GetInviteeFromID(string) (entities.Invitee, error)
	// GetInviteeFriendFromId fetches an invitee friend from the
	// database based on the supplied id
	GetInviteeFriendFromID(string) (entities.InviteeFriend, error)
	// UpdateInvitee updates an invitee in the database
	// with info from the passed in object
	UpdateInvitee(entities.Invitee) error
	// CreateInviteeFriend create and invitee guest from
	// a supplied invitee guest object
	CreateInviteeFriend(*entities.InviteeFriend) error
	// UpdateInviteeFriend updates an invitee guest in the
	// database with info from the passed in object
	UpdateInviteeFriend(entities.InviteeFriend) error
	// SetGuestMenuChoices sets the menu choices with the
	// supplied choices for the supplied guest id
	SetGuestMenuChoices(string, []entities.MenuChoice) ([]entities.MenuChoice, error)
	// SetGuestMenuNote sets the menu note for a guest
	// based on the supplied guest id
	SetGuestMenuNote(string, entities.MenuNote) (entities.MenuNote, error)
	// SetInviteeSeatingRequests sets the invitee seating requests for an invitee
	// based on the suppllied invitee id
	SetInviteeSeatingRequests(string, []entities.InviteeSeatingRequest) ([]entities.InviteeSeatingRequest, error)
	// GetSeatingRequestInviteesForEvent gets a list of invitees that only includes the needed info for
	// seating requests
	GetSeatingRequestInviteesForEvent(string) ([]entities.Invitee, error)
	// GetNumberOfInviteesForEvent gets the number of invitees for the supplied
	// event id
	GetNumberOfInviteesForEvent(string) int
}

// the invitee is a subset of the event object -
type inviteeService struct {
	da inviteeGateway
}

func newInviteeService(newDa inviteeGateway) inviteeService {
	return inviteeService{
		da: newDa,
	}
}

func (is inviteeService) GetInviteesForEvent(eventId string, p *PaginationService) ([]entities.Invitee, utils.Error) {
	invitees, err := is.da.GetAllInviteesForEvent(eventId, p.GetStartNumber(), p.GetSize())

	p.SetNumItems(is.da.GetNumberOfInviteesForEvent(eventId))

	if err != nil {
		return []entities.Invitee{}, utils.NewApiError(500, err.Error())
	}

	return invitees, nil
}

func (is inviteeService) CreateInviteeForEvent(invitee *entities.Invitee, event entities.Event) utils.Error {
	invitee.FkEventID = event.EventID

	err := is.da.CreateInvitee(invitee)

	if err != nil {
		return utils.NewApiError(500, err.Error())
	}

	return nil
}

func (is inviteeService) GetInviteeFromID(id string) (entities.Invitee, utils.Error) {
	invitee, err := is.da.GetInviteeFromID(id)

	if err != nil {
		return entities.Invitee{}, utils.NewApiError(500, err.Error())
	}

	return invitee, nil
}

func (is inviteeService) EditInvitee(updateMe entities.Invitee) utils.Error {
	err := is.da.UpdateInvitee(updateMe)

	if err != nil {
		return utils.NewApiError(500, err.Error())
	}

	return nil
}

func (is inviteeService) CreateInviteeFriend(friend *entities.InviteeFriend) utils.Error {
	err := is.da.CreateInviteeFriend(friend)

	if err != nil {
		return utils.NewApiError(500, err.Error())
	}

	return nil
}

func (is inviteeService) EditInviteeFriend(updateMe entities.InviteeFriend) utils.Error {
	err := is.da.UpdateInviteeFriend(updateMe)

	if err != nil {
		return utils.NewApiError(500, err.Error())
	}

	return nil
}

func (is inviteeService) SetGuestMenuChoices(guestID string, choices []entities.MenuChoice) ([]entities.MenuChoice, utils.Error) {
	// make sure that the FkGuestId is set correctly
	for key, _ := range choices {
		choices[key].FkGuestID = guestID
	}

	updatedChoices, err := is.da.SetGuestMenuChoices(guestID, choices)

	if err != nil {
		return []entities.MenuChoice{}, utils.NewApiError(500, err.Error())
	}

	return updatedChoices, nil
}

func (is inviteeService) GetInviteeFriendFromID(id string) (entities.InviteeFriend, utils.Error) {
	iFriend, err := is.da.GetInviteeFriendFromID(id)

	if err != nil {
		return entities.InviteeFriend{}, utils.NewApiError(500, err.Error())
	}

	return iFriend, nil
}

func (is inviteeService) SetGuestMenuNote(guestID string, note entities.MenuNote) (entities.MenuNote, utils.Error) {
	// make sure that the FkGuestId is set correctly
	note.FkGuestID = guestID

	updatedNote, err := is.da.SetGuestMenuNote(guestID, note)

	if err != nil {
		return entities.MenuNote{}, utils.NewApiError(500, err.Error())
	}

	return updatedNote, nil
}

func (is inviteeService) SetInviteeSeatingRequests(inviteeID string, requests []entities.InviteeSeatingRequest) ([]entities.InviteeSeatingRequest, utils.Error) {
	// make sure that the FkGuestId is set correctly
	for i := len(requests) - 1; i >= 0; i-- {
		requests[i].FkInviteeID = inviteeID

		// make sure someone is not setting themselves to be a friend
		if requests[i].FkInviteeRequestID == inviteeID {
			requests = append(requests[:i], requests[i+1:]...)
		} else {
			// clear any existing primary keys - remember, we are replacing, not
			// updating
			requests[i].InviteeSeatingRequestID = ""
		}
	}

	toReturn, err := is.da.SetInviteeSeatingRequests(inviteeID, requests)

	if err != nil {
		return []entities.InviteeSeatingRequest{}, utils.NewApiError(500, err.Error())
	}

	return toReturn, nil
}

func (is inviteeService) GetSeatingRequestInviteesForEvent(eventID string) ([]entities.Invitee, utils.Error) {
	toReturn, err := is.da.GetSeatingRequestInviteesForEvent(eventID)

	if err != nil {
		return []entities.Invitee{}, utils.NewApiError(500, err.Error())
	}

	return toReturn, nil
}
