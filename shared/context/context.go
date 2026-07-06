package context

import (
	"context"

	"github.com/google/uuid"
)

type traceIDKeyType struct{}

var TraceIDKey = traceIDKeyType{}

func GetTraceID(ctx context.Context) uuid.UUID {
	traceID, ok := ctx.Value(TraceIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil
	}
	return traceID
}
