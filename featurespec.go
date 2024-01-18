package gospec

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
)

type featureStepKind int

const (
	isUndefined featureStepKind = iota
	isFeature
	isBackground
	isScenario
	isGiven
	isWhen
	isThen
)

type featureStep struct {
	kind  featureStepKind
	title string
	cb    any
}

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

func (w *World) Get(name string) any {
	w.mu.Lock()
	defer w.mu.Unlock()
	value, ok := w.values[name]
	if !ok {
		w.T.Errorf("world does not have value set for '%s'", name)
	}
	return value
}

func (w *World) Set(name string, value any) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.values[name] = value
}

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

type FeatureSuite struct {
	t               testingInterface
	parallel        bool
	stack           []*featureStep
	backgroundStack []*featureStep
	suites          [][]*featureStep
	inBackground    bool
	atSuiteIndex    int
	out             io.Writer
	report          strings.Builder
	basePath        string
	printFilenames  bool
}

func NewFeatureSuite(t testingInterface, options ...SuiteOption) *FeatureSuite {
	fs := &FeatureSuite{
		t:        t,
		out:      os.Stdout,
		basePath: getBasePath(),
	}
	for _, o := range options {
		o(fs)
	}
	return fs
}

func (fs *FeatureSuite) API() (
	func(string, any),
	func(any),
	func(string, any),
	func(string, any),
	func(string, any),
	func(string, any),
	func(columns []string, items any),
) {
	return fs.Feature, fs.Background, fs.Scenario, fs.Given,
		fs.When, fs.Then, fs.Table
}

func (fs *FeatureSuite) prevKind() featureStepKind {
	fs.t.Helper()
	if len(fs.stack) == 0 {
		return isUndefined
	}
	return fs.stack[len(fs.stack)-1].kind
}

func (fs *FeatureSuite) Feature(title string, cb any) {
	fs.report = strings.Builder{}
	fs.t.Helper()
	if fs.prevKind() != isUndefined {
		fs.t.Errorf("invalid position for `Feature` function, it must be at top level")
		return
	}

	fs.print(fmt.Sprintf("Feature: %s", title))

	s := &featureStep{
		kind:  isFeature,
		title: title,
	}
	fs.pushStack(s)

	// TODO: validate cb is of correct type
	cb.(func())()

	// check if last there is a new suite added, if not, copy the stack here...
	fs.popBackgroundFromStackIfExists()
	fs.popStack(s)
	fs.backgroundStack = []*featureStep{}

	if len(fs.stack) > 0 {
		fs.t.Errorf("expected stack to be empty but it has %d steps", len(fs.stack))
		return
	}

	fs.start()

	fs.report.WriteString("\n")
	_, _ = fs.out.Write([]byte(fs.report.String()))
}

func (fs *FeatureSuite) pushStack(s *featureStep) {
	fs.t.Helper()
	fs.stack = append(fs.stack, s)
}

func (fs *FeatureSuite) pushToBackgroundStack(s *featureStep) {
	fs.t.Helper()
	fs.backgroundStack = append(fs.backgroundStack, s)
}

func (fs *FeatureSuite) popBackgroundFromStackIfExists() {
	fs.t.Helper()
	if len(fs.stack) == 0 {
		return
	}

	lastStep := fs.stack[len(fs.stack)-1]
	if lastStep.kind == isBackground {
		fs.stack = fs.stack[:len(fs.stack)-1]
	}
}

func (fs *FeatureSuite) popStack(s *featureStep) {
	fs.t.Helper()
	if len(fs.stack) == 0 {
		fs.t.Errorf("unexpected empty stack")
		return
	}

	lastStep := fs.stack[len(fs.stack)-1]
	if lastStep != s {
		fs.t.Errorf("unexpected step")
		return
	}

	fs.stack = fs.stack[:len(fs.stack)-1]
}

func (fs *FeatureSuite) popStackUntilStep(s *featureStep) {
	fs.t.Helper()
	if len(fs.stack) == 0 {
		fs.t.Errorf("unexpected empty stack")
		return
	}

	index := fs.findIndexOfStep(s)
	if index < 0 {
		return
	}

	if index+1 > len(fs.stack) {
		fs.t.Errorf("out of bound index")
		return
	}

	fs.stack = fs.stack[:index+1]
}

func (fs *FeatureSuite) findIndexOfStep(s *featureStep) int {
	fs.t.Helper()
	if len(fs.stack) == 0 {
		return -1
	}

	for i := len(fs.stack) - 1; i >= 0; i-- {
		if fs.stack[i] == s {
			return i
		}
	}

	return -1
}

func (fs *FeatureSuite) Background(cb any) {
	fs.t.Helper()
	if fs.prevKind() != isFeature {
		fs.t.Errorf("invalid position for `Background` function, it must be inside a `Feature` call")
		return
	}

	fs.print("\n\tBackground:")

	s := &featureStep{
		kind: isBackground,
	}

	fs.inBackground = true
	fs.pushToBackgroundStack(s)

	cb.(func())()

	fs.inBackground = false
}

func (fs *FeatureSuite) Scenario(title string, cb any) {
	fs.t.Helper()
	if fs.prevKind() != isFeature && fs.prevKind() != isBackground {
		fs.t.Errorf("invalid position for `Scenario` function, it must be inside a `Feature` call")
		return
	}

	fs.print(fmt.Sprintf("\n\tScenario: %s", title))

	s := &featureStep{
		kind:  isScenario,
		title: title,
	}
	fs.pushStack(s)

	cb.(func())()

	if len(fs.stack) > 0 {
		fs.copyStack()
		fs.popStackUntilStep(s)
	}

	fs.popStack(s)
}

func (fs *FeatureSuite) Given(title string, cb any) {
	fs.t.Helper()

	fs.print(fmt.Sprintf("\t\tGiven: %s", title))

	s := &featureStep{
		kind:  isGiven,
		title: title,
		cb:    cb,
	}
	if fs.inBackground {
		fs.pushToBackgroundStack(s)
	} else {
		fs.pushStack(s)
	}
}

func (fs *FeatureSuite) When(title string, cb any) {
	fs.t.Helper()

	fs.print(fmt.Sprintf("\t\tWhen: %s", title))

	s := &featureStep{
		kind:  isWhen,
		title: title,
		cb:    cb,
	}

	if fs.inBackground {
		fs.pushToBackgroundStack(s)
	} else {
		fs.pushStack(s)
	}
}

func (fs *FeatureSuite) Then(title string, cb any) {
	fs.t.Helper()

	fs.print(fmt.Sprintf("\t\tThen: %s", title))

	s := &featureStep{
		kind:  isThen,
		title: title,
		cb:    cb,
	}
	if fs.inBackground {
		fs.pushToBackgroundStack(s)
	} else {
		fs.pushStack(s)
	}
}

func (fs *FeatureSuite) copyStack() {
	fs.t.Helper()
	if len(fs.stack) <= 0 {
		return
	}

	var suite []*featureStep
	for _, s := range fs.stack[:1] {
		suite = append(suite, s)
	}
	for _, s := range fs.backgroundStack {
		suite = append(suite, s)
	}
	for _, s := range fs.stack[1:] {
		suite = append(suite, s)
	}
	fs.suites = append(fs.suites, suite)
}

func (fs *FeatureSuite) Table(columns []string, items interface{}) {
	fs.t.Helper()

	// TODO: validate table was called in valid call site

	items2 := reflect.ValueOf(items)

	if items2.Kind() != reflect.Slice {
		panic("EXPECTED SLICE...\n")
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
	fs.report.WriteString("\t\t\t|")
	for _, c := range columns {
		fs.report.WriteString(fmt.Sprintf(" %-"+strconv.Itoa(columnWidths[c])+"s ", c))
		fs.report.WriteString("|")
	}
	fs.report.WriteString("\n")

	for _, r := range rows {
		_ = r
		fs.report.WriteString("\t\t\t|")
		for _, c := range columns {
			fs.report.WriteString(fmt.Sprintf(" %-"+strconv.Itoa(columnWidths[c])+"s ", r[c]))
			_ = c

			fs.report.WriteString("|")
		}
		fs.report.WriteString("\n")
	}
}

func buildSuiteTitleForFeature(suite []*featureStep) string {
	var sb strings.Builder
	for i, s := range suite {
		if s.kind == isFeature || s.kind == isScenario {
			if i != 0 {
				sb.WriteString("/")
			}
			sb.WriteString(strings.TrimSpace(s.title))
		}
	}
	return sb.String()
}

func (fs *FeatureSuite) start() {
	for i := fs.atSuiteIndex; i < len(fs.suites); i++ {
		suite := fs.suites[i]
		fs.atSuiteIndex++
		fs.t.Run(buildSuiteTitleForFeature(suite), func(t *testing.T) {
			world := newWorld()
			world.T = t
			if fs.parallel {
				t.Parallel()
				for _, s := range suite {
					if s.kind == isGiven || s.kind == isWhen || s.kind == isThen {
						s.cb.(func(w *World))(world)
						continue
					}
					if s.cb != nil {
						s.cb.(func())()
					}
				}
				return
			}

			for _, s := range suite {
				if s.cb != nil {
					s.cb.(func())()
				}
			}
		})
	}
}

func (fs *FeatureSuite) print(title string) {
	pc, file, lineNo, ok := runtime.Caller(2)
	_ = pc
	_ = file
	_ = lineNo
	_ = ok

	if !fs.printFilenames {
		fs.report.WriteString(title + "\n")
		return
	}

	fs.report.WriteString(fmt.Sprintf("%s\t%s:%d\n",
		title, strings.TrimPrefix(file, fs.basePath), lineNo,
	))
}

func (fs *FeatureSuite) setOutput(w io.Writer) {
	fs.out = w
}

func (fs *FeatureSuite) setPrintFilenames() {
	fs.printFilenames = true
}

func (fs *FeatureSuite) setParallel() {
	fs.parallel = true
}
