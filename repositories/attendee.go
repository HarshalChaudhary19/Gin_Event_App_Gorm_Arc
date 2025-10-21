package repositories

import (
	"Gin_Event_App_Arc/models"
	"context"
	"time"

	"gorm.io/gorm"
)

type AttendeeRepo struct {
	DB *gorm.DB
}

func InitAttendeeRepo(db *gorm.DB) *AttendeeRepo {
	return &AttendeeRepo{
		DB: db,
	}
}

type AttendeeRepoI interface {
	GetByEventAndAttendee(int, int) (*models.Attendees, error)
	Insert(*models.Attendees) error
	GetAttendeesByEvent(int) (*[]models.User, error)
	Delete(int, int) error
	GetAllAttendeesandEventList() ([]*models.Attendees, error)
}

func (repo *AttendeeRepo) GetByEventAndAttendee(eventId, userId int) (*models.Attendees, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var attendee models.Attendees
	err := repo.DB.WithContext(ctx).Where("event_id=? AND user_id=?", eventId, userId).First(&attendee).Error
	return &attendee, err
}

func (repo *AttendeeRepo) Insert(attendee *models.Attendees) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return repo.DB.WithContext(ctx).Create(&attendee).Error
}

func (repo *AttendeeRepo) GetAttendeesByEvent(eventId int) (*[]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var userListByEventID *[]models.User
	err := repo.DB.WithContext(ctx).Table("users u").Select("u.id,u.name,u.email").Joins("JOIN attendees a ON u.id=a.user_id").Where("a.event_id=?", eventId).Scan(&userListByEventID).Error
	return userListByEventID, err
}

func (repo *AttendeeRepo) Delete(userId, eventId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return repo.DB.WithContext(ctx).Where("user_id=? AND event_id=?", userId, eventId).Delete(&models.Attendees{}).Error
}

func (repo *AttendeeRepo) GetAllAttendeesandEventList() ([]*models.Attendees, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var attendeeList []*models.Attendees
	err := repo.DB.WithContext(ctx).Find(&attendeeList).Error
	return attendeeList, err
}
