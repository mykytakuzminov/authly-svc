package context

import (
	"context"

	"github.com/google/uuid"
)

type traceIDKeyType struct{}
type userIDKeyType struct{}
type userRoleKeyType struct{}

var TraceIDKey = traceIDKeyType{}
var UserIDKey = userIDKeyType{}
var UserRoleKey = userRoleKeyType{}

func GetTraceID(ctx context.Context) uuid.UUID {
	traceID, ok := ctx.Value(TraceIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil
	}
	return traceID
}

func GetUserID(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
	return userID, ok
}

func GetUserRole(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(UserRoleKey).(string)
	return role, ok
}
