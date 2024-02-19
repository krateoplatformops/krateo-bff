package env

import (
	"os"
	"strconv"
	"strings"
	"time"
)

func String(key, defaultValue string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return val
}

func Int(key string, defaultValue int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	n := strings.TrimSpace(val)
	if len(n) == 0 {
		return defaultValue
	}

	res, err := strconv.Atoi(n)
	if err != nil {
		return defaultValue
	}
	return res
}

func Bool(key string, defaultValue bool) bool {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	res, err := strconv.ParseBool(strings.TrimSpace(val))
	if err != nil {
		return defaultValue
	}
	return res
}

func Duration(key string, defaultValue time.Duration) time.Duration {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	res, err := time.ParseDuration(strings.TrimSpace(val))
	if err != nil {
		return defaultValue
	}
	return res
}

func ServicePort(key string, defaultValue int) int {
	// "KRATEO_BFF_PORT=tcp://10.96.234.180:8081"
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	n := strings.TrimSpace(val)
	if len(n) == 0 {
		return defaultValue
	}

	if strings.HasPrefix(n, "tcp://") {
		idx := strings.LastIndexByte(n, ':')
		if idx < 0 || idx == len(n)-1 {
			return defaultValue
		}
		n = n[idx+1:]
	}

	res, err := strconv.Atoi(n)
	if err != nil {
		return defaultValue
	}
	return res
}
