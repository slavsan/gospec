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
	isTable
)

type Feature func(title string, cb any)
type Background func(any)
type Scenario func(title string, cb any)
type Given func(title string, cb any)
type When func(title string, cb any)
type Then func(title string, cb any)
type Table func(items any, columns ...string)

type featureStep struct {
	t        *testing.T
	kind     featureStepKind
	title    string
	printed  bool
	file     string
	lineNo   int
	failed   bool
	failedAt int
	executed bool
	cb       any
	done     func()
	n        *node2
}

// FeatureSuite is a test suite which is inspired by the Cucumber/Gherkin
// style of writing tests, in terms of defining: features, scenarios and
// given/when/then steps.
//
// Instead of using the Gherkin syntax though, the FeatureSuite exposes
// an API (methods) which resemble this way of structuring tests for
// defining the behaviour of production code.
//
// Those methods are: [FeatureSuite.Feature], [FeatureSuite.Background],
// [FeatureSuite.Scenario], [FeatureSuite.Given], [FeatureSuite.When],
// [FeatureSuite.Then].
type FeatureSuite struct {
	t               testingInterface
	parallel        bool
	done            func()
	stack           []*featureStep
	backgroundStack []*featureStep
	suites          [][]*featureStep
	inBackground    bool
	atSuiteIndex    int
	out             io.Writer
	basePath        string
	printFilenames  bool
	nodes           []*node2
	currNode        *node2
	nodesStack      []*node2
	wg              *sync.WaitGroup
	invalid         bool
	failedCount     int
	mu              sync.Mutex
	currentStep     *featureStep
}

// NewFeatureSuite returns a new [FeatureSuite] instance.
func newFeatureSuite(t *testing.T) *FeatureSuite {
	t.Helper()
	fs := &FeatureSuite{
		t:        t,
		out:      os.Stdout,
		basePath: getBasePath(),
	}
	return fs
}

// API returns the exposed methods on the [FeatureSuite] instance. It's intended usage is as follows:
//
//	feature, background, scenario, given, when, then, _ :=
//		gospec.NewFeatureSuite(t).API()
//
//	feature("my feature", func() {
//		var cart []string
//
//		background(func() {
//			given("some precondition", func() {
//				// ...
//			})
//		})
//
//		scenario("when ..", func() {
//			given("another nested precondition", func() {
//				// ...
//			})
//			when("something gets executed", func() {
//				// ...
//			})
//			then("something should happen", func() {
//				// ...
//			})
//		})
//	})
//
// ..
func (fs *FeatureSuite) API() (
	Feature,
	Background,
	Scenario,
	Given,
	When,
	Then,
	Table,
) {
	return fs.feature, fs.background, fs.scenario, fs.given,
		fs.when, fs.then, fs.table
}

func (fs *FeatureSuite) prevKind() featureStepKind {
	fs.t.Helper()
	if len(fs.stack) == 0 {
		return isUndefined
	}
	return fs.stack[len(fs.stack)-1].kind
}

func WithFeatureSuite(t *testing.T, f func(fs *FeatureSuite)) {
	fs := newFeatureSuite(t)

	defer fs.start()

	f(fs)
}

// Feature defines a feature block, this is the top-level block and should
// define a separate piece of functionality.
func (fs *FeatureSuite) feature(title string, cb any) {
	fs.t.Helper()
	if fs.prevKind() != isUndefined {
		fs.invalid = true
		fs.t.Errorf("invalid position for `Feature` function, it must be at top level")
		return
	}

	_, file, lineNo, _ := runtime.Caller(1)
	_ = file
	_ = lineNo

	n := &node2{
		// ..
	}

	fs.nodes = append(fs.nodes, n)

	fs.currNode = n

	s := &featureStep{
		kind:   isFeature,
		title:  title,
		file:   file,
		lineNo: lineNo,
	}

	n.step = s

	fs.pushStack(s)
	fs.pushStack2(n)

	// TODO: validate cb is of correct type
	cb.(func())()

	// check if last there is a new suite added, if not, copy the stack here...
	fs.popBackgroundFromStackIfExists()
	fs.popStack(s)
	fs.backgroundStack = []*featureStep{}

	fs.popStackUntilStep2(n)
	fs.popStack2(n)

	fs.currNode = nil

	if len(fs.stack) > 0 {
		fs.t.Errorf("expected stack to be empty but it has %d steps", len(fs.stack))
		return
	}
}

func (fs *FeatureSuite) pushStack(s *featureStep) {
	fs.t.Helper()
	fs.stack = append(fs.stack, s)
}

func (fs *FeatureSuite) pushStack2(n *node2) {
	fs.t.Helper()
	fs.nodesStack = append(fs.nodesStack, n)
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

func (fs *FeatureSuite) popStack2(n *node2) {
	fs.t.Helper()
	if len(fs.nodesStack) == 0 {
		fs.t.Errorf("unexpected empty node stack")
		return
	}

	lastStep := fs.nodesStack[len(fs.nodesStack)-1]
	if lastStep != n {
		fs.t.Errorf("unexpected node")
		return
	}

	fs.nodesStack = fs.nodesStack[:len(fs.nodesStack)-1]
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

// Background defines a block which would get executed before each [FeatureSuite.Scenario].
// It can contain multiple [FeatureSuite.Given] steps.
func (fs *FeatureSuite) background(cb any) {
	fs.t.Helper()
	if fs.prevKind() != isFeature {
		fs.t.Errorf("invalid position for `Background` function, it must be inside a `Feature` call")
		return
	}

	_, file, lineNo, _ := runtime.Caller(1)
	_ = file
	_ = lineNo

	n := &node2{
		// ..
	}

	s := &featureStep{
		kind:   isBackground,
		file:   file,
		lineNo: lineNo,
	}

	n.step = s

	fs.currNode.children = append(fs.currNode.children, n)
	fs.currNode = n

	fs.inBackground = true
	fs.pushToBackgroundStack(s)

	fs.pushStack2(n)

	cb.(func())()

	fs.popStackUntilStep2(n)
	fs.popStack2(n)

	fs.inBackground = false

	fs.currNode = fs.nodesStack[len(fs.nodesStack)-1]
}

func (fs *FeatureSuite) popStackUntilStep2(n *node2) {
	fs.t.Helper()
	if len(fs.nodesStack) == 0 {
		fs.t.Errorf("unexpected empty node stack")
		return
	}

	index := fs.findIndexOfNode(n)
	if index < 0 {
		return
	}

	if index+1 > len(fs.nodesStack) {
		fs.t.Errorf("out of bound index for node search")
		return
	}

	fs.nodesStack = fs.nodesStack[:index+1]
}

func (fs *FeatureSuite) findIndexOfNode(n *node2) int {
	fs.t.Helper()
	if len(fs.nodesStack) == 0 {
		return -1
	}

	for i := len(fs.nodesStack) - 1; i >= 0; i-- {
		if fs.nodesStack[i] == n {
			return i
		}
	}

	return -1
}

// Scenario defines a scenario block. It should test a particular feature in a particular
// scenario, provided a set of given/when/then steps.
func (fs *FeatureSuite) scenario(title string, cb any) {
	fs.t.Helper()
	if fs.prevKind() != isFeature && fs.prevKind() != isBackground {
		fs.invalid = true
		fs.t.Errorf("invalid position for `Scenario` function, it must be inside a `Feature` call")
		return
	}

	_, file, lineNo, _ := runtime.Caller(1)
	_ = file
	_ = lineNo

	n := &node2{
		// ..
	}

	s := &featureStep{
		kind:   isScenario,
		title:  title,
		lineNo: lineNo,
		file:   file,
	}
	fs.pushStack(s)

	n.step = s

	fs.currNode.children = append(fs.currNode.children, n)
	fs.currNode = n

	fs.pushStack2(n)

	cb.(func())()

	fs.popStack2(n)
	fs.currNode = fs.nodesStack[len(fs.nodesStack)-1]

	if len(fs.stack) > 0 {
		fs.copyStack()
		fs.popStackUntilStep(s)
	}

	fs.popStack(s)
}

// Given defines a block which is meant to build the prerequisites for a particular
// test. It's usual to have any test setup logic defined in a [FeatureSuite.Given]
// block.
func (fs *FeatureSuite) given(title string, cb any) {
	fs.t.Helper()

	_, file, lineNo, _ := runtime.Caller(1)
	_ = file
	_ = lineNo

	n := &node2{
		// ..
	}

	s := &featureStep{
		kind:   isGiven,
		title:  title,
		lineNo: lineNo,
		file:   file,
		//cb:     cb,
	}

	if fs.parallel {
		s.cb = func(w *World) {
			w.T.Helper()

			//defer s.done()

			cb.(func(*World))(w)
			//s.executed = true

			if w.T.Failed() {
				//s.failed = true
				//fs.failedCount++
				//s.failedAt = fs.failedCount
			}
		}
	} else {
		s.cb = func() {
			fs.t.Helper()
			cb.(func())()
			s.executed = true

			if fs.t.Failed() {
				s.failed = true
				fs.failedCount++
				s.failedAt = fs.failedCount
			}
		}
	}

	n.step = s
	fs.currNode.children = append(fs.currNode.children, n)

	s.n = n

	if fs.inBackground {
		fs.pushToBackgroundStack(s)
	} else {
		fs.pushStack(s)
	}
}

// When defines a block which should exercise the actual test.
func (fs *FeatureSuite) when(title string, cb any) {
	fs.t.Helper()

	_, file, lineNo, _ := runtime.Caller(1)
	_ = file
	_ = lineNo

	n := &node2{
		// ..
	}

	s := &featureStep{
		kind:   isWhen,
		title:  title,
		lineNo: lineNo,
		file:   file,
		cb:     cb,
	}

	n.step = s
	fs.currNode.children = append(fs.currNode.children, n)

	s.n = n

	if fs.inBackground {
		fs.pushToBackgroundStack(s)
	} else {
		fs.pushStack(s)
	}
}

// Then defines a block which should hold a set of assertions.
func (fs *FeatureSuite) then(title string, cb any) {
	fs.t.Helper()

	_, file, lineNo, _ := runtime.Caller(1)
	_ = file
	_ = lineNo

	n := &node2{
		// ..
	}

	s := &featureStep{
		kind:   isThen,
		title:  title,
		lineNo: lineNo,
		file:   file,
		cb:     cb,
	}

	n.step = s
	fs.currNode.children = append(fs.currNode.children, n)

	s.n = n

	if fs.inBackground {
		fs.pushToBackgroundStack(s)
	} else {
		fs.pushStack(s)
	}
}

func (fs *FeatureSuite) copyStack() {
	fs.t.Helper()
	if len(fs.stack) == 0 {
		return
	}

	suite := make([]*featureStep, 0, len(fs.stack)+len(fs.backgroundStack))
	suite = append(suite, fs.stack[:1]...)
	suite = append(suite, fs.backgroundStack...)
	suite = append(suite, fs.stack[1:]...)
	fs.suites = append(fs.suites, suite)
}

// Table is a utility function to only visualize test data in a table.
func (fs *FeatureSuite) table(items any, columns ...string) { //nolint:gocognit,cyclop
	fs.t.Helper()

	// TODO: detect when using in "parallel" context, and if yes, error. Should use the `World.Table` method in such cases.

	// TODO: validate table was called in valid call site

	var sb strings.Builder
	n := &node2{
		// ..
	}

	items2 := reflect.ValueOf(items)

	if items2.Kind() != reflect.Slice {
		fs.t.Errorf("expected items to be of type slice but was of type: %v", reflect.TypeOf(items))
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

	fs.currentStep.n.children = append(fs.currentStep.n.children, n)
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

func (fs *FeatureSuite) With(options ...SuiteOption) *FeatureSuite {
	for _, o := range options {
		o(fs)
	}
	return fs
}

func (fs *FeatureSuite) start() {
	fs.wg = &sync.WaitGroup{}
	fs.wg.Add(len(fs.suites))
	for i := fs.atSuiteIndex; i < len(fs.suites); i++ {
		suite := fs.suites[i]
		fs.atSuiteIndex++
		fs.t.Run(buildSuiteTitleForFeature(suite), func(t *testing.T) {
			t.Helper()

			if fs.t.Failed() {
				if fs.parallel {
					fs.wg.Done()
				}
				//t.Skip()
			}

			if fs.parallel {
				world := newWorld()
				world.T = t

				t.Parallel()
				for _, s := range suite {
					if s.kind == isGiven || s.kind == isWhen || s.kind == isThen {
						//s.done = func() {
						//	fs.wg.Done()
						//}

						world.currentFeatureStep = s

						s.cb.(func(w *World))(world)

						world.currentFeatureStep = nil

						continue
					}
					if s.cb != nil {
						s.cb.(func())()
					}
				}
				fs.wg.Done()
				return
			}

			for _, s := range suite {
				if s.cb != nil {
					if s.kind == isGiven || s.kind == isWhen || s.kind == isThen {
						s.t = t
					}

					fs.currentStep = s

					s.cb.(func())()

					fs.currentStep = nil
				}
			}
		})
	}

	if fs.invalid {
		return
	}

	if !fs.parallel {
		_, _ = fs.out.Write([]byte(tree2(fs.nodes).String(fs)))
		return
	}

	go func() {
		fs.wg.Wait()

		_, _ = fs.out.Write([]byte(tree2(fs.nodes).String(fs)))

		if fs.done != nil {
			fs.done()
		}
	}()
}

func (fs *FeatureSuite) setOutput(w io.Writer) {
	fs.out = w
}

func (fs *FeatureSuite) setPrintFilenames() {
	fs.printFilenames = true
}

func (fs *FeatureSuite) setParallel(done func()) {
	fs.parallel = true
	fs.done = done
}
