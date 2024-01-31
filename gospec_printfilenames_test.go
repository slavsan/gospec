package gospec

import (
	"bytes"
	"strings"
	"testing"

	"github.com/slavsan/gospec/internal/testing/helpers/assert"
)

func TestSuiteWithPrintedFilenames(t *testing.T) {
	var (
		out  bytes.Buffer
		spec *SpecSuite
		tm   = &mock{t: t}
	)

	func() {
		WithSpecSuite(t, func(s *SpecSuite) {
			spec = s
			s.t = tm

			describe, _, it := s.With(Output(&out), PrintedFilenames()).API()

			describe("describe 1", func() {
				it("it 1", func(t *T, w *W) {})
				it("it 2", func(t *T, w *W) {})
			})

			describe("describe 2", func() {
				it("it 3", func(t *T, w *W) {})
				it("it 4", func(t *T, w *W) {})
			})
		})
	}()

	assert.Equal(t, [][]any(nil), tm.calls)
	assert.Equal(t, []string{
		"describe 1/it 1",
		"describe 1/it 2",
		"describe 2/it 3",
		"describe 2/it 4",
	}, tm.testTitles)
	assert.Equal(t, 0, len(spec.stack))
	assert.Equal(t, 4, len(spec.suites))
	assert.Equal(t, 2, len(spec.suites[0]))
	assert.Equal(t, "describe 1", spec.suites[0][0].title)
	assert.Equal(t, "it 1", spec.suites[0][1].title)
	assert.Equal(t, 2, len(spec.suites[1]))
	assert.Equal(t, "describe 1", spec.suites[1][0].title)
	assert.Equal(t, "it 2", spec.suites[1][1].title)
	assert.Equal(t, 2, len(spec.suites[2]))
	assert.Equal(t, "describe 2", spec.suites[2][0].title)
	assert.Equal(t, "it 3", spec.suites[2][1].title)
	assert.Equal(t, 2, len(spec.suites[3]))
	assert.Equal(t, "describe 2", spec.suites[3][0].title)
	assert.Equal(t, "it 4", spec.suites[3][1].title)
	assert.Equal(t, strings.Join([]string{
		`describe 1	gospec_printfilenames_test.go:25`,
		`	✔ it 1	gospec_printfilenames_test.go:26`,
		`	✔ it 2	gospec_printfilenames_test.go:27`,
		``,
		`describe 2	gospec_printfilenames_test.go:30`,
		`	✔ it 3	gospec_printfilenames_test.go:31`,
		`	✔ it 4	gospec_printfilenames_test.go:32`,
		``,
		``,
	}, "\n"), out.String())
}

func TestFeatureSuiteWithPrintedFilenames(t *testing.T) {
	var (
		out  bytes.Buffer
		spec *FeatureSuite
		tm   = &mock{t: t}
	)

	func() {
		WithFeatureSuite(t, func(s *FeatureSuite) {
			spec = s
			s.t = tm
			feature, background, scenario, given, when, then, _ := s.With(Output(&out), PrintedFilenames()).API()

			feature("feature 1", func() {
				background(func() {
					given("given 1", func(t *T, w *W) {})
					given("given 2", func(t *T, w *W) {})
				})

				scenario("scenario 1", func() {
					given("given 3", func(t *T, w *W) {})
					when("when 1", func(t *T, w *W) {})
					then("then 1", func(t *T, w *W) {})
				})
			})

			feature("feature 2", func() {
				background(func() {
					given("given 12", func(t *T, w *W) {})
				})

				scenario("scenario 11", func() {
					given("given 13", func(t *T, w *W) {})
					when("when 11", func(t *T, w *W) {})
					then("then 11", func(t *T, w *W) {})
				})
			})
		})
	}()

	assert.Equal(t, [][]any(nil), tm.calls)
	assert.Equal(t, []string{
		"feature 1/scenario 1",
		"feature 2/scenario 11",
	}, tm.testTitles)
	assert.Equal(t, 0, len(spec.stack))
	assert.Equal(t, 2, len(spec.suites))
	assert.Equal(t, 8, len(spec.suites[0]))
	assert.Equal(t, "feature 1", spec.suites[0][0].title)
	assert.Equal(t, "", spec.suites[0][1].title)
	assert.Equal(t, "given 1", spec.suites[0][2].title)
	assert.Equal(t, "given 2", spec.suites[0][3].title)
	assert.Equal(t, "scenario 1", spec.suites[0][4].title)
	assert.Equal(t, "given 3", spec.suites[0][5].title)
	assert.Equal(t, "when 1", spec.suites[0][6].title)
	assert.Equal(t, "then 1", spec.suites[0][7].title)
	assert.Equal(t, 7, len(spec.suites[1]))
	assert.Equal(t, "feature 2", spec.suites[1][0].title)
	assert.Equal(t, "", spec.suites[1][1].title)
	assert.Equal(t, "given 12", spec.suites[1][2].title)
	assert.Equal(t, "scenario 11", spec.suites[1][3].title)
	assert.Equal(t, "given 13", spec.suites[1][4].title)
	assert.Equal(t, "when 11", spec.suites[1][5].title)
	assert.Equal(t, "then 11", spec.suites[1][6].title)
	assert.Equal(t, strings.Join([]string{
		`Feature: feature 1	gospec_printfilenames_test.go:84`,
		``,
		`	Background:	gospec_printfilenames_test.go:85`,
		`		Given given 1	gospec_printfilenames_test.go:86`,
		`		Given given 2	gospec_printfilenames_test.go:87`,
		``,
		`	Scenario: scenario 1	gospec_printfilenames_test.go:90`,
		`		Given given 3	gospec_printfilenames_test.go:91`,
		`		When when 1	gospec_printfilenames_test.go:92`,
		`		Then then 1	gospec_printfilenames_test.go:93`,
		``,
		`Feature: feature 2	gospec_printfilenames_test.go:97`,
		``,
		`	Background:	gospec_printfilenames_test.go:98`,
		`		Given given 12	gospec_printfilenames_test.go:99`,
		``,
		`	Scenario: scenario 11	gospec_printfilenames_test.go:102`,
		`		Given given 13	gospec_printfilenames_test.go:103`,
		`		When when 11	gospec_printfilenames_test.go:104`,
		`		Then then 11	gospec_printfilenames_test.go:105`,
		``,
		``,
	}, "\n"), out.String())
}
