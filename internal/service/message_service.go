package service

import (
	"context"
	"fmt"
	"time"

	"github.com/anilsenay/message-sending-system/internal/model"
)

type messageRepository interface {
	RetrieveAll(ctx context.Context, filters model.Message) ([]model.Message, error)
	RetrieveMessagesForProcess(ctx context.Context, limit int) ([]model.Message, error)
	Update(ctx context.Context, model *model.Message, updates map[string]interface{}) error
}

type messageSender interface {
	Start(context.Context, func() error)
	Stop()
}

type redisClient interface {
	SetJson(ctx context.Context, key string, val interface{}) error
}

type messageClient interface {
	Send(to, content string) (model.MessageClientResponse, error)
}

type MessageService struct {
	messageRepository messageRepository
	messageSender     messageSender
	redisClient       redisClient
	messageClient     messageClient
	messageLimit      int
}

func NewMessageService(
	messageRepository messageRepository,
	messageSender messageSender,
	redisClient redisClient,
	messageClient messageClient,
	messageLimit int,
) *MessageService {
	return &MessageService{
		messageRepository: messageRepository,
		messageSender:     messageSender,
		redisClient:       redisClient,
		messageClient:     messageClient,
		messageLimit:      messageLimit,
	}
}

func (s *MessageService) RetireveUnsentMessages(ctx context.Context) ([]model.Message, error) {
	return s.messageRepository.RetrieveAll(ctx, model.Message{
		Status: model.MESSAGE_UNSENT,
	})
}

func (s *MessageService) StartMessageSending(ctx context.Context) {
	s.messageSender.Start(ctx, s.processMessages)
}

func (s *MessageService) StopMessageSending(ctx context.Context) {
	s.messageSender.Stop()
}

func (s *MessageService) processMessages() error {
	ctx := context.Background()

	messages, err := s.messageRepository.RetrieveMessagesForProcess(ctx, s.messageLimit)
	if err != nil {
		return fmt.Errorf("RetrieveMessagesForProcess error: %s", err.Error())
	}

	for _, msg := range messages {
		resp, err := s.messageClient.Send(msg.RecipientPhoneNumber, msg.Content)
		if err != nil {
			return fmt.Errorf("messageClient.Send error for message: %d, err: %s", msg.Id, err.Error())
		}

		err = s.redisClient.SetJson(ctx, resp.MessageId, model.MessageRedisPayload{
			MessageId: resp.MessageId,
			Timestamp: int(time.Now().Unix()),
		})
		if err != nil {
			return fmt.Errorf("redisClient.Set error for message: %d, err: %s", msg.Id, err.Error())
		}

		err = s.messageRepository.Update(ctx, &msg, map[string]interface{}{"status": model.MESSAGE_SENT})
		if err != nil {
			return fmt.Errorf("messageRepository.Update error for message: %d, err: %s", msg.Id, err.Error())
		}
	}

	return nil
}
