package utils

import (
	"github.com/google/uuid"
	"time"
)

func GenerateUUID() string {
	return uuid.New().String()
}

func IsUUIDExpired(timestamp int64) bool {
	return time.Now().Unix()-timestamp > 86400
}
