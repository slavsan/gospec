// Package gospec is a testing library inspired by rspec, mocha but also Cucumber testing libraries which follow
// the BDD approach for testing.
//
// The library provides two flavours of syntax for defining tests: spec vs feature, i.e. rspec vs cucumber style.
//
// Both of those approaches can fit better in different scenarios or depending on the personal preference.
// Parallel execution of tests is supported but requires usage of the [SpecSuite.ParallelAPI] and [FeatureSuite.ParallelAPI] methods,
// and also the [World] construct, because state between defined steps needs to be passed using a "world" instance (per test)
// which carries the test's state with it.
//
// Supported options are:
//   - specifying a custom output writer via the [Output] option.
//   - print filepath:line in the test output via the [PrintedFilenames] option.
//
// Opposite to the Gherkin-first approach, which used in the Cucumber library, gospec follows a
// code-first approach by providing an API which exposes a "Feature", "Scenario", etc, functions,
// which assist in structuring tests so that they look as Gherkin, but are actually Go code. Gospec
// then generates Gherkin (albeit not exactly) as output.
package gospec
