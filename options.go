package gospec

import "io"

type SuiteOption func(suiteInterface SuiteInterface)

type SuiteInterface interface {
	setOutput(io.Writer)
	setParallel()
}

func WithOutput(w io.Writer) SuiteOption {
	return func(suite SuiteInterface) {
		suite.setOutput(w)
	}
}

func WithParallel() SuiteOption {
	return func(suite SuiteInterface) {
		suite.setParallel()
	}
}
