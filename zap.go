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
// logFormat may be either console or json.
// logLevel may be one of debug, info, warn or error. At debug level and console format output will be colored.
func NewDefaultZapAdapter(logFormat, logLevel string) (*ZapAdapter, error) {
	level := zap.NewAtomicLevel()

	encodeLevel := zapcore.LowercaseLevelEncoder
	if logFormat == "console" && logLevel == "debug" {
		encodeLevel = zapcore.CapitalColorLevelEncoder
	}

	config := zap.Config{ //nolint:exhaustivestruct
		Level:       level,
		Development: logLevel == "debug",
		Encoding:    logFormat,
		EncoderConfig: zapcore.EncoderConfig{ //nolint:exhaustivestruct
			CallerKey:      "caller",
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeLevel:    encodeLevel,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			LevelKey:       "level",
			MessageKey:     "msg",
			NameKey:        "logger",
			StacktraceKey:  "stack",
		},

		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := config.Build(zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))
	if err != nil {
		return nil, fmt.Errorf("build config: %w", err)
	}

	return NewZapAdapter(logger), nil
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
			zfs = append(zfs, zap.Error(field.value.(error)))
		case tString:
			zfs = append(zfs, zap.String(field.key, field.value.(string)))
		case tUint:
			zfs = append(zfs, zap.Uint(field.key, field.value.(uint)))
		case tInt:
			zfs = append(zfs, zap.Int(field.key, field.value.(int)))
		case tTime:
			zfs = append(zfs, zap.Time(field.key, field.value.(time.Time)))
		case tDuration:
			zfs = append(zfs, zap.Duration(field.key, field.value.(time.Duration)))
		case tUint8:
			zfs = append(zfs, zap.Uint8(field.key, field.value.(uint8)))
		case tUint64:
			zfs = append(zfs, zap.Uint64(field.key, field.value.(uint64)))
		}
	}

	return zfs
}
