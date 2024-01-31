package gospec

import (
	"errors"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"testing"
)

type W = World

type Describe func(title string, cb func())
type BeforeEach func(cb func(t *testing.T, w *World))
type It func(title string, cb func(t *testing.T, w *World))

// SpecSuite is a spec suite which follows the rspec syntax, i.e.
// describe, beforeEach, it blocks, etc. It has several methods
// that can be called on it: [SpecSuite.Describe], [SpecSuite.BeforeEach],
// and [SpecSuite.It].
type SpecSuite struct {
	t              testingInterface
	parallel       bool
	done           func()
	stack          []*step
	suites         [][]*step
	indent         int
	atSuiteIndex   int
	out            io.Writer
	basePath       string
	printFilenames bool
	nodes          []*node
	currNode       *node
	nodesStack     []*node
	failedCount    int
	calledDone     bool
	wg             *sync.WaitGroup
}

func WithSpecSuite(t *testing.T, f func(s *SpecSuite)) {
	s := newTestSuite(t)

	defer s.start2()

	f(s)
}

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
	//t        *testing.T
	indent   int
	block    block
	title    string
	printed  bool
	file     string
	lineNo   int
	failed   bool
	failedAt int
	executed bool
	cb       func(t *testing.T, w *World)
	cb2      func()
	done     func()
}

// NewTestSuite creates a new instance of SpecSuite.
func newTestSuite(t *testing.T) *SpecSuite {
	t.Helper()
	suite := &SpecSuite{
		t:        t,
		out:      os.Stdout,
		indent:   0,
		basePath: getBasePath(),
	}
	return suite
}

func (suite *SpecSuite) start2(onDone ...func()) {
	if suite.parallel && len(onDone) > 1 {
		suite.t.Errorf("invalid number of callbacks passed to start, expected 0 or 1, got %d", len(onDone))
		return
	}

	// TODO: make sure start is not called twice
	// TODO: detect start not called

	if suite.calledDone {
		panic("already calle done")
	}

	suite.calledDone = true

	suite.wg = &sync.WaitGroup{}
	suite.wg.Add(len(suite.suites[suite.atSuiteIndex:]))

	// TODO: add some default timeout perhaps and an option to override it ???
	// TODO: start and wait for all of the tasks to finish

	suite.start()

	if !suite.parallel {
		_, _ = suite.out.Write([]byte(tree(suite.nodes).String(suite)))
		return
	}

	go func() {
		suite.wg.Wait()

		_, _ = suite.out.Write([]byte(tree(suite.nodes).String(suite)))

		if len(onDone) > 0 {
			//onDone[0]()
		}
		if suite.done != nil {
			suite.done()
		}
	}()
}

// API returns the exposed methods on the [SpecSuite] instance. It's intended usage is as follows:
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
func (suite *SpecSuite) API() (
	Describe,
	BeforeEach,
	It,
) {
	return suite.describe, suite.beforeEach, suite.it
}

//func endsInItBlock(suite []*step) bool {
//	if len(suite) == 0 {
//		return false
//	}
//
//	lastStep := suite[len(suite)-1]
//
//	return lastStep.block == isIt
//}

func (suite *SpecSuite) start() {
	for i := suite.atSuiteIndex; i < len(suite.suites); i++ {
		suite2 := suite.suites[i]
		suite.atSuiteIndex++

		suite.t.Run(buildSuiteTitle(suite2), func(t *testing.T) {
			// TODO: check if the last step is an `it` block, and if not, skip this test

			if suite.t.Failed() {
				if suite.parallel {
					suite.wg.Done()
					t.Skip()
				}
			}

			world := newWorld()
			world.t = t

			if suite.parallel {
				t.Parallel()
				for _, s := range suite2 {
					if s.block == isIt || s.block == isBeforeEach {
						if s.block == isIt {
							s.done = func() {
								suite.wg.Done()
							}
						}
						s.cb(t, world)
						continue
					}
					if s.cb2 != nil {
						s.cb2()
					}
				}
				// suite.wg.Done() // TODO: perhaps just call this here ??
				return
			}

			for _, s := range suite2 {
				if s.cb != nil {
					if s.block == isIt || s.block == isBeforeEach {
						//s.t = t
						s.cb(t, world)
						continue
					}
				}
				if s.cb2 != nil {
					s.cb2()
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

	n := &node{
		// ..
	}

	if suite.isTopLevel() {
		// starting with a new top-level describe, so create a new router
		//suite.report = strings.Builder{}

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
func (suite *SpecSuite) beforeEach(cb func(*testing.T, *World)) {
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
func (suite *SpecSuite) it(title string, cb func(t *testing.T, w *World)) {
	suite.t.Helper()

	_, file, lineNo, _ := runtime.Caller(1)
	_ = file
	_ = lineNo

	// TODO: check if parallel and make sure `cb` is defined with *World as the first arg

	n := &node{
		// ..
	}

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

	s.cb = func(t *testing.T, w *World) {
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

	//if suite.parallel {
	//	s.cb = func(w *World) {
	//		w.T.Helper()
	//
	//		defer s.done()
	//
	//		cb(w)
	//		s.executed = true
	//
	//		if w.T.Failed() {
	//			s.failed = true
	//			suite.failedCount++
	//			s.failedAt = suite.failedCount
	//		}
	//	}
	//} else {
	//	s.cb = func(w *World) {
	//		suite.t.Helper()
	//		cb(w)
	//		s.executed = true
	//
	//		if suite.t.Failed() {
	//			s.failed = true
	//			suite.failedCount++
	//			s.failedAt = suite.failedCount
	//		}
	//	}
	//}

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

func (suite *SpecSuite) setOutput(w io.Writer) {
	suite.out = w
}

func (suite *SpecSuite) setPrintFilenames() {
	suite.printFilenames = true
}

func (suite *SpecSuite) setParallel(done func()) {
	suite.parallel = true
	suite.done = done
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
