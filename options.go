package gospec

import "io"

type SuiteOption func(suiteInterface SuiteInterface)

type SuiteInterface interface {
	setOutput(io.Writer)
	setParallel()
	setPrintFilenames()
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

func WithPrintedFilenames() SuiteOption {
	return func(suite SuiteInterface) {
		suite.setPrintFilenames()
	}
}
