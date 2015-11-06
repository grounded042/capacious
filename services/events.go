package services

import (
	"github.com/grounded042/capacious/entities"
	"github.com/grounded042/capacious/utils"
)

type eventsGateway interface {
	// GetAllEvents gets all of the events in the db
	// TODO: This will need to change eventually since
	// we will want to lockdown which events you can
	// get bases on your user id
	GetAllEvents() ([]entities.Event, error)
	// CreateEvent creates an event from a supplied
	// event object
	CreateEvent(*entities.Event) error
	// GetEventInfo gets the info for an event matching
	// the supplied event id
	GetEventInfo(eventId string) (entities.Event, error)
}

type eventsService struct {
	da eventsGateway
}

func newEventsService(newDa eventsGateway) eventsService {
	return eventsService{
		da: newDa,
	}
}

func (es eventsService) GetEvents() ([]entities.Event, utils.Error) {
	events, err := es.da.GetAllEvents()

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

func (es eventsService) CreateEvent(event *entities.Event) utils.Error {
	err := es.da.CreateEvent(event)

	if err != nil {
		return utils.NewApiError(500, err.Error())
	}

	return nil
}
