package gospec

import "testing"

type testingInterface interface {
	Helper()
	Parallel()
	Errorf(format string, args ...interface{})
	Run(name string, f func(t *testing.T)) bool
}
