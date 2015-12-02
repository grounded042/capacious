package services

import (
	"github.com/grounded042/capacious/entities"
	"github.com/grounded042/capacious/utils"
)

type eventsGateway interface {
	// GetAllEvents gets all of the events in the db that the passed in userID is
	// an admin of.
	GetAllEvents(userID string) ([]entities.Event, error)
	// CreateEvent creates an event from a supplied event object and adds the
	// specified user id as an owner of the event
	CreateEvent(*entities.Event, string) error
	// GetEventInfo gets the info for an event matching
	// the supplied event id
	GetEventInfo(eventId string) (entities.Event, error)
	// GetMenuItemsForEvent gets all of the menu items
	// for an event matching the supplied event id
	GetMenuItemsForEvent(eventID string) ([]entities.MenuItem, error)
	// GetEventAdminRecordForUserAndEventID gets the event admin record that
	// contains both the user id UserID and the event id EventID.
	GetEventAdminRecordForUserAndEventID(userID string, eventID string) (entities.EventAdmin, error)
	// GetNumAttendingForEvent gets the number of guests that are attending for
	// the specified event id
	GetNumAttendingForEvent(eventID string) (int, error)
}

type eventsService struct {
	da eventsGateway
}

type EventStats struct {
	NumInvitees  int `json:"num_invitees"`
	NumAttending int `json:"num_attending"`
}

func newEventsService(newDa eventsGateway) eventsService {
	return eventsService{
		da: newDa,
	}
}

func (es eventsService) GetEvents(userID string) ([]entities.Event, utils.Error) {
	events, err := es.da.GetAllEvents(userID)

	if err != nil {
		return []entities.Event{}, utils.NewApiError(500, err.Error())
	}

	return events, nil
}

func (es eventsService) GetEventInfo(eventId string) (entities.Event, utils.Error) {
	event, err := es.da.GetEventInfo(eventId)

	if err != nil {
		return entities.Event{}, utils.NewApiError(500, err.Error())
	}

	return event, nil
}

func (es eventsService) CreateEvent(event *entities.Event, userID string) utils.Error {
	err := es.da.CreateEvent(event, userID)

	if err != nil {
		return utils.NewApiError(500, err.Error())
	}

	return nil
}

// GetMenuItemsForEvent gets the menu items based on the event id
// eventID. It returns a slice of menu items and any errors that occured.
func (es eventsService) GetMenuItemsForEvent(eventID string) ([]entities.MenuItem, utils.Error) {
	items, err := es.da.GetMenuItemsForEvent(eventID)

	if err != nil {
		return []entities.MenuItem{}, utils.NewApiError(500, err.Error())
	}

	return items, nil
}

func (es eventsService) IsUserAnAdminForEvent(userID string, eventID string) (bool, utils.Error) {
	eAdmin, err := es.da.GetEventAdminRecordForUserAndEventID(userID, eventID)

	if err != nil && err.Error() != "record not found" {
		return false, utils.NewApiError(500, err.Error())
	}

	return eAdmin.EventAdminID != "", nil
}

func (es eventsService) GetNumAttendingForEvent(eventID string) (int, utils.Error) {
	num, err := es.da.GetNumAttendingForEvent(eventID)

	if err != nil {
		return 0, utils.NewApiError(500, err.Error())
	}

	return num, nil
}
