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
		"int":          {5, Field{"test", 0, 5}},
		"string":       {"lax", Field{"test", 0, "lax"}},
		"int slice":    {[]int{5}, Field{"test", 0, []int{5}}},
		"string slice": {[]string{"lax"}, Field{"test", 0, []string{"lax"}}},
		"int map":      {map[string]int{"lax": 5}, Field{"test", 0, map[string]int{"lax": 5}}},
		"string map":   {map[string]string{"lax": "lax"}, Field{"test", 0, map[string]string{"lax": "lax"}}},
		"struct":       {struct{ Test int }{5}, Field{"test", 0, struct{ Test int }{5}}},
		"struct ptr":   {&struct{ Test int }{5}, Field{"test", 0, &struct{ Test int }{5}}},
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
		"new":    {errors.New("test"), Field{"", 1, "test"}}, //nolint:goerr113
		"custom": {&customError{}, Field{"", 1, "test"}},
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

type customError struct{}

func (err *customError) Error() string {
	return "test"
}
