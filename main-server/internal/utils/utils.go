package utils

import (
	"github.com/google/uuid"
)

func GenerateRequestId() string {
	return uuid.New().String()
}
