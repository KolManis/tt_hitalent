package models

import "time"

type Chat struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"size:200;not null" json:"title"`
	CreatedAt time.Time `json:"created_at"`
}
