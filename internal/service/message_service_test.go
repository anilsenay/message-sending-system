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

type mockRepository struct{}

func (mockRepository) RetrieveAll(ctx context.Context, filters model.Message) ([]model.Message, error) {
	if filters.Status == model.MESSAGE_SENT {
		return []model.Message{
			{Id: 1, Status: model.MESSAGE_SENT, CreatedAt: time.Now()},
			{Id: 2, Status: model.MESSAGE_SENT, CreatedAt: time.Now()},
			{Id: 3, Status: model.MESSAGE_SENT, CreatedAt: time.Now()},
		}, nil
	}
	return []model.Message{}, nil
}
func (mockRepository) RetrieveMessagesForProcess(ctx context.Context, limit int) ([]model.Message, error) {
	return []model.Message{
		{Id: 4, Status: model.MESSAGE_UNSENT, CreatedAt: time.Now()},
		{Id: 5, Status: model.MESSAGE_UNSENT, CreatedAt: time.Now()},
	}, nil
}
func (mockRepository) Update(ctx context.Context, m *model.Message, updates map[string]interface{}) error {
	m.Status = updates["status"].(model.MessageStatus)
	return nil
}
func (mockRepository) Create(ctx context.Context, m *model.Message) error {
	return nil
}

type mockRedis struct{}

func (mockRedis) SetJson(ctx context.Context, key string, val interface{}) error {
	return nil
}

type mockMessageClient struct{}

func (mockMessageClient) Send(to, content string) (model.MessageClientResponse, error) {
	return model.MessageClientResponse{Message: "Accepted", MessageId: "123"}, nil
}

func TestNewMessageService(t *testing.T) {
	ms := worker.NewMessageSender(ticker.NewTimeTicker(), time.Second)
	s := service.NewMessageService(mockRepository{}, ms, mockRedis{}, mockMessageClient{}, 2)

	t.Run("1. Retireve Sent Messages", func(t *testing.T) {
		messages, err := s.RetireveSentMessages(context.Background())
		assert.NoError(t, err)
		assert.Len(t, messages, 3)
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
			s.StopMessageSending()
			time.Sleep(1 * time.Second)
			assert.Equal(t, false, ms.IsRunning())
		})
	})
}
