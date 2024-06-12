package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.With().Caller().Logger()
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}

func SetLogLevel(level string) {
	switch level {
	case "ERROR":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "INFO":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "WARN":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "DEBUG":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "DISABLED":
		zerolog.SetGlobalLevel(zerolog.Disabled)
	default:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

func Info() *zerolog.Event {
	return log.Info()
}

func Fatal() *zerolog.Event {
	return log.Fatal()
}

func Error() *zerolog.Event {
	return log.Error()
}

func Warn() *zerolog.Event {
	return log.Warn()
}

func Panic() *zerolog.Event {
	return log.Panic()
}

func Debug() *zerolog.Event {
	return log.Debug()
}

func Trace() *zerolog.Event {
	return log.Trace()
}

func Log() *zerolog.Event {
	return log.Log()
}
