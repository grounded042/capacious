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
}

type eventsService struct {
	da eventsGateway
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
