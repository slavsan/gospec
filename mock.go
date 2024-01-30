package gospec

import (
	"sync"
	"testing"
)

type mock struct {
	t          *testing.T
	calls      [][]any
	testTitles []string
	childMocks []*mock
}

func (m *mock) Helper() {
	// ..
}

func (m *mock) Parallel() {
	// ..
}

func (m *mock) Failed() bool {
	return len(m.calls) > 0
}

func (m *mock) Skipped() bool {
	return false
}

func (m *mock) Errorf(format string, args ...interface{}) {
	var call []any
	call = append(call, format)
	call = append(call, args...)
	m.calls = append(m.calls, call)
}

func (m *mock) Run(name string, f func(t *testing.T)) bool {
	m.testTitles = append(m.testTitles, name)
	m.t.Run(name, func(t *testing.T) {
		tm := &mock{t: t}
		m.childMocks = append(m.childMocks, tm)
		f(t)
	})
	return false
}

type assertMock struct {
	mu    sync.Mutex
	calls [][]any
}

func (m *assertMock) Assert(args ...any) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls = append(m.calls, args)
	// ..
}
