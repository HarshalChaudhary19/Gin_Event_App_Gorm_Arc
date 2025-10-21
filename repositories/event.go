package repositories

import (
	"Gin_Event_App_Arc/models"
	"context"
	"time"

	"gorm.io/gorm"
)

type EventRepoI interface {
	Insert(*models.Event) error
	Get(int) (*models.Event, error)
	GetAll(int, int) ([]*models.Event, error)
	Update(*models.Event, int) (*models.Event, error)
	Delete(int) error
	GetEventByAttendee(int) ([]*models.Event, error)
}

type EventRepo struct {
	DB *gorm.DB
}

func InitEventRepo(db *gorm.DB) *EventRepo {
	return &EventRepo{
		DB: db,
	}
}

func (repo *EventRepo) Insert(event *models.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return repo.DB.WithContext(ctx).Create(event).Error
}

func (repo *EventRepo) Get(id int) (*models.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var event *models.Event
	err := repo.DB.WithContext(ctx).Where("id=?", id).First(&event).Error
	return event, err
}

func (repo *EventRepo) GetAll(cursor, pageLimit int) ([]*models.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var events []*models.Event
	err := repo.DB.WithContext(ctx).Where("id>?", cursor).Limit(pageLimit).Find(&events).Error
	return events, err
}

func (repo *EventRepo) Update(event *models.Event, id int) (*models.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := repo.DB.WithContext(ctx).Where("id=?", event.Id).Model(&models.Event{}).Updates(models.Event{ //Basic Update Query hai ye
		Name:        event.Name,
		Description: event.Description,
		Date:        event.Date,
		Location:    event.Location,
	}).Error
	return event, err
}

func (repo *EventRepo) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return repo.DB.WithContext(ctx).Delete(&models.Event{}, id).Error
}

func (repo *EventRepo) GetEventByAttendee(attendeeID int) ([]*models.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var eventsList []*models.Event
	err := repo.DB.WithContext(ctx).Table("events e").Select("e.id,e.name,e.owner_id,e.description,e.date,e.location").
		Joins("JOIN attendees a ON a.event_id=e.id").Where("a.user_id = ?", attendeeID).Scan(&eventsList).Error
	return eventsList, err
}
