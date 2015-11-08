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

func (c Coordinator) GetMenuItemsForEvent(eventID string) ([]entities.MenuItem, utils.Error) {
	return c.events.GetMenuItemsForEvent(eventID)
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

func (c Coordinator) SetInviteeMenuChoices(inviteeID string, choices []entities.MenuChoice) ([]entities.MenuChoice, utils.Error) {
	invitee, err := c.invitees.GetInviteeFromId(inviteeID)

	if err != nil {
		return []entities.MenuChoice{}, err
	}

	return c.SetGuestMenuChoices(invitee.FkEventId, invitee.Self.GuestId, choices)
}

func (c Coordinator) SetInviteeFriendMenuChoices(iFriendID string, choices []entities.MenuChoice) ([]entities.MenuChoice, utils.Error) {
	iFriend, err := c.invitees.GetInviteeFriendFromId(iFriendID)

	if err != nil {
		return []entities.MenuChoice{}, err
	}

	invitee, err := c.invitees.GetInviteeFromId(iFriend.FkInviteeId)

	if err != nil {
		return []entities.MenuChoice{}, err
	}

	return c.SetGuestMenuChoices(invitee.FkEventId, iFriend.FkGuestId, choices)
}

func (c Coordinator) SetGuestMenuChoices(eventID string, guestID string, choices []entities.MenuChoice) ([]entities.MenuChoice, utils.Error) {
	// TODO: validate the menu choices
	items, err := c.events.GetMenuItemsForEvent(eventID)

	if err != nil {
		return []entities.MenuChoice{}, err
	}

	if !c.validateMenuChoicesWithMenuItems(choices, items) {
		return []entities.MenuChoice{}, utils.NewApiError(400, "So yeah...you had an error in the list of menu choices you sent in. Rejected!")
	}

	return c.invitees.SetGuestMenuChoices(guestID, choices)
}

// validateMenuChoicesWithMenuItems validates that the supplied choices match
// up with the supplied menu items. It returns a bool regarding the validity.
// TODO: unit test this sucker
func (c Coordinator) validateMenuChoicesWithMenuItems(choices []entities.MenuChoice, items []entities.MenuItem) bool {
	itemUsage := make(map[string]int)

	for _, cValue := range choices {
		itemMatch := false
		// what item does the current choice align to?
		for _, iValue := range items {
			for _, oValue := range iValue.Options {
				if oValue.MenuItemOptionId == cValue.FkMenuItemOptionId {
					// we found an item option that matches the choice id,
					// so we have validated that the choice is valid
					itemMatch = true

					// up the count for the number of times a choice has been used for an item
					itemUsage[iValue.MenuItemId] = itemUsage[iValue.MenuItemId] + 1

					// is the number of choices higher than the number of choices allowed?
					if itemUsage[iValue.MenuItemId] > iValue.NumChoices {
						return false
					}

				}
			}
		}

		// check and see if the choice matched up with an option in an item
		if !itemMatch {
			return false
		}

	}

	// if we made it this far, we are good!
	return true
}

//end invitee coordination
