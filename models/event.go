package models

// If we want cascading delete in postgres so we need to explicitly use it
type Event struct {
	// Id          int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Id          int    `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"not null" json:"name" binding:"required,min=3"`
	OwnerId     int    `gorm:"not null" json:"owner_id" binding:"required" `
	Description string `gorm:"not null" json:"description" binding:"required,min=10"`
	Date        string `gorm:"not null" json:"date" binding:"required"`
	Location    string `gorm:"not null" json:"location" binding:"required,min=3"`
}
