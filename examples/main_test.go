package examples_test

import (
	"log"
	"sync"
	"testing"
	"time"
)

var parallelTestsWg *sync.WaitGroup

func TestMain(m *testing.M) {
	done := make(chan struct{})
	parallelTestsWg = &sync.WaitGroup{}
	go func() {
		parallelTestsWg.Wait()
		close(done)
	}()

	code := m.Run()
	if code != 0 {
		log.Fatalf("some tests failed: %d", code)
	}

	select {
	case <-done:
		// all good
	case <-time.After(30 * time.Second):
		log.Fatalf("timed out whilst waiting for tests to finished")
	}
}
