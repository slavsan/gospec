package gospec

import "testing"

type testingInterface interface {
	Helper()
	Parallel()
	Skip(args ...any)
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...any)
	Failed() bool
	Skipped() bool
	Run(name string, f func(t *testing.T)) bool
}
