package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/anilsenay/message-sending-system/internal/model"
	"github.com/anilsenay/message-sending-system/internal/service"
	"github.com/anilsenay/message-sending-system/internal/worker"
	"github.com/anilsenay/message-sending-system/pkg/ticker"
	"github.com/stretchr/testify/assert"
)

var msgs = []model.Message{
	{Id: 1, Status: model.MESSAGE_UNSENT, CreatedAt: time.Now()},
	{Id: 2, Status: model.MESSAGE_UNSENT, CreatedAt: time.Now()},
	{Id: 3, Status: model.MESSAGE_UNSENT, CreatedAt: time.Now()},
}

type mockRepository struct{}

func (mockRepository) RetrieveAll(ctx context.Context, filters model.Message) ([]model.Message, error) {
	return msgs, nil
}
func (mockRepository) RetrieveMessagesForProcess(ctx context.Context, limit int) ([]model.Message, error) {
	return msgs[0:2], nil
}
func (mockRepository) Update(ctx context.Context, m *model.Message, updates map[string]interface{}) error {
	m.Status = updates["status"].(model.MessageStatus)
	return nil
}

type mockRedis struct{}

func (mockRedis) Set(key, value string) error {
	return nil
}

type mockMessageClient struct{}

func (mockMessageClient) Send(to, content string) (model.MessageResponse, error) {
	return model.MessageResponse{Message: "Accepted", MessageId: "123"}, nil
}

func TestNewMessageService(t *testing.T) {
	ms := worker.NewMessageSender(ticker.NewTimeTicker(), time.Second)
	s := service.NewMessageService(mockRepository{}, ms, mockRedis{}, mockMessageClient{}, 2)

	t.Run("1. Retireve Unsent Messages", func(t *testing.T) {
		messages, err := s.RetireveUnsentMessages(context.Background())
		assert.NoError(t, err)
		assert.Len(t, messages, len(msgs))
	})

	t.Run("2. Start processing", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		s.StartMessageSending(ctx)
		time.Sleep(6 * time.Second)
		assert.Equal(t, false, ms.IsRunning())
	})

	t.Run("3. Stop processing", func(t *testing.T) {
		t.Run("3.1. Start", func(t *testing.T) {
			s.StartMessageSending(context.Background())
			time.Sleep(1 * time.Second)
			assert.Equal(t, true, ms.IsRunning())
		})

		t.Run("3.2. Stop", func(t *testing.T) {
			s.StopMessageSending(context.Background())
			time.Sleep(1 * time.Second)
			assert.Equal(t, false, ms.IsRunning())
		})
	})
}
