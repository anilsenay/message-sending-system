package model

type MessageRedisPayload struct {
	MessageId string `json:"messageId"`
	Timestamp int    `json:"timestamp"`
}
