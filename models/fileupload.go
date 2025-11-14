package models

import "time"

type FileUpload struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	FileName    string    `json:"filename"`
	ContentType string    `json:"content_type"`
	Size        int       `json:"size"`
	Data        []byte    `json:"-"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
