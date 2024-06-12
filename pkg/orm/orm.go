package orm

import (
	"fmt"
	"time"

	"github.com/anilsenay/message-sending-system/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	gormLogger "gorm.io/gorm/logger"
)

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string

	MaxConns int

	MaxConnLifeTime   time.Duration
	MaxConnIdleTime   time.Duration
	HealthCheckPeriod time.Duration

	LogLevel string
}

type Database struct {
	Db *gorm.DB
}

var ErrRecordNotFound = gorm.ErrRecordNotFound
var ClauseReturning = clause.Returning{}

func NewDatabase(config DatabaseConfig) *Database {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d", config.Host, config.User, config.Password, config.Database, config.Port)
	db, err := gorm.Open(postgres.New(
		postgres.Config{
			DSN: dsn,
		},
	), &gorm.Config{
		Logger: getLogger(config.LogLevel),
	})
	if err != nil {
		logger.Panic().Msgf("Failed to connect database: %s", err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Panic().Msgf("Failed to connect database: %s", err.Error())
	}

	if config.MaxConnLifeTime != 0 {
		sqlDB.SetConnMaxLifetime(config.MaxConnLifeTime)
	}
	if config.MaxConnIdleTime != 0 {
		sqlDB.SetConnMaxIdleTime(config.MaxConnIdleTime)
	}
	if config.MaxConns != 0 {
		sqlDB.SetMaxOpenConns(int(config.MaxConns))
		sqlDB.SetMaxIdleConns(int(config.MaxConns))
	}
	if config.HealthCheckPeriod != 0 {
		sqlDB.SetConnMaxLifetime(config.HealthCheckPeriod)
	}

	return &Database{Db: db}
}

func (d *Database) GetConnection() *gorm.DB {
	return d.Db
}

func (d *Database) Ping() error {
	sqlDB, err := d.Db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func (d *Database) Close() {
	sqlDB, err := d.Db.DB()
	if err != nil {
		logger.Panic().Msgf("Failed to connect database: %s", err.Error())
	}
	sqlDB.Close()
}

func getLogger(level string) gormLogger.Interface {
	switch level {
	case "DEBUG":
		return gormLogger.Default.LogMode(gormLogger.Info)
	case "ERROR":
		return gormLogger.Default.LogMode(gormLogger.Error)
	default:
		return gormLogger.Default.LogMode(gormLogger.Silent)
	}
}
