package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvFallback(t *testing.T) {
	value := GetEnv("ENV_KEY", "test")
	assert.Equal(t, "test", value)

	value_int := GetEnvInt("ENV_KEY", 100)
	assert.Equal(t, 100, value_int)

	value_int64 := GetEnvInt64("ENV_KEY", 9223372036854775807)
	assert.Equal(t, int64(9223372036854775807), value_int64)

	value_uint := GetEnvUInt("ENV_KEY", 100)
	assert.Equal(t, uint(100), value_uint)

	value_uint64 := GetEnvUInt64("ENV_KEY", 9223372036854775807)
	assert.Equal(t, uint64(9223372036854775807), value_uint64)

	value_bool := GetEnvBool("ENV_KEY", true)
	assert.Equal(t, true, value_bool)

	value_duration := GetEnvDuration("ENV_KEY", 100*time.Second)
	assert.Equal(t, 100*time.Second, value_duration)

	value_list := GetEnvStrList("ENV_KEY", ",", "abc, cde, fge,qwe,zxc")
	assert.Equal(t, 5, len(value_list))
	assert.Equal(t, "abc", value_list[0])
	assert.Equal(t, "cde", value_list[1])
	assert.Equal(t, "zxc", value_list[4])

	value_list_2 := GetEnvStrPtrList("ENV_KEY", ",", "abc, cde, fge,qwe,zxc")
	assert.Equal(t, 5, len(value_list_2))
	assert.Equal(t, "abc", *value_list_2[0])
	assert.Equal(t, "cde", *value_list_2[1])
	assert.Equal(t, "zxc", *value_list_2[4])

	int_value_list := GetEnvIntList("ENV_KEY", ",", "1,2,3,4,5")
	assert.Equal(t, 5, len(int_value_list))
	assert.Equal(t, 1, int_value_list[0])
	assert.Equal(t, 3, int_value_list[2])
	assert.Equal(t, 5, int_value_list[4])

	int_value_list_2 := GetEnvIntPtrList("ENV_KEY", ",", "1,2,3,4,5")
	assert.Equal(t, 5, len(int_value_list_2))
	assert.Equal(t, 1, *int_value_list_2[0])
	assert.Equal(t, 3, *int_value_list_2[2])
	assert.Equal(t, 5, *int_value_list_2[4])
}

func TestGetEnv(t *testing.T) {
	defer func() {
		os.Unsetenv("ENV_VAR")
		os.Unsetenv("ENV_VAR_INT")
		os.Unsetenv("ENV_VAR_INT64")
		os.Unsetenv("ENV_VAR_UINT")
		os.Unsetenv("ENV_VAR_UINT64")
		os.Unsetenv("ENV_VAR_BOOL")
		os.Unsetenv("ENV_VAR_DURATION")
		os.Unsetenv("ENV_VAR_LIST")
		os.Unsetenv("ENV_VAR_LIST_AS_INT")
	}()
	os.Setenv("ENV_VAR", "envValue")
	os.Setenv("ENV_VAR_INT", "123")
	os.Setenv("ENV_VAR_INT64", "9223372036854775807")
	os.Setenv("ENV_VAR_UINT", "4294967295")
	os.Setenv("ENV_VAR_UINT64", "9223372036854775807")
	os.Setenv("ENV_VAR_BOOL", "true")
	os.Setenv("ENV_VAR_DURATION", "10s")
	os.Setenv("ENV_VAR_LIST", "deneme;123;test;abcd")
	os.Setenv("ENV_VAR_LIST_AS_INT", "1;2;3")

	value := GetEnv("ENV_VAR", "fallbackValue")
	assert.Equal(t, "envValue", value)
	assert.NotEqual(t, "fallbackValue", value)

	value_int := GetEnvInt("ENV_VAR_INT", 1)
	assert.Equal(t, 123, value_int)
	assert.NotEqual(t, 1, value_int)

	value_int64 := GetEnvInt64("ENV_VAR_INT64", 1)
	assert.Equal(t, int64(9223372036854775807), value_int64)
	assert.NotEqual(t, 1, value_int64)

	value_uint := GetEnvUInt("ENV_VAR_UINT", 1)
	assert.Equal(t, uint(4294967295), value_uint)
	assert.NotEqual(t, -1, value_uint)

	value_uint64 := GetEnvUInt64("ENV_VAR_UINT64", 1)
	assert.Equal(t, uint64(9223372036854775807), value_uint64)
	assert.NotEqual(t, -1, value_uint64)

	value_bool := GetEnvBool("ENV_VAR_BOOL", false)
	assert.Equal(t, true, value_bool)
	assert.NotEqual(t, false, value_bool)

	value_duration := GetEnvDuration("ENV_VAR_DURATION", 1*time.Second)
	assert.Equal(t, 10*time.Second, value_duration)
	assert.NotEqual(t, 1*time.Second, value_duration)

	value_list := GetEnvStrList("ENV_VAR_LIST", ";", "abc;1;324")
	assert.Equal(t, "deneme", value_list[0])
	assert.NotEqual(t, "1", value_list[0])

	value_list_2 := GetEnvStrPtrList("ENV_VAR_LIST", ";", "abc;1;324")
	assert.Equal(t, "deneme", *value_list_2[0])
	assert.NotEqual(t, "1", *value_list_2[0])

	int_value_list := GetEnvIntList("ENV_VAR_LIST_AS_INT", ";", "1;2;3")
	assert.Equal(t, 1, int_value_list[0])
	assert.NotEqual(t, "1", int_value_list[0])

	int_value_list_2 := GetEnvIntPtrList("ENV_VAR_LIST_AS_INT", ";", "1;2;3")
	assert.Equal(t, 1, *int_value_list_2[0])
	assert.NotEqual(t, "1", *int_value_list_2[0])
}

func TestGetEnvInvalid(t *testing.T) {
	defer func() {
		os.Unsetenv("ENV_VAR_INT")
		os.Unsetenv("ENV_VAR_INT64")
		os.Unsetenv("ENV_VAR_UINT")
		os.Unsetenv("ENV_VAR_UINT64")
		os.Unsetenv("ENV_VAR_BOOL")
		os.Unsetenv("ENV_VAR_DURATION")
		os.Unsetenv("ENV_VAR_LIST")
		os.Unsetenv("ENV_VAR_LIST_AS_INT")
	}()
	os.Setenv("ENV_VAR_INT", "notInt")
	os.Setenv("ENV_VAR_INT64", "notInt64")
	os.Setenv("ENV_VAR_UINT", "notUInt")
	os.Setenv("ENV_VAR_UINT64", "notUInt64")
	os.Setenv("ENV_VAR_BOOL", "notBool")
	os.Setenv("ENV_VAR_DURATION", "notDuration")
	os.Setenv("ENV_VAR_LIST", "abc;def;gh")
	os.Setenv("ENV_VAR_LIST_AS_INT", "1;2;3")

	value_int := GetEnvInt("ENV_VAR_INT", 1)
	assert.NotEqual(t, 123, value_int)
	assert.Equal(t, 1, value_int)

	value_int64 := GetEnvInt64("ENV_VAR_INT64", 1)
	assert.NotEqual(t, 9223372036854775807, value_int64)
	assert.Equal(t, int64(1), value_int64)

	value_uint := GetEnvUInt("ENV_VAR_UINT", 1)
	assert.NotEqual(t, 4294967295, value_uint)
	assert.Equal(t, uint(1), value_uint)

	value_uint64 := GetEnvUInt64("ENV_VAR_UINT64", 1)
	assert.NotEqual(t, 9223372036854775807, value_uint64)
	assert.Equal(t, uint64(1), value_uint64)

	value_bool := GetEnvBool("ENV_VAR_BOOL", false)
	assert.NotEqual(t, true, value_bool)
	assert.Equal(t, false, value_bool)

	value_duration := GetEnvDuration("ENV_VAR_DURATION", 1*time.Second)
	assert.NotEqual(t, 10*time.Second, value_duration)
	assert.Equal(t, 1*time.Second, value_duration)

	value_list := GetEnvStrList("ENV_VAR_LIST", ",", "abc;def;gh") // , insted of ;
	assert.Equal(t, "abc;def;gh", value_list[0])
	assert.NotEqual(t, "abc", value_list[0])

	value_list_2 := GetEnvStrPtrList("ENV_VAR_LIST", ",", "abc;def;gh") // , insted of ;
	assert.Equal(t, "abc;def;gh", *value_list_2[0])
	assert.NotEqual(t, "abc", *value_list_2[0])
}

func TestGetEnvInt_Panic(t *testing.T) {
	assert.Panics(t, func() {
		_ = GetEnvIntList("ENV_VAR_LIST", ",", "abc,def,gh")
	})

	assert.Panics(t, func() {
		_ = GetEnvIntPtrList("ENV_VAR_LIST", ",", "abc,def,gh")
	})
}
