package assert

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func Equal(t *testing.T, expected, actual any, msg ...string) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		message := ""
		if len(msg) > 0 {
			message = fmt.Sprintf("\t message: %s", strings.Join(msg, "\n"))
		}
		t.Errorf(
			//"equality assertion failed:\n\texpected:\n\n%s (%s)\n\n\t  actual:\n\n%s (%s)\n\n%s",
			"equality assertion failed:\n\texpected: %#v (%s)\n\t  actual: %#v (%s)\n%s",
			expected, reflect.TypeOf(expected),
			actual, reflect.TypeOf(actual),
			message,
		)
	}
}

//func NotNil(t *testing.T, value any) {
//	t.Helper()
//	if isNil(value) {
//		t.Errorf("expected %v to be nil but it is not", value)
//	}
//}
//
//func Nil(t *testing.T, value any) {
//	t.Helper()
//	if !isNil(value) {
//		t.Errorf("expected %v to be not nil but it is", value)
//	}
//}

//func isNil(value any) bool {
//	if value == nil {
//		return true
//	}
//	valueOf := reflect.ValueOf(value)
//	switch valueOf.Kind() {
//	case reflect.Chan, reflect.UnsafePointer, reflect.Ptr, reflect.Slice, reflect.Map, reflect.Func, reflect.Interface:
//		if valueOf.IsNil() {
//			return true
//		}
//	}
//	return false
//}
