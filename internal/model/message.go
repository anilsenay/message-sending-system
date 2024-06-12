package model

import (
	"errors"
	"regexp"
	"time"
)

type MessageStatus string

const (
	MESSAGE_UNSENT     MessageStatus = "unsent"
	MESSAGE_PROCESSING MessageStatus = "processing"
	MESSAGE_SENT       MessageStatus = "sent"
	MESSAGE_FAILED     MessageStatus = "failed"
)

type Message struct {
	Id                   int           `gorm:"column:id" json:"id"`
	Content              string        `gorm:"column:content" json:"content"`
	RecipientPhoneNumber string        `gorm:"column:recipient_phone_number" json:"recipient_phone_number"`
	Status               MessageStatus `gorm:"column:status" json:"status"`
	CreatedAt            time.Time     `gorm:"column:created_at" json:"created_at"`
	SentAt               *time.Time    `gorm:"column:sent_at" json:"sent_at"`
}

func (Message) TableName() string {
	return "message"
}

type MessageCreateRequest struct {
	Content              string `gorm:"column:content" json:"content"`
	RecipientPhoneNumber string `gorm:"column:recipient_phone_number" json:"recipient_phone_number"`
}

var phoneNumberRegex = regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)

func (m MessageCreateRequest) Validate() error {
	if len(m.Content) > 1000 {
		return errors.New("the length of 'content' cannot exceed 1000 characters")
	}
	if !phoneNumberRegex.MatchString(m.RecipientPhoneNumber) {
		return errors.New("invalid phone number")
	}
	return nil
}
