package models

type UserRole struct {
	RoleId int `gorm:"primaryKey;not null" json:"role_id"`
	UserId int `gorm:"primaryKey;not null" json:"user_id"`
}
