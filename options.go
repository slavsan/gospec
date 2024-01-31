package gospec

import "io"

// SuiteOption is a type defining an option for controlling the behaviour of [SpecSuite] or [FeatureSuite] instances.
// The available options are: [Output], [Indent] and [PrintedFilenames].
type SuiteOption func(suiteInterface SuiteInterface)

// SuiteInterface is an interface implemented by both [SpecSuite] and [FeatureSuite] suites. It is internal
// and is used by the available [SuiteOption] implementations.
type SuiteInterface interface {
	setOutput(out io.Writer)
	setPrintFilenames()
	setIndent(step string)
}

// Output is an option which provides the ability to capture the SpecSuite output in a custom
// [io.Writer]. By default, the output would get printed in [os.Stdout].
//
// [io.Writer]: https://pkg.go.dev/io#Writer
// [os.Stdout]: https://pkg.go.dev/os#Stdout
func Output(w io.Writer) SuiteOption {
	return func(suite SuiteInterface) {
		suite.setOutput(w)
	}
}

// PrintedFilenames is an option which enables additional printing of the filename and
// line number (`path/to/filename:line` format) which may come in handy in case your editor/IDE
// supports filepath recognition, with ability to navigate to the source code on click.
func PrintedFilenames() SuiteOption {
	return func(suite SuiteInterface) {
		suite.setPrintFilenames()
	}
}

// Indent is an option which sets the indent style, one of: [TwoSpaces], [FourSpaces], [OneTab].
func Indent(step string) SuiteOption {
	return func(suite SuiteInterface) {
		suite.setIndent(step)
	}
}
