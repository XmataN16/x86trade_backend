package utils

import (
	"os"
	"strconv"
)

// getEnvInt возвращает значение окружения как int, или дефолт, если пусто/непарсится
func GetEnvInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return i
}
