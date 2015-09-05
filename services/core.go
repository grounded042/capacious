package services

import "github.com/grounded042/capacious/dal"

type List struct {
	Events EventsService
}

func NewServicesList(newDa dal.DataHandler) List {
	return List{
		Events: NewEventsService(newDa),
	}
}
