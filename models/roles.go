package models

type Role struct {
	Id   int    `gorm:"primaryKey" json:"id"`
	Name string `gorm:"unique;not null" json:"name"`
}
