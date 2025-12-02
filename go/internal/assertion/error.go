package assertion

import (
	"reflect"
	"testing"
)

func Equal[T comparable](t testing.TB, expected, actual T) {
	t.Helper()

	if expected != actual {
		t.Errorf("Not equal: \nexpected: %#v\nactual  : %#v", expected, actual)
	}
}

func Len(t testing.TB, object any, length int) {
	t.Helper()

	v := reflect.ValueOf(object)
	if v.Len() != length {
		t.Errorf("\"%v\" should have %d item(s), but has %d", object, length, v.Len())
	}
}

func MustNoError(t testing.TB, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("Received unexpected error:\n%+v", err)
	}
}
