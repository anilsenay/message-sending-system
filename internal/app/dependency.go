package app

import (
	"github.com/anilsenay/message-sending-system/internal/client"
	"github.com/anilsenay/message-sending-system/internal/config"
	"github.com/anilsenay/message-sending-system/internal/handler"
	"github.com/anilsenay/message-sending-system/internal/repository"
	"github.com/anilsenay/message-sending-system/internal/service"
	"github.com/anilsenay/message-sending-system/internal/worker"
	"github.com/anilsenay/message-sending-system/pkg/orm"
	"github.com/anilsenay/message-sending-system/pkg/ticker"
	"github.com/redis/go-redis/v9"
)

var db = orm.NewDatabase(orm.DatabaseConfig{
	Host:              config.DB_HOST,
	Port:              config.DB_PORT,
	User:              config.DB_USER,
	Password:          config.DB_PASSWORD,
	Database:          config.DB_NAME,
	MaxConns:          config.DB_POOL_MAX_CONN,
	MaxConnLifeTime:   config.DB_POOL_MAX_CONN_LIFETIME,
	MaxConnIdleTime:   config.DB_POOL_MAX_CONN_IDLETIME,
	HealthCheckPeriod: config.DB_POOL_HEALTH_CHECK_PERIOD,
	LogLevel:          config.LOG_LEVEL,
})

var redisClient = client.NewRedis(
	redis.NewClient(&redis.Options{Addr: config.REDIS_HOST}),
)

var messageClient = client.NewMessageClient(config.WEBHOOK_URL, client.WithAuthKey(config.WEBHOOK_AUTH_KEY))

var msender = worker.NewMessageSender(ticker.NewTimeTicker(), config.WORKER_PERIOD)

var messageRepo = repository.NewMessageRepository(db)
var messageService = service.NewMessageService(messageRepo, msender, redisClient, messageClient, config.WORKER_MESSAGE_LIMIT)
var messageHandler = handler.NewMessageHandler(messageService)
