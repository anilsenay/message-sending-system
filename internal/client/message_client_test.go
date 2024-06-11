package client_test

import (
	"testing"

	"github.com/anilsenay/message-sending-system/internal/client"
	"github.com/stretchr/testify/assert"
)

func TestNewMessageClient(t *testing.T) {
	url := "https://webhook.site/9084ab6a-4827-45ce-b8fc-1fa50f3cbf10"
	c := client.NewMessageClient(url)
	resp, err := c.Send("+901111111111", "deneme")
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
}
