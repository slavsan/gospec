package gospec

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"
)

// W is an alias for the [World] struct type.
type W = World

// World is a structure which acts as a variables container
// in parallel tests. Changes to the contents of World are
// concurrently safe to make.
type World struct {
	t                  *testing.T
	values             map[string]any
	mu                 sync.Mutex
	currentFeatureStep *featureStep
}

func newWorld() *World {
	return &World{
		currentFeatureStep: nil,
		values:             map[string]any{},
	}
}

// Table is used for displaying a Gherkin style table in parallel tests.
// It does the same thing as the [Table] function, but It's used only in parallel tests
// when an instance of a test-scoped [World] struct is passed to the callback.
// The Table method is supposed to be used only in [FeatureSuite] tests.
func (w *World) Table(fs *FeatureSuite, items any, columns ...string) { //nolint:gocognit,cyclop
	w.t.Helper()

	// TODO: validate table was called in valid call site

	var sb strings.Builder
	n := &node2{}

	items2 := reflect.ValueOf(items)

	if items2.Kind() != reflect.Slice {
		w.t.Errorf("expected items to be of type slice but was of type: %v", reflect.TypeOf(items))
		return
	}

	columnWidths := make(map[string]int, items2.Len())
	_ = columnWidths

	for _, x := range columns {
		columnWidths[x] = len(x)
	}

	var rows []map[string]string

	for i := 0; i < items2.Len(); i++ {
		item := items2.Index(i)
		if item.Kind() == reflect.Struct {
			row := map[string]string{}
			v := reflect.Indirect(item)
			for j := 0; j < v.NumField(); j++ {
				name := v.Type().Field(j).Name
				value := v.Field(j).Interface()
				maxWidth, ok := columnWidths[name]
				_ = maxWidth
				if !ok {
					continue
				}
				switch z := value.(type) {
				case string:
					if len(z) > maxWidth {
						columnWidths[name] = len(z)
					}
					row[name] = z
				case float64, float32:
					ff := fmt.Sprintf("%.2f", z)
					if len(ff) > maxWidth {
						columnWidths[name] = len(ff)
					}
					row[name] = ff
				case int, int8, int16, int32, int64:
					ff := fmt.Sprintf("%d", z)
					if len(ff) > maxWidth {
						columnWidths[name] = len(ff)
					}
					row[name] = ff
				}
			}
			rows = append(rows, row)
		}
	}
	sb.WriteString("\t\t\t|")
	for _, c := range columns {
		sb.WriteString(fmt.Sprintf(" %-"+strconv.Itoa(columnWidths[c])+"s ", c))
		sb.WriteString("|")
	}
	sb.WriteString("\n")

	for i, r := range rows {
		_ = r
		if i != 0 {
			sb.WriteString("\n")
		}
		sb.WriteString("\t\t\t|")
		for _, c := range columns {
			sb.WriteString(fmt.Sprintf(" %-"+strconv.Itoa(columnWidths[c])+"s ", r[c]))
			_ = c

			sb.WriteString("|")
		}
	}

	s := &featureStep{
		kind:  isTable,
		title: sb.String(),
	}

	n.step = s

	fs.mu.Lock()
	w.mu.Lock()
	defer w.mu.Unlock()
	defer fs.mu.Unlock()

	if len(w.currentFeatureStep.n.children) > 0 {
		return
	}

	w.currentFeatureStep.n.children = append(w.currentFeatureStep.n.children, n)
}

// Get retrieves a value provided a name. If there is no such
// variable with such a name defined, an error will be reported.
func (w *World) Get(name string) any {
	w.mu.Lock()
	defer w.mu.Unlock()
	value, ok := w.values[name]
	if !ok {
		w.t.Errorf("world does not have value set for '%s'", name)
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
		w.t.Errorf("can not swap value, since world does not have value set for '%s', try setting it first", name)
	}
	newValue := f(value)
	w.values[name] = newValue
}
