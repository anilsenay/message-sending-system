package model

import "time"

type MessageStatus string

const (
	MESSAGE_UNSENT     MessageStatus = "unsent"
	MESSAGE_PROCESSING MessageStatus = "processing"
	MESSAGE_SENT       MessageStatus = "sent"
)

type Message struct {
	Id                   int           `gorm:"column:id" json:"id"`
	Content              string        `gorm:"column:content" json:"content"`
	RecipientPhoneNumber string        `gorm:"column:recipient_phone_number" json:"recipient_phone_number"`
	Status               MessageStatus `gorm:"column:status" json:"status"`
	CreatedAt            time.Time     `gorm:"column:created_at" json:"created_at"`
}

func (Message) TableName() string {
	return "message"
}
