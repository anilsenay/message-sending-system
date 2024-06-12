package service

import (
	"context"
	"fmt"
	"time"

	"github.com/anilsenay/message-sending-system/internal/model"
	"github.com/anilsenay/message-sending-system/pkg/logger"
)

type messageRepository interface {
	RetrieveAll(ctx context.Context, filters model.DbFilters[model.Message]) ([]model.Message, error)
	RetrieveMessagesForProcess(ctx context.Context, limit int) ([]model.Message, error)
	Update(ctx context.Context, model *model.Message, updates map[string]interface{}) error
	Create(ctx context.Context, model *model.Message) error
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

func (s *MessageService) RetireveSentMessages(ctx context.Context) ([]model.Message, error) {
	order := "sent_at DESC"

	return s.messageRepository.RetrieveAll(ctx, model.DbFilters[model.Message]{
		Order: &order,
		Model: &model.Message{
			Status: model.MESSAGE_SENT,
		},
	})
}

func (s *MessageService) CreateMessage(ctx context.Context, m *model.Message) error {
	m.CreatedAt = time.Now()
	m.Status = model.MESSAGE_UNSENT
	return s.messageRepository.Create(ctx, m)
}

func (s *MessageService) StartMessageSending(ctx context.Context) {
	s.messageSender.Start(ctx, s.processMessages)
}

func (s *MessageService) StopMessageSending() {
	s.messageSender.Stop()
}

func (s *MessageService) processMessages() error {
	ctx := context.Background()

	messages, err := s.messageRepository.RetrieveMessagesForProcess(ctx, s.messageLimit)
	if err != nil {
		return fmt.Errorf("RetrieveMessagesForProcess error: %s", err.Error())
	}

	logger.Info().Msgf("%d messages found", len(messages))

	for _, msg := range messages {
		err := s.processMessage(ctx, msg)
		if err != nil {
			logger.Error().Msgf("error while processing message: %d, err: %s", msg.Id, err.Error())

			err = s.messageRepository.Update(ctx, &msg, map[string]interface{}{"status": model.MESSAGE_FAILED})
			if err != nil {
				logger.Error().Msgf("error while updating message: %d as failed, err: %s", msg.Id, err.Error())
			}
		}
	}

	return nil
}

func (s *MessageService) processMessage(ctx context.Context, msg model.Message) error {
	resp, err := s.messageClient.Send(msg.RecipientPhoneNumber, msg.Content)
	if err != nil {
		return fmt.Errorf("messageClient.Send error for message: %d, err: %s", msg.Id, err.Error())
	}
	sentDate := time.Now()

	err = s.redisClient.SetJson(ctx, fmt.Sprintf("messages:%s", resp.MessageId), model.MessageRedisPayload{
		MessageId: resp.MessageId,
		Timestamp: int(sentDate.Unix()),
	})
	if err != nil {
		return fmt.Errorf("redisClient.Set error for message: %d, uuid:%s, err: %s", msg.Id, resp.MessageId, err.Error())
	}

	err = s.messageRepository.Update(ctx, &msg, map[string]interface{}{
		"status":  model.MESSAGE_SENT,
		"sent_at": sentDate,
	})
	if err != nil {
		return fmt.Errorf("messageRepository.Update error for message: %d, uuid:%s, err: %s", msg.Id, resp.MessageId, err.Error())
	}

	logger.Debug().Msgf("message with response uuid: %s processed", resp.MessageId)
	return nil
}
