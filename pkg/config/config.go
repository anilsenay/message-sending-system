package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}

func GetEnvInt64(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return fallback
}

func GetEnvUInt(key string, fallback uint) uint {
	if value, ok := os.LookupEnv(key); ok {
		if intValue, err := strconv.ParseUint(value, 10, 32); err == nil {
			return uint(intValue)
		}
	}
	return fallback
}

func GetEnvUInt64(key string, fallback uint64) uint64 {
	if value, ok := os.LookupEnv(key); ok {
		if intValue, err := strconv.ParseUint(value, 10, 64); err == nil {
			return intValue
		}
	}
	return fallback
}

func GetEnvBool(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return fallback
}

func GetEnvDuration(key string, fallback time.Duration) time.Duration {
	if value, ok := os.LookupEnv(key); ok {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return fallback
}

func GetEnvStrList(key, delimeter, fallback string) []string {
	var parsedArray []string
	if value, ok := os.LookupEnv(key); ok {
		parsedArray = strings.Split(value, delimeter)
	} else {
		parsedArray = strings.Split(fallback, delimeter)
	}

	var list = make([]string, 0, len(parsedArray))

	for _, item := range parsedArray {
		itemWithoutSpaces := strings.Trim(item, " ")
		list = append(list, itemWithoutSpaces)
	}

	return list
}

func GetEnvStrPtrList(key, delimeter, fallback string) []*string {
	var parsedArray []string
	if value, ok := os.LookupEnv(key); ok {
		parsedArray = strings.Split(value, delimeter)
	} else {
		parsedArray = strings.Split(fallback, delimeter)
	}

	var list []*string = make([]*string, len(parsedArray))

	for index, item := range parsedArray {
		itemWithoutSpaces := strings.Trim(item, " ")
		list[index] = &itemWithoutSpaces
	}

	return list
}

func GetEnvIntList(key, delimeter, fallback string) []int {
	var parsedArray []string
	if value, ok := os.LookupEnv(key); ok {
		parsedArray = strings.Split(value, delimeter)
	} else {
		parsedArray = strings.Split(fallback, delimeter)
	}

	var list = make([]int, 0, len(parsedArray))

	for _, item := range parsedArray {
		itemWithoutSpaces := strings.Trim(item, " ")
		itemAsInt, err := strconv.Atoi(itemWithoutSpaces)
		if err != nil {
			log.Panic().Msg("Invalid integer list element in environment variable")
		}
		list = append(list, itemAsInt)
	}

	return list
}

func GetEnvIntPtrList(key, delimeter, fallback string) []*int {
	var parsedArray []string
	if value, ok := os.LookupEnv(key); ok {
		parsedArray = strings.Split(value, delimeter)
	} else {
		parsedArray = strings.Split(fallback, delimeter)
	}

	var list = make([]*int, 0, len(parsedArray))

	for _, item := range parsedArray {
		itemWithoutSpaces := strings.Trim(item, " ")
		itemAsInt, err := strconv.Atoi(itemWithoutSpaces)
		if err != nil {
			log.Panic().Msg("Invalid integer list element in environment variable")
		}
		list = append(list, &itemAsInt)
	}

	return list
}
