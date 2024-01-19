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
			"equality assertion failed:\n"+
				"\texpected: %#v (%s)\n"+
				"\t  actual: %#v (%s)\n"+
				"%s",
			expected, reflect.TypeOf(expected),
			actual, reflect.TypeOf(actual),
			message,
		)
	}
}
