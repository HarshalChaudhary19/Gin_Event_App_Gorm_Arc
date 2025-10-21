package services

import (
	"Gin_Event_App_Arc/models"
	"Gin_Event_App_Arc/repositories"
)

type EventServiceI interface {
	Insert(*models.Event) error
	Get(int) (*models.Event, error)
	GetAll(int) ([]*models.Event, int, error)
	Update(*models.Event, int) (*models.Event, error)
	Delete(int) error
	GetEventByAttendee(int) ([]*models.Event, error)
}

type EventService struct {
	EventServ repositories.EventRepoI
}

func NewEventService(eventRepo repositories.EventRepoI) EventServiceI {
	return &EventService{EventServ: eventRepo}
}

//Services are down below

func (service *EventService) Insert(event *models.Event) error {
	err := service.EventServ.Insert(event)
	if err != nil {
		return err
	}
	return nil
}

// Get One Event
func (service *EventService) Get(id int) (*models.Event, error) {
	event, err := service.EventServ.Get(id)
	if err != nil {
		return nil, err
	}
	return event, nil
}

// Get All Events
func (service *EventService) GetAll(cursor int) ([]*models.Event, int, error) {
	pageLimit := 5
	eventsList, err := service.EventServ.GetAll(cursor, pageLimit)
	if err != nil {
		return nil, cursor, err
	}
	nextCursor := 0
	if len(eventsList) > 0 {
		nextCursor = eventsList[len(eventsList)-1].Id
	}
	return eventsList, nextCursor, nil
}

//Update Event

func (service *EventService) Update(event *models.Event, id int) (*models.Event, error) {

	eventNew := &models.Event{
		Name:        event.Name,
		Description: event.Description,
		Date:        event.Date,
		Location:    event.Location,
	}
	eventLatest, err := service.EventServ.Update(eventNew, id)
	if err != nil {
		return nil, err
	}
	return eventLatest, nil

}

func (service *EventService) Delete(id int) error {
	tx := service.EventServ.Delete(id)
	if tx.Error != nil {
		return tx
	}
	// if tx.RowsAffected == 0 {
	// 	return fmt.Errorf("No event found with id %d", id)
	// }
	return nil
}

func (service *EventService) GetEventByAttendee(attendeeId int) ([]*models.Event, error) {
	eventsList, err := service.EventServ.GetEventByAttendee(attendeeId)
	if err != nil {
		return nil, err
	}
	return eventsList, nil
}
