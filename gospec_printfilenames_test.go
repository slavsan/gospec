package gospec

import (
	"bytes"
	"strings"
	"testing"

	"github.com/slavsan/gospec/internal/testing/helpers/assert"
)

func TestWithPrintedFilenames(t *testing.T) {
	var (
		out      bytes.Buffer
		tm       = &mock{t: t}
		spec     = NewTestSuite(tm, WithOutput(&out), WithPrintedFilenames())
		describe = spec.Describe
		it       = spec.It
	)

	describe("describe 1", func() {
		it("it 1", func() {})
		it("it 2", func() {})
	})

	describe("describe 2", func() {
		it("it 3", func() {})
		it("it 4", func() {})
	})

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
		`describe 1	gospec_printfilenames_test.go:20`,
		`	it 1	gospec_printfilenames_test.go:21`,
		`	it 2	gospec_printfilenames_test.go:22`,
		``,
		`describe 2	gospec_printfilenames_test.go:25`,
		`	it 3	gospec_printfilenames_test.go:26`,
		`	it 4	gospec_printfilenames_test.go:27`,
		``,
		``,
	}, "\n"), out.String())
}
