package logger

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	Key = "logger"

	RequestID = "request_id"
)

type Logger struct {
	l *zap.Logger
}

func New(ctx context.Context) (context.Context, error) {

	logger, err := zap.NewDevelopment()
	if err != nil {
		return ctx, err
	}

	return context.WithValue(ctx, Key, &Logger{logger}), nil
}

func GetLoggerFromCtx(ctx context.Context) *Logger {
	return ctx.Value(Key).(*Logger)
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx.Value(RequestID) != nil {
		fields = append(fields, zap.String(RequestID, ctx.Value(RequestID).(string)))
	}
	l.l.Info(msg, fields...)
}

func (l *Logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx.Value(RequestID) != nil {
		fields = append(fields, zap.String(RequestID, ctx.Value(RequestID).(string)))
	}
	l.l.Fatal(msg, fields...)
}

func (l *Logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx.Value(RequestID) != nil {
		fields = append(fields, zap.String(RequestID, ctx.Value(RequestID).(string)))
	}
	l.l.Debug(msg, fields...)
}

func Interceptor(ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	next grpc.UnaryHandler,
) (any, error) {

	guid := uuid.New().String()
	ctx = context.WithValue(ctx, RequestID, guid)

	if ctx.Value(Key) == nil {
		var err error
		ctx, err = New(ctx)
		if err != nil {
			return nil, err
		}
	}

	GetLoggerFromCtx(ctx).Info(ctx,
		"request", zap.String("method", info.FullMethod),
		zap.Time("request time", time.Now()),
	)

	return next(ctx, req)
}
