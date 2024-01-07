package gospec

import (
	"fmt"
	"strings"
	"testing"
)

type testingInterface interface {
	Helper()
	Errorf(format string, args ...interface{})
	Run(name string, f func(t *testing.T)) bool
}

type Suite struct {
	t      testingInterface
	steps  []*step
	indent int
}

type block int

const (
	isDescribe block = iota
	isBeforeEach
	isIt
)

type step struct {
	indent int
	block  block
	title  string
	cb     func()
}

func NewTestSuite(t *testing.T) *Suite {
	return &Suite{
		t:      t,
		steps:  []*step{},
		indent: 0,
	}
}

func (suite *Suite) report() {
	for _, s := range suite.steps {
		if s.block == isBeforeEach {
			continue
		}
		if s.block == isIt {
			fmt.Printf("%s✓ %s\n", strings.Repeat("  ", s.indent+1), s.title)
			continue
		}
		fmt.Printf("%s%s\n", strings.Repeat("  ", s.indent), s.title)
	}
	fmt.Printf("\n")
}

func (suite *Suite) buildSuites() [][]*step {
	var suites [][]*step
	stack := []*step{}
	indent := 0

	copyStack := func() {
		suite := []*step{}
		if len(suites) > 0 {
			lastSuite := suites[len(suites)-1]
			if stack[len(stack)-1] == lastSuite[len(lastSuite)-1] {
				return
			}
		}
		for _, step := range stack {
			suite = append(suite, step)
		}
		suites = append(suites, suite)
	}

	lastSuiteEndsWithIt := func() bool {
		suite := suites[len(suites)-1]
		return suite[len(suite)-1].block == isIt
	}

	lastSuiteStartsWithStep := func(s *step) bool {
		if len(suites) == 0 {
			return false
		}
		suite := suites[len(suites)-1]
		return suite[0] == s
	}

	findLastSiblingIndex := func(s *step) int {
		lastDescribeIndex := len(stack)
		for i := len(stack) - 1; i >= 0; i-- {
			step := stack[i]
			lastDescribeIndex--
			if step.indent == s.indent && step.block == isDescribe {
				break
			}
		}
		return lastDescribeIndex
	}

	isNextStepSibling := func(s *step, i int) bool {
		if i+1 < len(suite.steps) { // if there is a next step
			next := suite.steps[i+1] // peek
			if next.indent == s.indent && next.block == isDescribe {
				return true
			}
		}
		return false
	}

	for i := 0; i < len(suite.steps); i++ {
		s := suite.steps[i]
		if s.block == isDescribe && s.indent < indent {
			if !lastSuiteEndsWithIt() {
				copyStack()
			}
			lastDescribeIndex := findLastSiblingIndex(s)
			if lastDescribeIndex != 0 {
				stack = stack[:lastDescribeIndex]
				indent = stack[len(stack)-1].indent
			}
			stack = append(stack, s)
			continue
		}
		if s.block == isDescribe && s.indent == indent {
			if !lastSuiteEndsWithIt() {
				copyStack()
			}
			stack = stack[:findLastSiblingIndex(s)]
			stack = append(stack, s)
			continue
		}
		if s.block == isDescribe && s.indent > indent {
			stack = append(stack, s)
			indent++
			if isNextStepSibling(s, i) {
				copyStack()
			}
			continue
		}
		if s.block == isBeforeEach {
			stack = append(stack, s)
			continue
		}
		if s.block == isIt {
			stack = append(stack, s)
			copyStack()
			stack = stack[:len(stack)-1]
			if i+1 < len(suite.steps) {
				next := suite.steps[i+1] // peek
				if len(stack) > 0 && next.indent < stack[len(stack)-1].indent {
					stack = stack[:findLastSiblingIndex(s)]
				}
			} else {
				stack = stack[:findLastSiblingIndex(s)]
			}
			continue
		}
	}

	if len(stack) > 0 {
		if len(stack) == 1 && lastSuiteStartsWithStep(stack[0]) {
		} else {
			copyStack()
		}
	}

	return suites
}

func (suite *Suite) Start() {
	suite.report()

	suites := suite.buildSuites()

	// debugSuitesAndSteps(suites)

	for _, childSuite := range suites {
		suite.t.Run(buildSuiteTitle(childSuite), func(t *testing.T) {
			for _, step := range childSuite {
				step.cb()
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

// func debugSuitesAndSteps(suites [][]*step) {
// 	for i, suite := range suites {
// 		fmt.Printf("SUITE: %d\n", i)
// 		for _, step := range suite {
// 			fmt.Printf("STEP: %#v\n", step)
// 		}
// 		fmt.Println()
// 	}
// }

func (suite *Suite) Describe(title string, cb func()) {
	suite.indent++
	suite.steps = append(suite.steps, &step{
		title:  title,
		indent: suite.indent,
		block:  isDescribe,
		cb:     cb,
	})

	cb()
	suite.indent--
}

func (suite *Suite) BeforeEach(cb func()) {
	suite.steps = append(suite.steps, &step{
		indent: suite.indent,
		block:  isBeforeEach,
		cb:     cb,
	})
}

func (suite *Suite) It(title string, cb func()) {
	suite.steps = append(suite.steps, &step{
		title:  title,
		indent: suite.indent,
		block:  isIt,
		cb:     cb,
	})
}
