package gospec

import (
	"errors"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

const (
	bold    = "\033[1m"
	noBold  = "\033[0m"
	noColor = "\033[0m"
	gray    = "\033[1;30m"
	red     = "\033[0;31m"
	green   = "\033[0;32m"
	yellow  = "\033[0;33m"
	blue    = "\033[0;34m"
	purple  = "\033[0;35m"
	cyan    = "\033[0;36m"
)

// SpecOption is used for changing the behaviour of [SpecSuite] tests.
// The available options are: [Only].
type SpecOption int

const (
	undefinedSpecOption SpecOption = iota

	// Only allows for specifying which suites to run. Once there is at least
	// one `It` block with this option passed, all other ones will be skipped.
	Only

	invalidSpecOption
)

// Describe is used to define a describe block in a [SpecSuite].
type Describe func(title string, cb func())

// BeforeEach defines a block of code to be executed before all `it` ([It] or [ParallelIt]) blocks.
type BeforeEach func(cb func(t *testing.T))

// It defines a block which gets executed.
type It func(title string, cb func(t *testing.T), options ...SpecOption)

// ParallelBeforeEach is the same as [BeforeEach] but is used for parallel tests. It accepts
// additionally a *[World] instance which is used for passing state between the different steps
// of a test suite.
type ParallelBeforeEach func(cb func(t *testing.T, w *World))

// ParallelIt is the same as [It] but is used for parallel tests. It accepts
// additionally a *[World] instance which is used for passing state between the different steps
// of a test suite.
type ParallelIt func(title string, cb func(t *testing.T, w *World))

// SpecSuite is a spec suite which follows the rspec syntax, i.e.
// describe, beforeEach, it blocks, etc. It has several methods
// that can be called on it: [SpecSuite.Describe], [SpecSuite.BeforeEach],
// and [SpecSuite.It].
type SpecSuite struct {
	t           testingInterface
	parallel    bool
	done        func()
	stack       []*step
	suites      [][]*step
	indent      int
	outputs     []output1
	basePath    string
	nodes       []*node
	currNode    *node
	nodesStack  []*node
	failedCount int
	only        bool
	wg          *sync.WaitGroup
}

// WithSpecSuite defines a new [SpecSuite] instance, by passing that new instance through the callback.
//
// Example:
//
//	WithSpecSuite(t, func(s *gospec.SpecSuite) {
//		/* use the SpecSuite s here for defining one or more specs */
//	})
func WithSpecSuite(t *testing.T, callback func(s *SpecSuite)) {
	t.Helper()

	s := newSpecSuite(t)

	defer s.start2()

	callback(s)
}

// With is used for setting the options for a [SpecSuite]. It will error if called twice.
func (suite *SpecSuite) With(options ...SuiteOption) *SpecSuite {
	for _, o := range options {
		o(suite)
	}
	return suite
}

type block int

const (
	isDescribe block = iota
	isBeforeEach
	isIt
)

var (
	basePath      string //nolint:gochecknoglobals
	isBasePathSet bool   //nolint:gochecknoglobals
)

type step struct {
	t          testingInterface
	indent     int
	block      block
	title      string
	file       string
	lineNo     int
	failed     bool
	failedAt   int
	executed   bool
	only       bool
	cb         func(t *testing.T)
	parallelCb func(t *testing.T, w *World)
	timeSpent  time.Duration
	done       func()
}

type output1 struct {
	out            io.Writer
	colorful       bool
	durations      bool
	printFilenames bool
	indent         OutputOption
	indentStep     string
}

func (o *output1) renderSpec(s *SpecSuite) (int, error) {
	return o.out.Write([]byte(tree(s.nodes).String(o)))
}

func (o *output1) renderFeature(fs *FeatureSuite) (int, error) {
	return o.out.Write([]byte(tree2(fs.nodes).String(o)))
}

// NewTestSuite creates a new instance of SpecSuite.
func newSpecSuite(t *testing.T) *SpecSuite {
	t.Helper()
	suite := &SpecSuite{
		t:        t,
		indent:   0,
		basePath: getBasePath(),
	}
	return suite
}

func (suite *SpecSuite) start2() {
	if !suite.parallel {
		suite.start()
		for _, out := range suite.outputs {
			_, _ = out.renderSpec(suite)
		}
		return
	}

	suite.wg = &sync.WaitGroup{}
	suite.wg.Add(len(suite.suites))

	suite.start()

	go func() {
		suite.wg.Wait()

		for _, out := range suite.outputs {
			_, _ = out.renderSpec(suite)
		}

		if suite.done != nil {
			suite.done()
		}
	}()
}

// API returns the exposed methods on the [SpecSuite] instance. It's intended usage is as follows:
//
//	describe, beforeEach, it := s.API()
//
//	describe("my feature", func() {
//		beforeEach(func(t *testing.T) {
//			/* execute any preconditions */
//			/* or execute code under test */
//		})
//
//		it("should do this and that", func(t *testing.T) {
//			/* execute code under test and assert */
//			/* or just assert */
//		})
//	})
//
// The reason for this lies in the influence gospec has from rspec, mocha, and other
// BDD frameworks, which have a similar API and general look and feel (namely the
// lowercase `describe`, `beforeEach`, and `it` functions).
func (suite *SpecSuite) API() (
	Describe,
	BeforeEach,
	It,
) {
	if len(suite.outputs) == 0 {
		o := defaultOutput()
		o.durations = true
		suite.outputs = append(suite.outputs, o)
	}
	return suite.describe, suite.beforeEach, suite.it
}

// ParallelAPI returns the exposed functions for defining spec suites which are meant
// to run in parallel. The [ParallelBeforeEach] and [ParallelIt] functions
// accept a *[World] instance in the callback arguments. This is necessary
// so that the steps of the test (suite) can pass the test-scoped state.
func (suite *SpecSuite) ParallelAPI(done func()) (
	Describe,
	ParallelBeforeEach,
	ParallelIt,
) {
	suite.parallel = true
	suite.done = done
	if len(suite.outputs) == 0 {
		do := defaultOutput()
		do.durations = true
		suite.outputs = append(suite.outputs, do)
	}
	return suite.describe, suite.parallelBeforeEach, suite.parallelIt
}

func (suite *SpecSuite) start() { //nolint:gocognit,cyclop
	suite.t.Helper()

	for _, suite2 := range suite.suites {
		suite2 := suite2

		suite.t.Run(buildSuiteTitle(suite2), func(t *testing.T) {
			t.Helper()
			start := time.Now()
			if !suite.parallel {
				defer func() {
					lastStep := suite2[len(suite2)-1]
					if lastStep.block == isIt {
						lastStep.timeSpent = time.Since(start)
					}
				}()
			}

			world := newWorld()
			world.t = t

			if suite.parallel {
				t.Parallel()
				for _, s := range suite2 {
					if s.block == isIt {
						s.t = t
					}
					if s.block == isIt || s.block == isBeforeEach {
						if s.block == isIt {
							s.done = func() {
								suite.wg.Done()
							}
						}
						s.parallelCb(t, world)
						continue
					}
				}
				// suite.wg.Done() // TODO: perhaps just call this here ??
				return
			}

			for _, s := range suite2 {
				if s.cb == nil {
					continue
				}
				if s.block == isIt || s.block == isBeforeEach {
					if s.block == isIt {
						s.t = t
					}
					if s.block == isIt && suite.only && !s.only && s.t != nil {
						s.t.Skip()
					}
					s.cb(t)
				}
			}
		})
	}
}

func buildSuiteTitle(suite []*step) string {
	var sb strings.Builder
	for i, s := range suite {
		if s.block == isDescribe || s.block == isIt {
			if i != 0 {
				sb.WriteString("/")
			}
			sb.WriteString(strings.TrimSpace(s.title))
		}
	}
	return sb.String()
}

func (suite *SpecSuite) pushStack(s *step) {
	suite.t.Helper()
	suite.stack = append(suite.stack, s)
}

func (suite *SpecSuite) pushStack2(n *node) {
	suite.t.Helper()
	suite.nodesStack = append(suite.nodesStack, n)
}

func (suite *SpecSuite) popStack(s *step) {
	suite.t.Helper()
	if len(suite.stack) == 0 {
		suite.t.Errorf("unexpected empty stack")
		return
	}

	lastStep := suite.stack[len(suite.stack)-1]
	if lastStep != s {
		suite.t.Errorf("unexpected step")
		return
	}

	suite.stack = suite.stack[:len(suite.stack)-1]
}

func (suite *SpecSuite) popStack2(n *node) {
	suite.t.Helper()
	if len(suite.nodesStack) == 0 {
		suite.t.Errorf("unexpected empty node stack")
		return
	}

	lastNode := suite.nodesStack[len(suite.nodesStack)-1]
	if lastNode != n {
		suite.t.Errorf("unexpected node")
		return
	}

	suite.nodesStack = suite.nodesStack[:len(suite.nodesStack)-1]
}

func (suite *SpecSuite) popStackUntilStep(s *step) {
	suite.t.Helper()
	if len(suite.stack) == 0 {
		suite.t.Errorf("unexpected empty stack")
		return
	}

	index := suite.findIndexOfStep(s)
	if index < 0 {
		return
	}

	if index+1 > len(suite.stack) {
		suite.t.Errorf("out of bound index")
		return
	}

	suite.stack = suite.stack[:index+1]
}

func (suite *SpecSuite) popStackUntilStep2(n *node) {
	suite.t.Helper()
	if len(suite.nodesStack) == 0 {
		suite.t.Errorf("unexpected empty node stack")
		return
	}

	index := suite.findIndexOfNode(n)
	if index < 0 {
		return
	}

	if index+1 > len(suite.nodesStack) {
		suite.t.Errorf("out of bound index for node search")
		return
	}

	suite.nodesStack = suite.nodesStack[:index+1]
}

func (suite *SpecSuite) findIndexOfStep(s *step) int {
	suite.t.Helper()
	if len(suite.stack) == 0 {
		return -1
	}

	for i := len(suite.stack) - 1; i >= 0; i-- {
		if suite.stack[i] == s {
			return i
		}
	}

	return -1
}

func (suite *SpecSuite) findIndexOfNode(n *node) int {
	suite.t.Helper()
	if len(suite.nodesStack) == 0 {
		return -1
	}

	for i := len(suite.nodesStack) - 1; i >= 0; i-- {
		if suite.nodesStack[i] == n {
			return i
		}
	}

	return -1
}

// Describe is a function which describes a feature or contextual logical
// block, which may contain inner describe, beforeEach, or it blocks.
//
// It's important to note that Describe blocks should not mutate any state.
// The expected usage is for variables to be defined in a Describe block but
// not initialized or assigned there.
//
// When the [Parallel] is used, even declaring variables in Describe is not
// expected or advised since it would lead to undefined behaviour or race conditions.
// In those cases, just use the [World] construct which would get passed to
// all [SpecSuite.BeforeEach] and [SpecSuite.It] function calls.
func (suite *SpecSuite) describe(title string, cb func()) {
	suite.t.Helper()

	// TODO: add checks for order of blocks

	_, file, lineNo, _ := runtime.Caller(1)
	_ = file
	_ = lineNo

	n := &node{}

	if suite.isTopLevel() {
		// top level, so the node should go at the suites level
		suite.nodes = append(suite.nodes, n)
	} else {
		suite.currNode.children = append(suite.currNode.children, n)
	}

	suite.currNode = n

	suite.indent++

	s := &step{
		title:  title,
		indent: suite.indent,
		block:  isDescribe,
		file:   file,
		lineNo: lineNo,
	}

	n.step = s

	suite.pushStack(s)
	suite.pushStack2(n)

	cb()

	suite.indent--

	// TODO: check if last suite starts with the same describe we have
	// TODO: check whether the last suite contains this step
	if len(suite.stack) > 0 && !suite.lastSuiteContainsStep(s) {
		suite.copyStack()
	}

	suite.popStackUntilStep(s)
	suite.popStack(s)

	suite.popStackUntilStep2(n)
	suite.popStack2(n)

	if len(suite.nodesStack) >= 1 {
		suite.currNode = suite.nodesStack[len(suite.nodesStack)-1]
	}

	// closing top-level describe, therefore write the output
	if suite.isTopLevel() {
		suite.currNode = nil
	}
}

func (suite *SpecSuite) isTopLevel() bool {
	return suite.indent == 0
}

func (suite *SpecSuite) lastSuiteContainsStep(step *step) bool {
	if len(suite.suites) == 0 {
		return false
	}

	lastSuite := suite.suites[len(suite.suites)-1]
	for _, s := range lastSuite {
		if s == step {
			return true
		}
	}

	return false
}

// BeforeEach is a function which executes before each [SpecSuite.It] or [SpecSuite.Describe]
// block which is defined after it.
//
// It is used for assigning values to variables which are then used in the following
// blocks. If the [SpecSuite.BeforeEach] block is not followed by a [SpecSuite.It] block, it
// will not get executed.
func (suite *SpecSuite) parallelBeforeEach(cb func(*testing.T, *World)) {
	suite.t.Helper()

	s := &step{
		indent:     suite.indent,
		block:      isBeforeEach,
		parallelCb: cb,
	}

	suite.pushStack(s)
}

func (suite *SpecSuite) beforeEach(cb func(*testing.T)) {
	suite.t.Helper()

	s := &step{
		indent: suite.indent,
		block:  isBeforeEach,
		cb:     cb,
	}

	suite.pushStack(s)
}

// It defines a block which gets executed in a test suite as the last step. [SpecSuite.It] blocks
// can not be nested.
func (suite *SpecSuite) parallelIt(title string, cb func(t *testing.T, w *World)) {
	suite.t.Helper()

	_, file, lineNo, _ := runtime.Caller(1)
	_ = file
	_ = lineNo

	// TODO: check if parallel and make sure `cb` is defined with *World as the first arg

	n := &node{}

	s := &step{
		title:  title,
		indent: suite.indent,
		block:  isIt,
		file:   file,
		lineNo: lineNo,
		cb:     nil,
	}

	n.step = s

	suite.currNode.children = append(suite.currNode.children, n)

	s.parallelCb = func(t *testing.T, w *World) {
		w.t.Helper()

		if suite.parallel {
			defer s.done()
		}

		cb(t, w)

		s.executed = true

		if w.t.Failed() {
			s.failed = true
			suite.failedCount++
			s.failedAt = suite.failedCount
		}
	}

	suite.pushStack(s)

	if len(suite.stack) > 0 {
		suite.copyStack()
	}

	suite.popStack(s)
}

func (suite *SpecSuite) it(title string, cb func(t *testing.T), options ...SpecOption) {
	suite.t.Helper()

	_, file, lineNo, _ := runtime.Caller(1)
	_ = file
	_ = lineNo

	// TODO: check if parallel and make sure `cb` is defined with *World as the first arg

	n := &node{}

	s := &step{
		title:  title,
		indent: suite.indent,
		block:  isIt,
		file:   file,
		lineNo: lineNo,
		cb:     nil,
	}

	for _, o := range options {
		if o <= undefinedSpecOption || o >= invalidSpecOption {
			suite.t.Errorf("invalid spec option passed")
		}
		if o == Only {
			s.only = true
			suite.only = true
		}
	}

	n.step = s

	suite.currNode.children = append(suite.currNode.children, n)

	s.cb = func(t *testing.T) {
		t.Helper()

		if suite.parallel {
			defer s.done()
		}

		cb(t)

		s.executed = true

		if t.Failed() {
			s.failed = true
			suite.failedCount++
			s.failedAt = suite.failedCount
		}
	}

	suite.pushStack(s)

	if len(suite.stack) > 0 {
		suite.copyStack()
	}

	suite.popStack(s)
}

func (suite *SpecSuite) copyStack() {
	suite.t.Helper()
	if len(suite.stack) == 0 {
		return
	}

	suiteCopy := make([]*step, 0, len(suite.stack))
	suiteCopy = append(suiteCopy, suite.stack...)
	suite.suites = append(suite.suites, suiteCopy)
}

// OutputOption is an option for controlling the formatting of custom outputs.
type OutputOption int

const (
	undefinedOption OutputOption = iota

	// Colorful enables colorful output on Unix/Linux terminals. No support for Windows currently.
	Colorful

	// Durations is an option for displaying the durations for individual test suites.
	Durations

	// PrintFilenames is an option for showing the filename:line to the specific test suite step, the location of its definition.
	PrintFilenames

	// IndentTwoSpaces is an option for enabling two spaces as the preferred indentation step.
	IndentTwoSpaces

	// IndentFourSpaces is an option for enabling four spaces as the preferred indentation step.
	IndentFourSpaces

	// IndentOneTab is an option for enabling one tab as the preferred indentation step.
	IndentOneTab

	invalidOption
)

const (
	indentTwoSpaces  = "  "
	indentFourSpaces = "    "
	indentOneTab     = "	"
)

func indentStep(o OutputOption) string {
	switch o { //nolint:exhaustive
	case IndentTwoSpaces:
		return indentTwoSpaces
	case IndentFourSpaces:
		return indentFourSpaces
	case IndentOneTab:
		return indentOneTab
	}
	return indentTwoSpaces
}

func (o OutputOption) string() string {
	switch o { //nolint:exhaustive
	case Colorful:
		return "colors enabled option"
	case PrintFilenames:
		return "print filenames option"
	case IndentTwoSpaces:
		return "two spaces indentation"
	case IndentFourSpaces:
		return "four spaces indentation"
	case IndentOneTab:
		return "one tab indentation"
	default:
		return "invalid option"
	}
}

func (suite *SpecSuite) setOutput(w io.Writer, outputOptions ...OutputOption) {
	suite.t.Helper()
	suite.outputs = append(suite.outputs, setOutput(suite.t, w, outputOptions...))
}

func getBasePath() string {
	if isBasePathSet {
		return basePath
	}

	isBasePathSet = true

	cwd, err := os.Getwd()
	if err != nil {
		return basePath
	}

	parts := strings.Split(cwd, "/")
	for i := len(parts); i >= 0; i-- {
		pathWithoutGoMod := strings.Join(parts[:i], "/") + "/"
		goModFile := pathWithoutGoMod + "go.mod"
		_, err = os.Stat(goModFile)
		if err == nil {
			basePath = pathWithoutGoMod
			break
		}
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		basePath = pathWithoutGoMod
	}

	return basePath
}
