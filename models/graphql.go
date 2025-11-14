package models

type CategoryGraph struct {
	ID          int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string `gorm:"not null" json:"name"`
	Description string `gorm:"not null" json:"description"`
}

type CourseGraph struct {
	ID          int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string `gorm:"not null" json:"name"`
	Description string `gorm:"not null" json:"description"`
}
