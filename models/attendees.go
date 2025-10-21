package models

type Attendees struct {
	//AutoIncrement will not be there in postgres0 0a0nd need to add cascade too
	// Id      int `gorm:"primaryKey;autoIncrement" json:"id"`.	Id      int `gorm:"primaryKey" json:"id"`
	UserId  int `gorm:"not null" json:"userId"`
	EventId int `gorm:"not null" json:"eventId"`
}
