package common

import (
	"context"
	"os"
)

var (
	Ctx context.Context
)

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
