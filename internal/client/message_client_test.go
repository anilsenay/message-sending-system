package client_test

import (
	"testing"

	"github.com/anilsenay/message-sending-system/internal/client"
	"github.com/stretchr/testify/assert"
)

func TestNewMessageClient(t *testing.T) {
	url := "https://webhook.site/da02ccdc-d02d-4a41-89ac-4938daca524e"
	c := client.NewMessageClient(url)
	resp, err := c.Send("+901111111111", "deneme")
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
}
