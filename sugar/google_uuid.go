package sugar

import (
	"context"

	"github.com/google/uuid"
)

func GetUUID() string {
	return uuid.New().String()
}

func NewContext() context.Context {
	ctx := context.WithValue(context.Background(), "traceId", GetUUID())
	return ctx
}
