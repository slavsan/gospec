package gospec

import "io"

// SuiteOption is a type defining an option for controlling the behaviour of [SpecSuite] or [FeatureSuite] instances.
// The available option is: [Output].
type SuiteOption func(suiteInterface SuiteInterface)

// SuiteInterface is an interface implemented by both [SpecSuite] and [FeatureSuite] suites. It is internal
// and is used by the available [SuiteOption] implementations.
type SuiteInterface interface {
	setOutput(out io.Writer, outputOptions ...OutputOption)
}

// Output is an option which provides the ability to capture the SpecSuite output in a custom
// [io.Writer].
//
// The available options for [Output] are:
//   - [Colorful]
//   - [Durations]
//   - [PrintFilenames]
//   - [IndentTwoSpaces]
//   - [IndentFourSpaces]
//   - [IndentOneTab]
//
// If there is no Output option specified, by default, the output would get printed in [os.Stdout], with the [Colorful], [Durations] and [IndentTwoSpaces] enabled.
// When a single Output option is defined, it will overwrite the default setting entirely.
//
// One can set multiple outputs, to enable both stdout output and writing to a file, in case saving the spec output is a desired option. This can be done like this:
//
//	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0600)
//	check(err)
//	defer func() { _ = f.Close() }()
//
//	options := []gospec.OutputOptions{
//		gospec.Output(os.Stdout, gospec.Colorful, gospec.Durations),
//		gospec.Output(f),
//	}
//	describe, beforeEach, it := s.With(options...).API()
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
