package lax

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapAdapter implements Logger and pgx.Logger interfaces for zap logger.
type ZapAdapter struct {
	l *zap.Logger
}

// NewZapAdapter creates zap logger adapter.
func NewZapAdapter(l *zap.Logger) *ZapAdapter {
	return &ZapAdapter{l: l}
}

// NewDefaultZapAdapter creates zap logger adapter using ectobit.com proposed configuration.
//
// format may be either console or json.
// level may be one of debug, info, warn or error. For level debug and format console output will be colored.
//
// Logger will show stack trace just in error case.
func NewDefaultZapAdapter(format, level string) (*ZapAdapter, error) {
	log, err := NewZapLogger(format, level, 1)
	if err != nil {
		return nil, fmt.Errorf("new zap logger: %w", err)
	}

	return NewZapAdapter(log), nil
}

// Debug sends a message to logger at debug level.
func (l *ZapAdapter) Debug(message string, fields ...Field) {
	l.l.Debug(message, l.toZapFields(fields)...)
}

// Info sends a message to logger at info level.
func (l *ZapAdapter) Info(message string, fields ...Field) {
	l.l.Info(message, l.toZapFields(fields)...)
}

// Warn sends a message to logger at warn level.
func (l *ZapAdapter) Warn(message string, fields ...Field) {
	l.l.Warn(message, l.toZapFields(fields)...)
}

// Error sends a message to logger at error level.
func (l *ZapAdapter) Error(message string, fields ...Field) {
	l.l.Error(message, l.toZapFields(fields)...)
}

// Flush flushes buffer ignoring eventual error.
func (l *ZapAdapter) Flush() {
	_ = l.l.Sync()
}

// Log is the interface used to get logging from pgx internals.
func (l *ZapAdapter) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	fields := make([]zapcore.Field, 0, len(data))

	for k, v := range data {
		fields = append(fields, zap.Any(k, v))
	}

	switch level {
	case pgx.LogLevelTrace:
		l.l.Debug(msg, append(fields, zap.Stringer("PGX_LOG_LEVEL", level))...)
	case pgx.LogLevelDebug:
		l.l.Debug(msg, fields...)
	case pgx.LogLevelInfo:
		l.l.Info(msg, fields...)
	case pgx.LogLevelWarn:
		l.l.Warn(msg, fields...)
	case pgx.LogLevelError:
		l.l.Error(msg, fields...)
	default:
		l.l.Error(msg, append(fields, zap.Stringer("PGX_LOG_LEVEL", level))...)
	}
}

func (l *ZapAdapter) toZapFields(fields []Field) []zapcore.Field { //nolint:cyclop
	zfs := make([]zapcore.Field, 0, len(fields))

	for _, field := range fields {
		switch field.vType {
		case tAny:
			zfs = append(zfs, zap.Any(field.key, field.value))
		case tError:
			zfs = append(zfs, zap.Error(field.value.(error))) //nolint:forcetypeassert
		case tString:
			zfs = append(zfs, zap.String(field.key, field.value.(string))) //nolint:forcetypeassert
		case tUint:
			zfs = append(zfs, zap.Uint(field.key, field.value.(uint))) //nolint:forcetypeassert
		case tInt:
			zfs = append(zfs, zap.Int(field.key, field.value.(int))) //nolint:forcetypeassert
		case tTime:
			zfs = append(zfs, zap.Time(field.key, field.value.(time.Time))) //nolint:forcetypeassert
		case tDuration:
			zfs = append(zfs, zap.Duration(field.key, field.value.(time.Duration))) //nolint:forcetypeassert
		case tUint8:
			zfs = append(zfs, zap.Uint8(field.key, field.value.(uint8))) //nolint:forcetypeassert
		case tUint64:
			zfs = append(zfs, zap.Uint64(field.key, field.value.(uint64))) //nolint:forcetypeassert
		}
	}

	return zfs
}

// NewZapLogger creates zap logger using ectobit.com proposed configuration.
//
// format may be either console or json.
// level may be one of debug, info, warn or error. For level debug and format console output will be colored.
// callerSkip defines how many callers should be skipped when logging. Useful when zap logger is wrapped.
//
// Logger will show stack trace just in error case.
func NewZapLogger(format, level string, callerSkip int) (*zap.Logger, error) {
	encodeLevel := zapcore.LowercaseLevelEncoder
	if format == "console" && level == "debug" {
		encodeLevel = zapcore.CapitalColorLevelEncoder
	}

	zapLevel := zap.NewAtomicLevel()

	if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
		return nil, fmt.Errorf("parse loglevel: %w", err)
	}

	config := zap.Config{ //nolint:exhaustivestruct
		Level:    zapLevel,
		Encoding: format,
		EncoderConfig: zapcore.EncoderConfig{ //nolint:exhaustivestruct
			CallerKey:      "caller",
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeLevel:    encodeLevel,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			FunctionKey:    zapcore.OmitKey,
			LevelKey:       "level",
			LineEnding:     zapcore.DefaultLineEnding,
			MessageKey:     "msg",
			NameKey:        "logger",
			StacktraceKey:  "stack",
			TimeKey:        "ts",
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := config.Build(zap.AddCallerSkip(callerSkip), zap.AddStacktrace(zapcore.ErrorLevel))
	if err != nil {
		return nil, fmt.Errorf("build config: %w", err)
	}

	return logger, nil
}
