package handler_test

import (
	"testing"
	"time"

	"github.com/anilsenay/message-sending-system/internal/model"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewMessageHandler(t *testing.T) {
	var app = fiber.New()
	messageHandler.SetRoutes(app)

	messageCount := 10

	t.Run("1. Create messages", func(t *testing.T) {
		for i := 0; i < messageCount; i++ {
			msg := model.MessageCreateRequest{Content: "Content", RecipientPhoneNumber: "+905551112233"}
			status, resp, fail := testRequest[model.Message](t, app, "POST", "/messages", msg, nil)
			assert.Equal(t, 200, status)
			assert.Empty(t, fail)
			assert.NotEmpty(t, resp)
			assert.Equal(t, msg.Content, resp.Content)
		}
	})

	t.Run("2. Get sent messages", func(t *testing.T) {
		status, resp, fail := testRequest[[]model.Message](t, app, "GET", "/messages", nil, nil)
		assert.Equal(t, 200, status)
		assert.Empty(t, fail)
		assert.Empty(t, resp)
		assert.Len(t, resp, 0)
	})

	t.Run("3. Start worker", func(t *testing.T) {
		status, resp, fail := testRequest[string](t, app, "POST", "/messages/worker/start", nil, nil)
		assert.Equal(t, 200, status)
		assert.Empty(t, fail)
		assert.NotEmpty(t, resp)
		assert.Equal(t, "OK", resp)
	})

	t.Run("4. Wait for messages to process", func(t *testing.T) {
		// because of network delays, sleep time is more than it should be
		time.Sleep(15 * time.Second)
	})

	t.Run("5. Get sent messages", func(t *testing.T) {
		status, resp, fail := testRequest[[]model.Message](t, app, "GET", "/messages", nil, nil)
		assert.Equal(t, 200, status)
		assert.Empty(t, fail)
		assert.NotEmpty(t, resp)
		assert.Len(t, resp, messageCount)
	})

	t.Run("6. Stop worker", func(t *testing.T) {
		status, resp, fail := testRequest[string](t, app, "POST", "/messages/worker/stop", nil, nil)
		assert.Equal(t, 200, status)
		assert.Empty(t, fail)
		assert.NotEmpty(t, resp)
		assert.Equal(t, "OK", resp)
	})
}
