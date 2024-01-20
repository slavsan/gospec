package gospec

import (
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"testing"
)

// Suite is a spec suite which follows the rspec syntax, i.e.
// describe, beforeEach, it blocks, etc. It has several methods
// that can be called on it: [Suite.Describe], [Suite.BeforeEach],
// and [Suite.It].
type Suite struct {
	t              testingInterface
	parallel       bool
	stack          []*step
	suites         [][]*step
	indent         int
	atSuiteIndex   int
	out            io.Writer
	report         strings.Builder
	basePath       string
	printFilenames bool
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
	indent int
	block  block
	title  string
	cb     any
}

// NewTestSuite creates a new instance of Suite.
func NewTestSuite(t testingInterface, options ...SuiteOption) *Suite {
	suite := &Suite{
		t:        t,
		out:      os.Stdout,
		indent:   0,
		basePath: getBasePath(),
	}
	for _, o := range options {
		o(suite)
	}
	return suite
}

// API returns the exposed methods on the [Suite] instance. It's intended usage is as follows:
//
//	describe, beforeEach, it := gospec.NewTestSuite(t).API()
//
//	describe("my feature", func() {
//		beforeEach(func() {
//			// ..
//		})
//
//		// ..
//	})
//
// The reason for this lies in the influence gospec has from rspec, mocha, and other
// BDD frameworks, which have a similar API and general look and feel (namely the
// lowercase describe, beforeEach, it functions).
//
// Alternatively you can export the Describe, BeforeEach, and It like so
//
//	var (
//		spec = gospec.NewTestSuite(t)
//		describe = spec.Describe
//		beforeEach = spec.BeforeEach
//		it = spec.It
//	)
//
// Or, you can just instantiate the suite instance and use the public methods on it directly
//
//	spec := gospec.NewTestSuite(t)
//
//	spec.Describe("my feature", func() {
//		spec.BeforeEach(func() {
//			// ..
//		})
//
//		// ..
//	})
func (suite *Suite) API() (
	func(string, any),
	func(any),
	func(string, any),
) {
	return suite.Describe, suite.BeforeEach, suite.It
}

func (suite *Suite) start() {
	for i := suite.atSuiteIndex; i < len(suite.suites); i++ {
		suite2 := suite.suites[i]
		suite.atSuiteIndex++
		suite.t.Run(buildSuiteTitle(suite2), func(t *testing.T) {
			t.Helper()
			world := newWorld()
			world.T = t
			if suite.parallel {
				t.Parallel()
				for _, s := range suite2 {
					if s.block == isIt || s.block == isBeforeEach {
						s.cb.(func(w *World))(world)
						continue
					}
					if s.cb != nil {
						s.cb.(func())()
					}
				}
				return
			}

			for _, s := range suite2 {
				if s.cb != nil {
					s.cb.(func())()
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

func (suite *Suite) pushStack(s *step) {
	suite.t.Helper()
	suite.stack = append(suite.stack, s)
}

func (suite *Suite) popStack(s *step) {
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

func (suite *Suite) popStackUntilStep(s *step) {
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

func (suite *Suite) findIndexOfStep(s *step) int {
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

func (suite *Suite) print(title string) {
	pc, file, lineNo, ok := runtime.Caller(2)
	_ = pc
	_ = file
	_ = lineNo
	_ = ok

	if !suite.printFilenames {
		suite.report.WriteString(fmt.Sprintf("%s%s\n", strings.Repeat("\t", suite.indent), title))
		return
	}

	suite.report.WriteString(fmt.Sprintf("%s%s\t%s:%d\n",
		strings.Repeat("\t", suite.indent), title,
		strings.TrimPrefix(file, suite.basePath), lineNo,
	))
}

// Describe is a function which describes a feature or contextual logical
// block, which may contain inner describe, beforeEach, or it blocks.
//
// It's important to note that Describe blocks should not mutate any state.
// The expected usage is for variables to be defined in a Describe block but
// not initialized or assigned there.
//
// When the [WithParallel] is used, even declaring variables in Describe is not
// expected or advised since it would lead to undefined behaviour or race conditions.
// In those cases, just use the [World] construct which would get passed to
// all [Suite.BeforeEach] and [Suite.It] function calls.
func (suite *Suite) Describe(title string, cb any) {
	suite.t.Helper()

	// TODO: add checks for order of blocks

	if suite.indent == 0 {
		// starting with a new top-level describe, so create a new router
		suite.report = strings.Builder{}
	}

	suite.print(title)

	suite.indent++

	s := &step{
		title:  title,
		indent: suite.indent,
		block:  isDescribe,
	}

	suite.pushStack(s)

	cb.(func())()

	suite.indent--
	suite.popStackUntilStep(s)
	suite.popStack(s)

	// closing top-level describe, therefore write the output
	if suite.indent == 0 {
		suite.start()
		suite.report.WriteString("\n")
		_, _ = suite.out.Write([]byte(suite.report.String()))
	}
}

// BeforeEach is a function which executes before each [Suite.It] or [Suite.Describe]
// block which is defined after it.
//
// It is used for assigning values to variables which are then used in the following
// blocks. If the [Suite.BeforeEach] block is not followed by a [Suite.It] block, it
// will not get executed.
func (suite *Suite) BeforeEach(cb any) {
	suite.t.Helper()

	s := &step{
		indent: suite.indent,
		block:  isBeforeEach,
		cb:     cb,
	}

	suite.pushStack(s)
}

// It defines a block which gets executed in a test suite as the last step. [Suite.It] blocks
// can not be nested.
func (suite *Suite) It(title string, cb any) {
	suite.t.Helper()

	suite.print(title)

	// TODO: check if parallel and make sure `cb` is defined with *World as the first arg

	s := &step{
		title:  title,
		indent: suite.indent,
		block:  isIt,
		cb:     cb,
	}

	suite.pushStack(s)

	if len(suite.stack) > 0 {
		suite.copyStack()
	}

	suite.popStack(s)
}

func (suite *Suite) copyStack() {
	suite.t.Helper()
	if len(suite.stack) == 0 {
		return
	}

	suiteCopy := make([]*step, 0, len(suite.stack))
	suiteCopy = append(suiteCopy, suite.stack...)
	suite.suites = append(suite.suites, suiteCopy)
}

func (suite *Suite) setOutput(w io.Writer) {
	suite.out = w
}

func (suite *Suite) setPrintFilenames() {
	suite.printFilenames = true
}

func (suite *Suite) setParallel() {
	suite.parallel = true
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
