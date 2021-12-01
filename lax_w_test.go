package lax //nolint:testpackage

import (
	"errors"
	"reflect"
	"testing"
)

func TestAny(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		in   interface{}
		want Field
	}{
		"int":          {5, Field{"test", tAny, 5}},
		"string":       {"lax", Field{"test", tAny, "lax"}},
		"int slice":    {[]int{5}, Field{"test", tAny, []int{5}}},
		"string slice": {[]string{"lax"}, Field{"test", tAny, []string{"lax"}}},
		"int map":      {map[string]int{"lax": 5}, Field{"test", tAny, map[string]int{"lax": 5}}},
		"string map":   {map[string]string{"lax": "lax"}, Field{"test", tAny, map[string]string{"lax": "lax"}}},
		"struct":       {struct{ Test int }{5}, Field{"test", tAny, struct{ Test int }{5}}},
		"struct ptr":   {&struct{ Test int }{5}, Field{"test", tAny, &struct{ Test int }{5}}},
	}

	for n, test := range tests { //nolint:paralleltest
		test := test

		t.Run(n, func(t *testing.T) {
			t.Parallel()

			got := Any("test", test.in)
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("want %#vv, got %#v", test.want, got)
			}
		})
	}
}

func TestError(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		in   error
		want Field
	}{
		"new":    {errors.New("test"), Field{"", tError, "test"}}, //nolint:goerr113
		"custom": {&customError{}, Field{"", tError, "test"}},
	}

	for n, test := range tests { //nolint:paralleltest
		test := test

		t.Run(n, func(t *testing.T) {
			t.Parallel()

			got := Error(test.in)
			if got.key != test.want.key {
				t.Errorf("want %q, got %q", test.want.key, got.key)
			}

			if got.vType != test.want.vType {
				t.Errorf("want %d, got %d", test.want.vType, got.vType)
			}

			if got.value.(error).Error() != test.want.value {
				t.Errorf("want %q, got %q", test.want.value, got.value.(error).Error())
			}
		})
	}
}

func TestString(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		in   string
		want Field
	}{
		"empty":  {"", Field{"test", tString, ""}},
		"string": {"lax", Field{"test", tString, "lax"}},
	}

	for n, test := range tests { //nolint:paralleltest
		test := test

		t.Run(n, func(t *testing.T) {
			t.Parallel()

			got := String("test", test.in)
			if got.key != test.want.key {
				t.Errorf("want %q, got %q", test.want.key, got.key)
			}

			if got.vType != test.want.vType {
				t.Errorf("want %d, got %d", test.want.vType, got.vType)
			}

			if got.value.(string) != test.want.value {
				t.Errorf("want %q, got %q", test.want.value, got.value.(string))
			}
		})
	}
}

func TestUint(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		in   uint
		want Field
	}{
		"zero": {0, Field{"test", tUint, uint(0)}},
		"some": {5, Field{"test", tUint, uint(5)}},
	}

	for n, test := range tests { //nolint:paralleltest
		test := test

		t.Run(n, func(t *testing.T) {
			t.Parallel()

			got := Uint("test", test.in)
			if got.key != test.want.key {
				t.Errorf("want %q, got %q", test.want.key, got.key)
			}

			if got.vType != test.want.vType {
				t.Errorf("want %d, got %d", test.want.vType, got.vType)
			}

			if got.value != test.want.value {
				t.Errorf("want %d, got %d", test.want.value, got.value)
			}
		})
	}
}

func TestUint8(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		in   uint8
		want Field
	}{
		"zero": {0, Field{"test", tUint, uint8(0)}},
		"some": {5, Field{"test", tUint, uint8(5)}},
	}

	for n, test := range tests { //nolint:paralleltest
		test := test

		t.Run(n, func(t *testing.T) {
			t.Parallel()

			got := Uint8("test", test.in)
			if got.key != test.want.key {
				t.Errorf("want %q, got %q", test.want.key, got.key)
			}

			if got.vType != test.want.vType {
				t.Errorf("want %d, got %d", test.want.vType, got.vType)
			}

			if got.value != test.want.value {
				t.Errorf("want %d, got %d", test.want.value, got.value)
			}
		})
	}
}

func TestInt(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		in   int
		want Field
	}{
		"negative": {-5, Field{"test", tInt, -5}},
		"zero":     {0, Field{"test", tInt, 0}},
		"positive": {5, Field{"test", tInt, 5}},
	}

	for n, test := range tests { //nolint:paralleltest
		test := test

		t.Run(n, func(t *testing.T) {
			t.Parallel()

			got := Int("test", test.in)
			if got.key != test.want.key {
				t.Errorf("want %q, got %q", test.want.key, got.key)
			}

			if got.vType != test.want.vType {
				t.Errorf("want %d, got %d", test.want.vType, got.vType)
			}

			if got.value != test.want.value {
				t.Errorf("want %d, got %d", test.want.value, got.value)
			}
		})
	}
}

type customError struct{}

func (err *customError) Error() string {
	return "test"
}
