package common

import (
	"bufio"
	"context"
	"os"
	"strings"
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

func ReadTrimString(reader *bufio.Reader) (string, error) {
	str, err := reader.ReadString('\n')
	return strings.Split(strings.TrimSpace(str), "\n")[0], err
}
