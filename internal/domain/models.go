package domain

import "time"

type Chat struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Title     string    `json:"title" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	Message   []Message `json:"message,omitempty" gorm:"foreignKey:ChatID;constraint:OnDelete:CASCADE"`
}

type Message struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	ChatID    uint      `json:"chat_id" gorm:"not null"`
	Text      string    `json:"text" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
}
