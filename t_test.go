package gospec_test

import (
	"sync"
	"testing"
)

// this test is just meant to define this global variable which
// will get reused between the other Example tests which are only
// meant to be used for documentation purposes

var (
	t               *testing.T     //nolint:gochecknoglobals
	parallelTestsWg sync.WaitGroup //nolint:gochecknoglobals
)
