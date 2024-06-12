package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/anilsenay/message-sending-system/internal/model"
)

type messageBody struct {
	To      string `json:"to"`
	Content string `json:"content"`
}

type MessageClient struct {
	url     string
	authKey string
	timeout time.Duration
}

type option func(*MessageClient)

func WithTimeout(t time.Duration) option {
	return func(mc *MessageClient) {
		mc.timeout = t
	}
}

func WithAuthKey(key string) option {
	return func(mc *MessageClient) {
		mc.authKey = key
	}
}

func NewMessageClient(url string, opts ...option) *MessageClient {
	mc := &MessageClient{url: url}

	for _, opt := range opts {
		opt(mc)
	}

	return mc
}

func (mc *MessageClient) Send(to, content string) (model.MessageClientResponse, error) {
	body := messageBody{To: to, Content: content}
	bodyByteArr, err := json.Marshal(body)
	if err != nil {
		return model.MessageClientResponse{}, err
	}

	req, err := http.NewRequest(http.MethodPost, mc.url, bytes.NewBuffer(bodyByteArr))
	if err != nil {
		return model.MessageClientResponse{}, fmt.Errorf("could not create request: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-ins-auth-key", mc.authKey)

	client := http.Client{Timeout: mc.timeout}

	res, err := client.Do(req)
	if err != nil {
		return model.MessageClientResponse{}, fmt.Errorf("could not send request: %s", err)
	}

	respBody, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return model.MessageClientResponse{}, fmt.Errorf("could not read response body: %s", err)
	}

	messageResponse := model.MessageClientResponse{}
	err = json.Unmarshal(respBody, &messageResponse)
	if err != nil {
		return model.MessageClientResponse{}, fmt.Errorf("could not unmarshal response body: %s", err)
	}

	if messageResponse.MessageId == "" {
		return model.MessageClientResponse{}, fmt.Errorf("messageId could not received")
	}

	return messageResponse, nil
}
