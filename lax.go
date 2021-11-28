// Package lax contains logger abstraction and adapters for the most common loggers.
package lax

// Logger is an abstraction of typical logger methods.
type Logger interface {
	Debug(string, ...Field)
	Info(string, ...Field)
	Warn(string, ...Field)
	Error(string, ...Field)
}

type vType uint8

const (
	tAny vType = iota
	tError
	tString
	tUint
	tInt
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

// Int creates field with an attribute of int type.
func Int(key string, value int) Field {
	return Field{
		key:   key,
		vType: tInt,
		value: value,
	}
}
