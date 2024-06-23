package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/anilsenay/message-sending-system/internal/client"
	"github.com/anilsenay/message-sending-system/internal/handler"
	"github.com/anilsenay/message-sending-system/internal/repository"
	"github.com/anilsenay/message-sending-system/internal/service"
	"github.com/anilsenay/message-sending-system/internal/worker"
	"github.com/anilsenay/message-sending-system/pkg/orm"
	"github.com/anilsenay/message-sending-system/pkg/ticker"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	testcontainers_redis "github.com/testcontainers/testcontainers-go/modules/redis"
)

var dockerDatabase *orm.Database

var messageHandler *handler.MessageHandler

var messageCountPerProcess = 5

func TestMain(m *testing.M) {
	dockerDatabase = orm.NewDockerDatabase(orm.DockerDatabaseConfig{
		MigrationQueryPath: "../../db.sql",
	})

	redisContainer, err := testcontainers_redis.RunContainer(context.Background(),
		testcontainers.WithImage("docker.io/redis:7"),
		testcontainers_redis.WithSnapshotting(10, 1),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}
	redisUrl, _ := redisContainer.ConnectionString(context.Background())
	redisClient := client.NewRedis(redis.NewClient(&redis.Options{Addr: strings.ReplaceAll(redisUrl, "redis://", "")}))

	messageClient := client.NewMessageClient("https://webhook.site/9c867dd2-b25e-446f-accc-bef9988fc035")

	msender := worker.NewMessageSender(ticker.NewTimeTicker(), 5*time.Second)

	messageRepo := repository.NewMessageRepository(dockerDatabase)
	messageService := service.NewMessageService(messageRepo, msender, redisClient, messageClient, messageCountPerProcess)
	messageHandler = handler.NewMessageHandler(messageService)

	os.Exit(m.Run())
}

func testRequest[T any](t *testing.T, app *fiber.App, method string, url string, body any, headers map[string]string) (statusCode int, success T, fail handler.FailDetails) {
	var resBody struct {
		Status  int                 `json:"status"`
		Success T                   `json:"success"`
		Fail    handler.FailDetails `json:"fail"`
	}

	reqBody, err := json.Marshal(body)
	assert.NoError(t, err)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	assert.NoError(t, err)

	req.Header.Add("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	res, err := app.Test(req, 5000)
	assert.NoError(t, err)
	jsonDataFromHttp, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	err = json.Unmarshal(jsonDataFromHttp, &resBody)
	assert.NoError(t, err)
	assert.NotEmpty(t, resBody)
	return res.StatusCode, resBody.Success, resBody.Fail
}
