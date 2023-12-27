package utils

import (
	"fmt"
	"os"
	"strconv"
)

type Env interface {
	uint | string | bool
}

func GetEnv[T Env](key string, defaultValue string, required bool) T {
	var val T

	value, ok := os.LookupEnv(key)
	if !ok {
		if required {
			panic(fmt.Sprintf("env variable %s is required", key))
		}

		value = defaultValue
	}

	switch ptr := any(&val).(type) {
	case *string:
		*ptr = value
	case *uint:
		target, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			panic(fmt.Sprintf("env variable %s=%v can not be paresed to int", key, value))
		}

		*ptr = uint(target)
	case *bool:
		target, err := strconv.ParseBool(value)
		if err != nil {
			panic(fmt.Sprintf("env variable %s=%v can not be paresed to bool", key, value))
		}

		*ptr = target
	}

	return val
}
