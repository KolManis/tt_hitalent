package models

import "time"

type Message struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ChatID    uint      `gorm:"not null;index" json:"chat_id"`
	Text      string    `gorm:"size:5000;not null" json:"text"`
	CreatedAt time.Time `json:"created_at"`
}
