package client_test

import (
	"testing"

	"github.com/anilsenay/message-sending-system/internal/client"
	"github.com/stretchr/testify/assert"
)

func TestNewMessageClient(t *testing.T) {
	url := "https://webhook.site/9c867dd2-b25e-446f-accc-bef9988fc035"
	c := client.NewMessageClient(url)
	resp, err := c.Send("+901111111111", "deneme")
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
}
