package gospec

import "io"

// SuiteOption is a type defining an option for controlling the behaviour of [Suite] or [FeatureSuite] instances.
// The available options are: [WithOutput], [WithParallel], and [WithPrintedFilenames].
type SuiteOption func(suiteInterface SuiteInterface)

// SuiteInterface is an interface implemented by both [Suite] and [FeatureSuite] suites. It is internal
// and is used by the available [SuiteOption] implementations.
type SuiteInterface interface {
	setOutput(out io.Writer)
	setParallel(done func())
	setPrintFilenames()
}

// WithOutput is an option which provides the ability to capture the Suite output in a custom
// [io.Writer]. By default, the output would get printed in [os.Stdout].
//
// [io.Writer]: https://pkg.go.dev/io#Writer
// [os.Stdout]: https://pkg.go.dev/os#Stdout
func WithOutput(w io.Writer) SuiteOption {
	return func(suite SuiteInterface) {
		suite.setOutput(w)
	}
}

// WithParallel is an option which enables parallel execution of tests. It's equivalent to calling `t.Parallel()`
// in standard Go tests. All the tests included in a top-level [Suite] or [FeatureSuite] block would execute
// in parallel once this option is enabled for the suite, e.g. via:
//
//	spec := gospec.NewTestSuite(t, gospec.WithParallel())
//
// One important requirement for using this option is to also use the [World] instance which gets passed to each
// block's function, at least the ones that are supposed to execute code.
func WithParallel(done func()) SuiteOption {
	return func(suite SuiteInterface) {
		suite.setParallel(done)
	}
}

// WithPrintedFilenames is an option which enables additional printing of the filename and
// line number (`path/to/filename:line` format) which may come in handy in case your editor/IDE
// supports filepath recognition, with ability to navigate to the source code on click.
func WithPrintedFilenames() SuiteOption {
	return func(suite SuiteInterface) {
		suite.setPrintFilenames()
	}
}
