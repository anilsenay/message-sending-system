package repository_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/anilsenay/message-sending-system/internal/model"
	"github.com/anilsenay/message-sending-system/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestMessageRepository_GetMessagesForProcess(t *testing.T) {
	repo := repository.NewMessageRepository(dockerDatabase)

	numOfMessagesToProcess := 2
	created := []model.Message{}
	t.Run("1. Create some messages", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			message := model.Message{
				Content:              fmt.Sprintf("content-%d", i),
				RecipientPhoneNumber: "+90111111111",
				Status:               model.MESSAGE_UNSENT,
				CreatedAt:            time.Now().Add(time.Duration(i+1) * time.Second),
			}
			result := dockerDatabase.GetConnection().Create(&message)
			assert.NoError(t, result.Error)
			created = append(created, message)
		}

		assert.NotEmpty(t, created)
	})

	t.Run("2. Get all messages", func(t *testing.T) {
		messages, err := repo.RetrieveAll(context.Background(), model.Message{})
		assert.NoError(t, err)
		assert.Len(t, messages, len(created))
	})

	t.Run("3. Retrieve messages for processing", func(t *testing.T) {
		messages, err := repo.RetrieveMessagesForProcess(context.Background(), numOfMessagesToProcess)
		assert.NoError(t, err)
		assert.Len(t, messages, numOfMessagesToProcess)
		assert.Equal(t, model.MESSAGE_PROCESSING, messages[0].Status)
		assert.Equal(t, "content-0", messages[0].Content)
	})

	t.Run("4. Get all unsent messages", func(t *testing.T) {
		messages, err := repo.RetrieveAll(context.Background(), model.Message{Status: model.MESSAGE_UNSENT})
		assert.NoError(t, err)
		assert.Len(t, messages, len(created)-numOfMessagesToProcess)
	})
}
