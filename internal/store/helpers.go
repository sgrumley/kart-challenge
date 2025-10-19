package store

import (
	"time"

	"github.com/google/uuid"
)

func TimeStampNow() int {
	return int(time.Now().UTC().UnixNano())
}

func GenerateUUIDv4() uuid.UUID {
	return uuid.New()
}
