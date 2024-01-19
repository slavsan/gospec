package gospec

import (
	"sync"
	"testing"
)

// World is a structure which acts as a variables container
// in parallel tests. Changes to the contents of World are
// concurrently safe to make.
type World struct {
	T      *testing.T
	values map[string]any
	mu     sync.Mutex
}

func newWorld() *World {
	return &World{
		values: map[string]any{},
	}
}

// Get retrieves a value provided a name. If there is no such
// variable with such a name defined, an error will be reported.
func (w *World) Get(name string) any {
	w.mu.Lock()
	defer w.mu.Unlock()
	value, ok := w.values[name]
	if !ok {
		w.T.Errorf("world does not have value set for '%s'", name)
	}
	return value
}

// Set adds a new variable value for a given name.
// It is meant to be used only when initially initializing
// the variable.
func (w *World) Set(name string, value any) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.values[name] = value
}

// Swap updates a given variable in [World]. It errors if there
// is no variable with such a name already defined.
func (w *World) Swap(name string, f func(any) any) {
	w.mu.Lock()
	defer w.mu.Unlock()
	value, ok := w.values[name]
	if !ok {
		w.T.Errorf("can not swap value, since world does not have value set for '%s', try setting it first", name)
	}
	newValue := f(value)
	w.values[name] = newValue
}
