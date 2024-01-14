package gospec

import (
	"bytes"
	"strings"
	"testing"

	"github.com/slavsan/gospec/internal/testing/helpers/assert"
)

func TestDescribesCanBeSetAtTopLevel(t *testing.T) {
	var (
		out      bytes.Buffer
		tm       = &mock{t: t}
		spec     = NewTestSuite(tm, WithOutput(&out))
		describe = spec.Describe
	)

	describe("describe 1", func() {})
	describe("describe 2", func() {})

	assert.Equal(t, [][]any(nil), tm.calls)
	assert.Equal(t, []string(nil), tm.testTitles)
	assert.Equal(t, 0, len(spec.stack))
	assert.Equal(t, 0, len(spec.suites))
	assert.Equal(t, strings.Join([]string{
		`describe 1`,
		``,
		`describe 2`,
		``,
		``,
	}, "\n"), out.String())
}

func TestDescribesCanBeNested(t *testing.T) {
	var (
		out      bytes.Buffer
		tm       = &mock{t: t}
		spec     = NewTestSuite(tm, WithOutput(&out))
		describe = spec.Describe
	)

	describe("Checkout", func() {
		describe("nested describe", func() {
		})
	})

	assert.Equal(t, [][]any(nil), tm.calls)
	assert.Equal(t, []string(nil), tm.testTitles)
	assert.Equal(t, 0, len(spec.stack))
	assert.Equal(t, 0, len(spec.suites))
	assert.Equal(t, strings.Join([]string{
		`Checkout`,
		`	nested describe`,
		``,
		``,
	}, "\n"), out.String())
}

func TestTwoTopLevelDescribesWithTwoNestedDescribes(t *testing.T) {
	var (
		out      bytes.Buffer
		tm       = &mock{t: t}
		spec     = NewTestSuite(tm, WithOutput(&out))
		describe = spec.Describe
	)

	describe("describe 1", func() {
		describe("nested 1", func() {})
		describe("nested 2", func() {})
	})

	describe("describe 2", func() {
		describe("nested 3", func() {})
		describe("nested 4", func() {})
	})

	assert.Equal(t, [][]any(nil), tm.calls)
	assert.Equal(t, []string(nil), tm.testTitles)
	assert.Equal(t, 0, len(spec.stack))
	assert.Equal(t, 0, len(spec.suites))
	assert.Equal(t, strings.Join([]string{
		`describe 1`,
		`	nested 1`,
		`	nested 2`,
		``,
		`describe 2`,
		`	nested 3`,
		`	nested 4`,
		``,
		``,
	}, "\n"), out.String())
}

func TestTwoTopLevelDescribesWithThreeLevelsNestedDescribes(t *testing.T) {
	var (
		out      bytes.Buffer
		tm       = &mock{t: t}
		spec     = NewTestSuite(tm, WithOutput(&out))
		describe = spec.Describe
	)

	describe("describe 1", func() {
		describe("nested 1", func() {
			describe("nested 2", func() {})
		})
	})

	describe("describe 2", func() {
		describe("nested 3", func() {
			describe("nested 4", func() {})
		})
	})

	assert.Equal(t, [][]any(nil), tm.calls)
	assert.Equal(t, []string(nil), tm.testTitles)
	assert.Equal(t, 0, len(spec.stack))
	assert.Equal(t, 0, len(spec.suites))
	assert.Equal(t, strings.Join([]string{
		`describe 1`,
		`	nested 1`,
		`		nested 2`,
		``,
		`describe 2`,
		`	nested 3`,
		`		nested 4`,
		``,
		``,
	}, "\n"), out.String())
}

func TestDescribesNestingComplexExample(t *testing.T) {
	var (
		out      bytes.Buffer
		tm       = &mock{t: t}
		spec     = NewTestSuite(tm, WithOutput(&out))
		describe = spec.Describe
	)

	describe("describe 1", func() {
		describe("nested 1", func() {
			describe("nested 2", func() {})
			describe("nested 3", func() {
				describe("nested 4", func() {})
				describe("nested 5", func() {
					describe("nested 6", func() {})
				})
			})
			describe("nested 7", func() {})
			describe("nested 8", func() {})
		})
	})

	describe("describe 9", func() {
		describe("nested 10", func() {
			describe("nested 11", func() {})
		})
		describe("nested 12", func() {})
		describe("nested 13", func() {})
		describe("nested 14", func() {
			describe("nested 15", func() {
				describe("nested 16", func() {})
			})
			describe("nested 17", func() {})
		})
		describe("nested 18", func() {})
	})

	assert.Equal(t, [][]any(nil), tm.calls)
	assert.Equal(t, []string(nil), tm.testTitles)
	assert.Equal(t, 0, len(spec.stack))
	assert.Equal(t, 0, len(spec.suites))
	assert.Equal(t, strings.Join([]string{
		`describe 1`,
		`	nested 1`,
		`		nested 2`,
		`		nested 3`,
		`			nested 4`,
		`			nested 5`,
		`				nested 6`,
		`		nested 7`,
		`		nested 8`,
		``,
		`describe 9`,
		`	nested 10`,
		`		nested 11`,
		`	nested 12`,
		`	nested 13`,
		`	nested 14`,
		`		nested 15`,
		`			nested 16`,
		`		nested 17`,
		`	nested 18`,
		``,
		``,
	}, "\n"), out.String())
}

func TestDescribesWithBeforeEach(t *testing.T) {
	var (
		out        bytes.Buffer
		tm         = &mock{t: t}
		spec       = NewTestSuite(tm, WithOutput(&out))
		describe   = spec.Describe
		beforeEach = spec.BeforeEach
	)

	describe("describe 1", func() {
		beforeEach(func() {})
		describe("nested 1", func() {})
		describe("nested 2", func() {})
	})

	describe("describe 2", func() {
		beforeEach(func() {})
		describe("nested 3", func() {})
		describe("nested 4", func() {})
	})

	assert.Equal(t, [][]any(nil), tm.calls)
	assert.Equal(t, []string(nil), tm.testTitles)
	assert.Equal(t, 0, len(spec.stack))
	assert.Equal(t, 0, len(spec.suites))
	assert.Equal(t, strings.Join([]string{
		`describe 1`,
		`	nested 1`,
		`	nested 2`,
		``,
		`describe 2`,
		`	nested 3`,
		`	nested 4`,
		``,
		``,
	}, "\n"), out.String())
}

func TestSingleDescribeWithSingleItBlock(t *testing.T) {
	var (
		out      bytes.Buffer
		tm       = &mock{t: t}
		spec     = NewTestSuite(tm, WithOutput(&out))
		describe = spec.Describe
		it       = spec.It
	)

	describe("describe 1", func() {
		it("it 1", func() {})
	})

	assert.Equal(t, [][]any(nil), tm.calls)
	assert.Equal(t, []string{"describe 1/it 1"}, tm.testTitles)
	assert.Equal(t, 0, len(spec.stack))
	assert.Equal(t, 1, len(spec.suites))
	assert.Equal(t, 2, len(spec.suites[0]))
	assert.Equal(t, "describe 1", spec.suites[0][0].title)
	assert.Equal(t, "it 1", spec.suites[0][1].title)
	assert.Equal(t, strings.Join([]string{
		`describe 1`,
		`	it 1`,
		``,
		``,
	}, "\n"), out.String())
}

func TestSingleDescribeWithTwoItBlocks(t *testing.T) {
	var (
		out      bytes.Buffer
		tm       = &mock{t: t}
		spec     = NewTestSuite(tm, WithOutput(&out))
		describe = spec.Describe
		it       = spec.It
	)

	describe("describe 1", func() {
		it("it 1", func() {})
		it("it 2", func() {})
	})

	assert.Equal(t, [][]any(nil), tm.calls)
	assert.Equal(t, []string{"describe 1/it 1", "describe 1/it 2"}, tm.testTitles)
	assert.Equal(t, 0, len(spec.stack))
	assert.Equal(t, 2, len(spec.suites))
	assert.Equal(t, 2, len(spec.suites[0]))
	assert.Equal(t, "describe 1", spec.suites[0][0].title)
	assert.Equal(t, "it 1", spec.suites[0][1].title)
	assert.Equal(t, 2, len(spec.suites[1]))
	assert.Equal(t, "describe 1", spec.suites[1][0].title)
	assert.Equal(t, "it 2", spec.suites[1][1].title)
	assert.Equal(t, strings.Join([]string{
		`describe 1`,
		`	it 1`,
		`	it 2`,
		``,
		``,
	}, "\n"), out.String())
}

func TestTwoDescribeBlocksWithTwoItBlocks(t *testing.T) {
	var (
		out      bytes.Buffer
		tm       = &mock{t: t}
		spec     = NewTestSuite(tm, WithOutput(&out))
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
		`describe 1`,
		`	it 1`,
		`	it 2`,
		``,
		`describe 2`,
		`	it 3`,
		`	it 4`,
		``,
		``,
	}, "\n"), out.String())
}

func TestTwoDescribeBlocksWithBothNestedDescribesAndItBlocks(t *testing.T) {
	var (
		out        bytes.Buffer
		tm         = &mock{t: t}
		spec       = NewTestSuite(tm, WithOutput(&out))
		describe   = spec.Describe
		beforeEach = spec.BeforeEach
		it         = spec.It
	)

	describe("describe 1", func() {
		it("it 1", func() {})
		it("it 2", func() {})

		describe("describe 2", func() {
			beforeEach(func() {})
			it("it 3", func() {})
			it("it 4", func() {})
		})
	})

	describe("describe 3", func() {
		it("it 5", func() {})
		it("it 6", func() {})
	})

	assert.Equal(t, [][]any(nil), tm.calls)
	assert.Equal(t, []string{
		"describe 1/it 1",
		"describe 1/it 2",
		"describe 1/describe 2/it 3",
		"describe 1/describe 2/it 4",
		"describe 3/it 5",
		"describe 3/it 6",
	}, tm.testTitles)
	assert.Equal(t, 0, len(spec.stack))
	assert.Equal(t, 6, len(spec.suites))
	assert.Equal(t, 2, len(spec.suites[0]))
	assert.Equal(t, "describe 1", spec.suites[0][0].title)
	assert.Equal(t, "it 1", spec.suites[0][1].title)
	assert.Equal(t, 2, len(spec.suites[1]))
	assert.Equal(t, "describe 1", spec.suites[1][0].title)
	assert.Equal(t, "it 2", spec.suites[1][1].title)
	assert.Equal(t, 4, len(spec.suites[2]))
	assert.Equal(t, "describe 1", spec.suites[2][0].title)
	assert.Equal(t, "describe 2", spec.suites[2][1].title)
	assert.Equal(t, "", spec.suites[2][2].title)
	assert.Equal(t, isBeforeEach, spec.suites[2][2].block)
	assert.Equal(t, "it 3", spec.suites[2][3].title)
	assert.Equal(t, 4, len(spec.suites[3]))
	assert.Equal(t, "describe 1", spec.suites[3][0].title)
	assert.Equal(t, "describe 2", spec.suites[3][1].title)
	assert.Equal(t, "", spec.suites[3][2].title)
	assert.Equal(t, isBeforeEach, spec.suites[3][2].block)
	assert.Equal(t, "it 4", spec.suites[3][3].title)
	assert.Equal(t, 2, len(spec.suites[4]))
	assert.Equal(t, "describe 3", spec.suites[4][0].title)
	assert.Equal(t, "it 5", spec.suites[4][1].title)
	assert.Equal(t, 2, len(spec.suites[5]))
	assert.Equal(t, "describe 3", spec.suites[5][0].title)
	assert.Equal(t, "it 6", spec.suites[5][1].title)
	assert.Equal(t, strings.Join([]string{
		`describe 1`,
		`	it 1`,
		`	it 2`,
		`	describe 2`,
		`		it 3`,
		`		it 4`,
		``,
		`describe 3`,
		`	it 5`,
		`	it 6`,
		``,
		``,
	}, "\n"), out.String())
}

func TestSequentialExecution(t *testing.T) {
	var (
		out         bytes.Buffer
		testingMock = &mock{t: t}
		mockAssert  = &assertMock{}
		spec        = NewTestSuite(testingMock, WithOutput(&out))
		describe    = spec.Describe
		beforeEach  = spec.BeforeEach
		it          = spec.It
	)

	describe("describe 1", func() {
		var nums []int

		beforeEach(func() {
			nums = []int{}
			nums = append(nums, 1)
		})

		describe("describe 2", func() {
			beforeEach(func() {
				nums = append(nums, 2)
			})

			it("it 1", func() {
				nums = append(nums, 3)
				mockAssert.Assert(nums)
			})

			it("it 2", func() {
				nums = append(nums, 4)
				mockAssert.Assert(nums)
			})
		})
	})

	describe("describe 3", func() {
		var nums []int

		beforeEach(func() {
			nums = []int{}
			nums = append(nums, 11)
		})

		describe("describe 4", func() {
			beforeEach(func() {
				nums = append(nums, 12)
			})

			it("it 3", func() {
				nums = append(nums, 13)
				mockAssert.Assert(nums)
			})

			it("it 4", func() {
				nums = append(nums, 14)
				mockAssert.Assert(nums)
			})
		})
	})

	assert.Equal(t, [][]any(nil), testingMock.calls)
	assert.Equal(t, 4, len(mockAssert.calls))
	assert.Equal(t, []any{[]int{1, 2, 3}}, mockAssert.calls[0])
	assert.Equal(t, []any{[]int{1, 2, 4}}, mockAssert.calls[1])
	assert.Equal(t, []any{[]int{11, 12, 13}}, mockAssert.calls[2])
	assert.Equal(t, []any{[]int{11, 12, 14}}, mockAssert.calls[3])

	assert.Equal(t, []string{
		"describe 1/describe 2/it 1",
		"describe 1/describe 2/it 2",
		"describe 3/describe 4/it 3",
		"describe 3/describe 4/it 4",
	}, testingMock.testTitles)

	assert.Equal(t, strings.Join([]string{
		`describe 1`,
		`	describe 2`,
		`		it 1`,
		`		it 2`,
		``,
		`describe 3`,
		`	describe 4`,
		`		it 3`,
		`		it 4`,
		``,
		``,
	}, "\n"), out.String())
}
