package models

type User struct {
	// Id       int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Id       int    `gorm:"primaryKey" json:"id"`
	Email    string `gorm:"unique;not null" json:"email"`
	Name     string `gorm:"not null" json:"name"`
	Password string `json:"-"`
}
