package models

import "time"

type Message struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	SenderID   uint      `json:"sender_id" gorm:"not null"`
	Sender     User      `json:"sender" gorm:"foreignKey:SenderID"`
	ReceiverID uint      `json:"receiver_id" gorm:"not null"`
	Receiver   User      `json:"receiver" gorm:"foreignKey:ReceiverID"`
	Content    string    `json:"content" gorm:"not null"`
	Read       bool      `json:"read" gorm:"default:false"`
	CreatedAt  time.Time `json:"created_at"`
}
