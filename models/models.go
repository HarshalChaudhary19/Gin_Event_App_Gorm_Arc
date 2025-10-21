// package models
package models

// import (
// 	"Gin_Event_App_Arc/config"

// 	"gorm.io/gorm"
// )

// type Models struct {
// 	Users     UserModel
// 	Events    EventModel
// 	Attendees AttendeesModel
// }

// type Application struct {
// 	Port   int
// 	Models Models
// 	Cache  *config.RedisClient
// }

// func NewModels(db *gorm.DB) Models {
// 	return Models{
// 		Users:     UserModel{DB: db},
// 		Events:    EventModel{DB: db},
// 		Attendees: AttendeesModel{DB: db},
// 	}
// }

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required,min=3"`
}
