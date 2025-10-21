package services

import (
	"Gin_Event_App_Arc/models"
	"Gin_Event_App_Arc/repositories"
	"fmt"

	"gorm.io/gorm"
)

type AttendeeServiceI interface {
	GetByEventAndAttendee(int, int) (*models.Attendees, error)
	Insert(*models.Attendees) (*models.Attendees, error)
	Delete(int, int) error
	GetAttendeesByEvent(int) (*[]models.User, error)
	GetAllAttendeesandEventList() ([]*models.Attendees, error)
}

type AttendeeService struct {
	AttendeeServe repositories.AttendeeRepoI
}

func NewAttendeeService(attendeeRepo repositories.AttendeeRepoI) AttendeeServiceI {
	return &AttendeeService{AttendeeServe: attendeeRepo}
}

func (service *AttendeeService) GetByEventAndAttendee(eventID int, userId int) (*models.Attendees, error) {
	attendee, err := service.AttendeeServe.GetByEventAndAttendee(eventID, userId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return attendee, nil

}

func (service *AttendeeService) Insert(attendee *models.Attendees) (*models.Attendees, error) {
	err := service.AttendeeServe.Insert(attendee)
	if err != nil {
		fmt.Println("Error yahan se hai")
		return nil, err
	}
	return nil, nil
}

func (service *AttendeeService) GetAttendeesByEvent(eventID int) (*[]models.User, error) {
	usersList, err := service.AttendeeServe.GetAttendeesByEvent(eventID)
	if err != nil {
		return nil, err
	}
	return usersList, err
}

func (service *AttendeeService) Delete(eventId, userId int) error {

	err := service.AttendeeServe.Delete(userId, eventId)

	if err != nil {
		return err
	}
	return nil
}

func (service *AttendeeService) GetAllAttendeesandEventList() ([]*models.Attendees, error) {
	attendeeList, err := service.AttendeeServe.GetAllAttendeesandEventList()
	if err != nil {
		return nil, err
	}
	return attendeeList, nil
}
