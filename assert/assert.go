package assert

import (
	"reflect"
	"testing"
)

func True(t *testing.T, result bool) {
	if !result {
		t.Fail()
	}
}

func DeepEqual(t *testing.T, x interface{}, y interface{}) {
	True(t, reflect.DeepEqual(x, y))
}
