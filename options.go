package gospec

import "io"

// SuiteOption is a type defining an option for controlling the behaviour of [SpecSuite] or [FeatureSuite] instances.
// The available options are: [Output], [Indent] and [PrintedFilenames].
type SuiteOption func(suiteInterface SuiteInterface)

// SuiteInterface is an interface implemented by both [SpecSuite] and [FeatureSuite] suites. It is internal
// and is used by the available [SuiteOption] implementations.
type SuiteInterface interface {
	setOutput(out io.Writer, outputOptions ...OutputOption)
	setIndent(step string)
}

// Output is an option which provides the ability to capture the SpecSuite output in a custom
// [io.Writer]. By default, the output would get printed in [os.Stdout].
//
// [io.Writer]: https://pkg.go.dev/io#Writer
// [os.Stdout]: https://pkg.go.dev/os#Stdout
func Output(w io.Writer, outputOptions ...OutputOption) SuiteOption {
	return func(suite SuiteInterface) {
		switch s := suite.(type) {
		case *SpecSuite:
			s.t.Helper()
		case *FeatureSuite:
			s.t.Helper()
		}
		suite.setOutput(w, outputOptions...)
	}
}

// Indent is an option which sets the indent style, one of: [TwoSpaces], [FourSpaces], [OneTab].
func Indent(step string) SuiteOption {
	return func(suite SuiteInterface) {
		suite.setIndent(step)
	}
}
