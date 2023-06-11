package gospec

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type Suite struct {
	t      *testing.T
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
		indent: -1,
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

func (suite *Suite) Start() {
	suite.report()

	var suites [][]*step
	stack := []*step{}
	indent := -1

	copyStack := func() {
		suite := []*step{}
		for _, step := range stack {
			suite = append(suite, step)
		}
		suites = append(suites, suite)
	}

	lastSuiteEndsWithIt := func() bool {
		suite := suites[len(suites)-1]
		return suite[len(suite)-1].block == isIt
	}

	for _, s := range suite.steps {
		if s.block == isDescribe && s.indent < indent {
			if !lastSuiteEndsWithIt() {
				copyStack()
			}

			lastDescribeIndex := len(stack) - 1
			for i := len(stack) - 1; i >= 0; i-- {
				step := stack[i]
				lastDescribeIndex--
				// TODO: only when we reach the describe block ?
				if step.indent == s.indent {
					break
				}
			}

			if lastDescribeIndex == -1 {
				stack = []*step{}
				indent = -1
			} else {
				stack = stack[:lastDescribeIndex]
				indent = stack[len(stack)-1].indent + 1
			}
			stack = append(stack, s)
			continue
		}
		if s.block == isDescribe && s.indent == indent {
			copyStack()
			stack = stack[:len(stack)-1]
			stack = append(stack, s)
			continue
		}
		if s.block == isDescribe && s.indent > indent {
			stack = append(stack, s)
			indent++
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
			continue
		}
		indent--
	}

	debugSuitesAndSteps(suites)

	for _, childSuite := range suites {
		suite.t.Run("", func(t *testing.T) {
			for _, step := range childSuite {
				step.cb()
			}
		})
	}
}

func debugSuitesAndSteps(suites [][]*step) {
	for i, suite := range suites {
		fmt.Printf("SUITE: %d\n", i)
		for _, step := range suite {
			fmt.Printf("STEP: %#v\n", step)
		}
		fmt.Println()
	}
}

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

type Chain struct {
	To        *Chain
	Be        *Chain
	Of        *Chain
	Have      *Chain
	Not       *Chain
	Contain   *Chain
	True      func()
	False     func()
	EqualTo   func(expected any)
	LengthOf  func(expected int)
	Property  func(expected any)
	Element   func(expected any)
	Substring func(expected string)
	Type      func(expected any)
	Nil       func()
}

func (suite *Suite) Expect(value any) *Chain {
	suite.t.Helper()
	return &Chain{
		Not: &Chain{
			To: &Chain{
				Be: &Chain{
					Nil: func() {
						// TODO: implement me
					},
				},
			},
		},
		To: &Chain{
			Contain: &Chain{
				Substring: func(sub string) {
					// TODO: implement me
				},
				Element: func(elem any) {
					// TODO: implement me
				},
			},
			Have: &Chain{
				LengthOf: func(length int) {
					suite.t.Helper()

					kind := reflect.TypeOf(value).Kind()

					if kind != reflect.Slice && kind != reflect.Array && kind != reflect.String {
						suite.t.Errorf("expected target to be slice/array/slice but it was %s", kind)
					}

					if kind == reflect.String {
						reflectValue := reflect.ValueOf(value)
						if reflectValue.Len() != length {
							suite.t.Errorf("expected %s to have length %d but it has %d", value, reflectValue.Len(), length)
						}
						return
					}

					reflectValue := reflect.ValueOf(value)
					if reflectValue.Len() != length {
						suite.t.Errorf("expected %s to have length %d but it has %d", value, reflectValue.Len(), length)
					}
				},
				Property: func(prop any) {
					// TODO: implement me
				},
			},
			Be: &Chain{
				Of: &Chain{
					Type: func(expected any) {
						// TODO: implement me
					},
				},
				Nil: func() {
					// TODO: implement me
				},
				True: func() {
					suite.t.Helper()
					v, ok := value.(bool)
					if !ok {
						suite.t.Errorf("expected test target to be bool but it was %s", reflect.TypeOf(value))
						return
					}
					if v == false {
						suite.t.Errorf("expected false but got true")
					}
				},
				False: func() {
					suite.t.Helper()
					v, ok := value.(bool)
					if !ok {
						suite.t.Errorf("expected test target to be bool but it was %s", reflect.TypeOf(value))
						return
					}
					if v != false {
						suite.t.Errorf("expected false but got true")
					}
				},
				EqualTo: func(expected any) {
					// TODO: implement me
				},
			},
		},
	}
}
