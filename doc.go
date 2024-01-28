// Package gospec is a testing library inspired by rspec, mocha but also Cucumber testing libraries which follow
// the BDD approach for testing.
//
// The library provides two flavours of syntax for defining tests: spec vs feature, i.e. rspec vs cucumber style.
//
// Both of those approaches can fit better in different scenarios or depending on the personal preference of the
// user. Parallel execution of tests is supported but requires usage of the [World] construct and enabling the
// [Parallel] option.
//
// Supported options are:
//   - specifying a custom output writer via the [Output] option.
//   - parallel tests execution using the [Parallel] option.
//   - print filepath:line in the test output via the [PrintedFilenames] option.
//
// Opposite to the Gherkin-first approach used in Cucumber tests, gospec follows the
// code-first approach which then generates Gherkin looking output.
package gospec
