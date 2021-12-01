// Package lax contains logger abstraction and adapters for the most common loggers.
package lax

import "time"

// Logger is an abstraction of typical logger methods.
type Logger interface {
	// Debug sends a message to logger at debug level.
	Debug(string, ...Field)
	// Info sends a message to logger at info level.
	Info(string, ...Field)
	// Warn sends a message to logger at warn level.
	Warn(string, ...Field)
	// Error sends a message to logger at error level.
	Error(string, ...Field)
	// Flush flushes buffer ignoring eventual error.
	Flush()
}

type vType uint8

const (
	tAny vType = iota
	tError
	tString
	tUint
	tInt
	tTime
	tDuration
	tUint8
)

// Field is a general type to log typed values as log message attributes.
type Field struct {
	key   string
	vType vType
	value interface{}
}

// Any creates field with an attribute value of any type.
func Any(key string, value interface{}) Field {
	return Field{
		key:   key,
		vType: tAny,
		value: value,
	}
}

// Error creates field with an attribute of error type.
func Error(value error) Field {
	return Field{ //nolint:exhaustivestruct
		vType: tError,
		value: value,
	}
}

// String creates field with an attribute of string type.
func String(key string, value string) Field {
	return Field{
		key:   key,
		vType: tString,
		value: value,
	}
}

// Uint creates field with an attribute of uint type.
func Uint(key string, value uint) Field {
	return Field{
		key:   key,
		vType: tUint,
		value: value,
	}
}

// Uint8 creates field with an attribute of uint8 type.
func Uint8(key string, value uint8) Field {
	return Field{
		key:   key,
		vType: tUint8,
		value: value,
	}
}

// Int creates field with an attribute of int type.
func Int(key string, value int) Field {
	return Field{
		key:   key,
		vType: tInt,
		value: value,
	}
}

// Time creates field with an attribute of time.Time type.
func Time(key string, value time.Time) Field {
	return Field{
		key:   key,
		vType: tTime,
		value: value,
	}
}

// Duration creates field with an attribute of time.Duration type.
func Duration(key string, value time.Duration) Field {
	return Field{
		key:   key,
		vType: tDuration,
		value: value,
	}
}
