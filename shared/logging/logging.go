package logging

import (
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func FailedParseJSON(logger *zap.SugaredLogger, traceID uuid.UUID, err error) {
	logger.Errorw("failed to parse JSON data", "trace_id", traceID, "error", err)
}

func FailedValidateData(logger *zap.SugaredLogger, traceID uuid.UUID, err error) {
	logger.Warnw("failed to validate data", "trace_id", traceID, "error", err)
}

func FailedWriteResponse(logger *zap.SugaredLogger, traceID uuid.UUID, err error) {
	logger.Errorw("failed to write response", "trace_id", traceID, "error", err)
}

func Failed(logger *zap.SugaredLogger, traceID uuid.UUID, what string, err error) {
	logger.Errorw(fmt.Sprintf("%s failed", what), "trace_id", traceID, "error", err)
}

func Declined(logger *zap.SugaredLogger, traceID uuid.UUID, what string, err error) {
	logger.Warnw(fmt.Sprintf("%s declined", what), "trace_id", traceID, "error", err)
}

func Success(logger *zap.SugaredLogger, traceID uuid.UUID, what string, fields ...any) {
	args := append([]any{"trace_id", traceID}, fields...)
	logger.Infow(fmt.Sprintf("%s successfully completed", what), args...)
}
