package lax

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ Logger = (*ZapAdapter)(nil)

// ZapAdapter implements Logger interface for zap logger.
type ZapAdapter struct {
	l *zap.Logger
}

// NewZapAdapter creates Logger adapter for zap logger.
func NewZapAdapter(l *zap.Logger) *ZapAdapter {
	return &ZapAdapter{l: l}
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

func (l *ZapAdapter) toZapFields(fields []Field) []zapcore.Field {
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
		}
	}

	return zfs
}
