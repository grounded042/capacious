package services

import (
	"github.com/grounded042/capacious/entities"
	"github.com/grounded042/capacious/utils"
)

type EventsGateway interface {
	// GetAllEvents gets all of the events in the db
	// TODO: This will need to change eventually since
	// we will want to lockdown which events you can
	// get bases on your user id
	GetAllEvents() ([]entities.Event, error)
}

type EventsService struct {
	da EventsGateway
}

func NewEventsService(newDa EventsGateway) EventsService {
	return EventsService{
		da: newDa,
	}
}

func (es EventsService) GetEvents() ([]entities.Event, utils.Error) {
	events, err := es.da.GetAllEvents()

	if err != nil {
		return []entities.Event{}, utils.NewApiError(500, err.Error())
	}

	return events, nil
}
