package config

import (
	"time"

	"github.com/anilsenay/message-sending-system/pkg/config"
)

// app config
var APP_NAME = config.GetEnv("APP_NAME", "MESSAGE_SENDING_SERVICE")
var APP_VERSION = config.GetEnv("APP_VERSION", "0.1.0")
var LOG_LEVEL = config.GetEnv("LOG_LEVEL", "INFO")

// server config
var SERVER_HOST = config.GetEnv("SERVER_HOST", "localhost")
var SERVER_PORT = config.GetEnvInt("SERVER_PORT", 8080)

// worker
var WORKER_PERIOD = config.GetEnvDuration("WORKER_PERIOD", 2*time.Minute)
var WORKER_MESSAGE_LIMIT = config.GetEnvInt("WORKER_MESSAGE_LIMIT", 2)

// database
var DB_HOST = config.GetEnv("DB_HOST", "localhost")
var DB_PORT = config.GetEnvInt("DB_PORT", 5432)
var DB_USER = config.GetEnv("DB_USER", "user")
var DB_PASSWORD = config.GetEnv("DB_PASSWORD", "password")
var DB_NAME = config.GetEnv("DB_NAME", "database")
var DB_POOL_MIN_CONN = config.GetEnvInt("DB_POOL_MIN_CONN", 3)
var DB_POOL_MAX_CONN = config.GetEnvInt("DB_POOL_MAX_CONN", 5)
var DB_POOL_MAX_CONN_LIFETIME = config.GetEnvDuration("DB_POOL_MAX_CONN_LIFETIME", 24*time.Hour)
var DB_POOL_MAX_CONN_IDLETIME = config.GetEnvDuration("DB_POOL_MAX_CONN_IDLETIME", 24*time.Hour)
var DB_POOL_HEALTH_CHECK_PERIOD = config.GetEnvDuration("DB_POOL_HEALTH_CHECK_PERIOD", 1*time.Minute)

// redis
var REDIS_HOST = config.GetEnv("REDIS_HOST", "localhost:6379")
var REDIS_DB = config.GetEnvInt("REDIS_DB", 0)

// webhook
var WEBHOOK_URL = config.GetEnv("WEBHOOK_URL", "https://webhook.site/9c867dd2-b25e-446f-accc-bef9988fc035")
var WEBHOOK_AUTH_KEY = config.GetEnv("WEBHOOK_AUTH_KEY", "")
