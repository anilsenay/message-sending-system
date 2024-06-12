package client_test

import (
	"testing"

	"github.com/anilsenay/message-sending-system/internal/client"
	"github.com/stretchr/testify/assert"
)

func TestNewMessageClient(t *testing.T) {
	url := "https://webhook.site/e72753b5-938a-4957-a12c-2c4f8b8c009c"
	c := client.NewMessageClient(url)
	resp, err := c.Send("+901111111111", "deneme")
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
}
