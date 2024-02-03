package gospec

import (
	"bytes"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/slavsan/gospec/internal/testing/helpers/assert"
)

type T = testing.T

func TestDescribesCanBeSetAtTopLevel(t *testing.T) {
	var (
		out  bytes.Buffer
		tm   = &mock{t: t}
		spec *SpecSuite
	)

	func() {
		WithSpecSuite(t, func(s *SpecSuite) {
			spec = s
			s.t = tm
			describe, _, _ := s.With(Output(&out)).API()

			describe("describe 1", func() {})
			describe("describe 2", func() {})
		})
	}()

	assert.Equal(t, [][]any(nil), tm.calls)
	assert.Equal(t, []string{"describe 1", "describe 2"}, tm.testTitles)
	assert.Equal(t, 0, len(spec.stack))
	assert.Equal(t, 2, len(spec.suites))
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
		out  bytes.Buffer
		tm   = &mock{t: t}
		spec *SpecSuite
	)

	func() {
		WithSpecSuite(t, func(s *SpecSuite) {
			spec = s
			spec.t = tm
			describe, _, _ := s.With(Output(&out)).API()

			describe("Checkout", func() {
				describe("nested describe", func() {
				})
			})
		})
	}()

	assert.Equal(t, [][]any(nil), tm.calls)
	assert.Equal(t, []string{"Checkout/nested describe"}, tm.testTitles)
	assert.Equal(t, 0, len(spec.stack))
	assert.Equal(t, 1, len(spec.suites))
	assert.Equal(t, strings.Join([]string{
		`Checkout`,
		`  nested describe`,
		``,
		``,
	}, "\n"), out.String())
}

func TestTwoTopLevelDescribesWithTwoNestedDescribes(t *testing.T) {
	var (
		out  bytes.Buffer
		tm   = &mock{t: t}
		spec *SpecSuite
	)

	func() {
		WithSpecSuite(t, func(s *SpecSuite) {
			spec = s
			s.t = tm
			describe, _, _ := s.With(Output(&out)).API()

			describe("describe 1", func() {
				describe("nested 1", func() {})
				describe("nested 2", func() {})
			})

			describe("describe 2", func() {
				describe("nested 3", func() {})
				describe("nested 4", func() {})
			})
		})
	}()

	assert.Equal(t, [][]any(nil), tm.calls)
	assert.Equal(t, []string{
		"describe 1/nested 1",
		"describe 1/nested 2",
		"describe 2/nested 3",
		"describe 2/nested 4",
	}, tm.testTitles)
	assert.Equal(t, 0, len(spec.stack))
	assert.Equal(t, 4, len(spec.suites))
	assert.Equal(t, strings.Join([]string{
		`describe 1`,
		`  nested 1`,
		`  nested 2`,
		``,
		`describe 2`,
		`  nested 3`,
		`  nested 4`,
		``,
		``,
	}, "\n"), out.String())
}

func TestTwoTopLevelDescribesWithThreeLevelsNestedDescribes(t *testing.T) {
	var (
		out  bytes.Buffer
		tm   = &mock{t: t}
		spec *SpecSuite
	)

	func() {
		WithSpecSuite(t, func(s *SpecSuite) {
			spec = s
			s.t = tm
			describe, _, _ := s.With(Output(&out)).API()

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
		})
	}()

	assert.Equal(t, [][]any(nil), tm.calls)
	assert.Equal(t, []string{
		"describe 1/nested 1/nested 2",
		"describe 2/nested 3/nested 4",
	}, tm.testTitles)
	assert.Equal(t, 0, len(spec.stack))
	assert.Equal(t, 2, len(spec.suites))
	assert.Equal(t, strings.Join([]string{
		`describe 1`,
		`  nested 1`,
		`    nested 2`,
		``,
		`describe 2`,
		`  nested 3`,
		`    nested 4`,
		``,
		``,
	}, "\n"), out.String())
}

func TestDescribesNestingComplexExample(t *testing.T) {
	var (
		out  bytes.Buffer
		spec *SpecSuite
		tm   = &mock{t: t}
	)

	func() {
		WithSpecSuite(t, func(s *SpecSuite) {
			spec = s
			s.t = tm
			describe, _, _ := s.With(Output(&out)).API()

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
		})
	}()

	assert.Equal(t, [][]any(nil), tm.calls)
	assert.Equal(t, []string{
		"describe 1/nested 1/nested 2",
		"describe 1/nested 1/nested 3/nested 4",
		"describe 1/nested 1/nested 3/nested 5/nested 6",
		"describe 1/nested 1/nested 7",
		"describe 1/nested 1/nested 8",
		"describe 9/nested 10/nested 11",
		"describe 9/nested 12",
		"describe 9/nested 13",
		"describe 9/nested 14/nested 15/nested 16",
		"describe 9/nested 14/nested 17",
		"describe 9/nested 18",
	}, tm.testTitles)
	assert.Equal(t, 0, len(spec.stack))
	assert.Equal(t, 11, len(spec.suites))
	assert.Equal(t, strings.Join([]string{
		`describe 1`,
		`  nested 1`,
		`    nested 2`,
		`    nested 3`,
		`      nested 4`,
		`      nested 5`,
		`        nested 6`,
		`    nested 7`,
		`    nested 8`,
		``,
		`describe 9`,
		`  nested 10`,
		`    nested 11`,
		`  nested 12`,
		`  nested 13`,
		`  nested 14`,
		`    nested 15`,
		`      nested 16`,
		`    nested 17`,
		`  nested 18`,
		``,
		``,
	}, "\n"), out.String())
}

func TestDescribesWithBeforeEach(t *testing.T) {
	var (
		out  bytes.Buffer
		tm   = &mock{t: t}
		spec *SpecSuite
	)

	func() {
		WithSpecSuite(t, func(s *SpecSuite) {
			spec = s
			s.t = tm
			describe, beforeEach, _ := s.With(Output(&out)).API()

			describe("describe 1", func() {
				beforeEach(func(t *T) {})
				describe("nested 1", func() {})
				describe("nested 2", func() {})
			})

			describe("describe 2", func() {
				beforeEach(func(t *T) {})
				describe("nested 3", func() {})
				describe("nested 4", func() {})
			})
		})
	}()

	assert.Equal(t, [][]any(nil), tm.calls)
	assert.Equal(t, []string{
		"describe 1/nested 1",
		"describe 1/nested 2",
		"describe 2/nested 3",
		"describe 2/nested 4",
	}, tm.testTitles)
	assert.Equal(t, 0, len(spec.stack))
	assert.Equal(t, 4, len(spec.suites))
	assert.Equal(t, strings.Join([]string{
		`describe 1`,
		`  nested 1`,
		`  nested 2`,
		``,
		`describe 2`,
		`  nested 3`,
		`  nested 4`,
		``,
		``,
	}, "\n"), out.String())
}

func TestSingleDescribeWithSingleItBlock(t *testing.T) {
	var (
		out  bytes.Buffer
		spec *SpecSuite
		tm   = &mock{t: t}
	)

	func() {
		WithSpecSuite(t, func(s *SpecSuite) {
			spec = s
			s.t = tm
			describe, _, it := s.With(Output(&out)).API()

			describe("describe 1", func() {
				it("it 1", func(t *T) {})
			})
		})
	}()

	assert.Equal(t, [][]any(nil), tm.calls)
	assert.Equal(t, []string{"describe 1/it 1"}, tm.testTitles)
	assert.Equal(t, 0, len(spec.stack))
	assert.Equal(t, 1, len(spec.suites))
	assert.Equal(t, 2, len(spec.suites[0]))
	assert.Equal(t, "describe 1", spec.suites[0][0].title)
	assert.Equal(t, "it 1", spec.suites[0][1].title)
	assert.Equal(t, strings.Join([]string{
		`describe 1`,
		`  ✔ it 1`,
		``,
		``,
	}, "\n"), out.String())
}

func TestSingleDescribeWithTwoItBlocks(t *testing.T) {
	var (
		out  bytes.Buffer
		spec *SpecSuite
		tm   = &mock{t: t}
	)

	func() {
		WithSpecSuite(t, func(s *SpecSuite) {
			spec = s
			s.t = tm
			describe, _, it := s.With(Output(&out)).API()

			describe("describe 1", func() {
				it("it 1", func(t *T) {})
				it("it 2", func(t *T) {})
			})
		})
	}()

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
		`  ✔ it 1`,
		`  ✔ it 2`,
		``,
		``,
	}, "\n"), out.String())
}

func TestSingleDescribeWithTwoOutputs(t *testing.T) {
	var (
		out  bytes.Buffer
		out2 bytes.Buffer
		spec *SpecSuite
		tm   = &mock{t: t}
	)

	func() {
		WithSpecSuite(t, func(s *SpecSuite) {
			spec = s
			s.t = tm
			describe, _, it := s.With(Output(&out), Output(&out2)).API()

			describe("describe 1", func() {
				it("it 1", func(t *T) {})
				it("it 2", func(t *T) {})
			})
		})
	}()

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
		`  ✔ it 1`,
		`  ✔ it 2`,
		``,
		``,
	}, "\n"), out.String())
	assert.Equal(t, strings.Join([]string{
		`describe 1`,
		`  ✔ it 1`,
		`  ✔ it 2`,
		``,
		``,
	}, "\n"), out2.String())
}

func TestTwoDescribeBlocksWithTwoItBlocks(t *testing.T) {
	var (
		out  bytes.Buffer
		spec *SpecSuite
		tm   = &mock{t: t}
	)

	func() {
		WithSpecSuite(t, func(s *SpecSuite) {
			spec = s
			s.t = tm
			describe, _, it := s.With(Output(&out)).API()

			describe("describe 1", func() {
				it("it 1", func(t *T) {})
				it("it 2", func(t *T) {})
			})

			describe("describe 2", func() {
				it("it 3", func(t *T) {})
				it("it 4", func(t *T) {})
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
		`describe 1`,
		`  ✔ it 1`,
		`  ✔ it 2`,
		``,
		`describe 2`,
		`  ✔ it 3`,
		`  ✔ it 4`,
		``,
		``,
	}, "\n"), out.String())
}

func TestTwoDescribeBlocksWithBothNestedDescribesAndItBlocks(t *testing.T) {
	var (
		out  bytes.Buffer
		spec *SpecSuite
		tm   = &mock{t: t}
	)

	func() {
		WithSpecSuite(t, func(s *SpecSuite) {
			spec = s
			s.t = tm
			describe, beforeEach, it := s.With(Output(&out)).API()

			describe("describe 1", func() {
				it("it 1", func(t *T) {})
				it("it 2", func(t *T) {})

				describe("describe 2", func() {
					beforeEach(func(t *T) {})
					it("it 3", func(t *T) {})
					it("it 4", func(t *T) {})
				})
			})

			describe("describe 3", func() {
				it("it 5", func(t *T) {})
				it("it 6", func(t *T) {})
			})
		})
	}()

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
		`  ✔ it 1`,
		`  ✔ it 2`,
		`  describe 2`,
		`    ✔ it 3`,
		`    ✔ it 4`,
		``,
		`describe 3`,
		`  ✔ it 5`,
		`  ✔ it 6`,
		``,
		``,
	}, "\n"), out.String())
}

func TestSequentialExecution(t *testing.T) {
	var (
		out         bytes.Buffer
		testingMock = &mock{t: t}
		mockAssert  = &assertMock{}
	)

	func() {
		WithSpecSuite(t, func(s *SpecSuite) {
			s.t = testingMock
			describe, beforeEach, it := s.With(Output(&out)).API()

			describe("describe 1", func() {
				var nums []int

				beforeEach(func(t *T) {
					nums = []int{}
					nums = append(nums, 1)
				})

				describe("describe 2", func() {
					beforeEach(func(t *T) {
						nums = append(nums, 2)
					})

					it("it 1", func(t *T) {
						nums = append(nums, 3)
						mockAssert.Assert(nums)
					})

					it("it 2", func(t *T) {
						nums = append(nums, 4)
						mockAssert.Assert(nums)
					})
				})
			})

			describe("describe 3", func() {
				var nums []int

				beforeEach(func(t *T) {
					nums = []int{}
					nums = append(nums, 11)
				})

				describe("describe 4", func() {
					beforeEach(func(t *T) {
						nums = append(nums, 12)
					})

					it("it 3", func(t *T) {
						nums = append(nums, 13)
						mockAssert.Assert(nums)
					})

					it("it 4", func(t *T) {
						nums = append(nums, 14)
						mockAssert.Assert(nums)
					})
				})
			})
		})
	}()

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
		`  describe 2`,
		`    ✔ it 1`,
		`    ✔ it 2`,
		``,
		`describe 3`,
		`  describe 4`,
		`    ✔ it 3`,
		`    ✔ it 4`,
		``,
		``,
	}, "\n"), out.String())
}

func TestSpecSuitesGetExecutedInParallel(t *testing.T) {
	t.Parallel()

	var (
		out         bytes.Buffer
		testingMock = &mock{t: t}
		mockAssert  = &assertMock{}
		done        = make(chan bool, 1)
	)

	t.Run("run parallel tests", func(t *testing.T) {
		WithSpecSuite(t, func(s *SpecSuite) {
			s.t = testingMock
			describe, beforeEach, it := s.With(Output(&out)).ParallelAPI(func() { close(done) })

			describe("Checkout 1", func() {
				describe("given 1", func() {
					beforeEach(func(t *T, w *World) {
						w.Set("nums", []int{1})
					})

					describe("when 1", func() {
						beforeEach(func(t *T, w *World) {
							w.Swap("nums", func(current any) any {
								return append(current.([]int), 2)
							})
							w.Swap("nums", func(current any) any {
								return append(current.([]int), 3)
							})
						})

						describe("scenario 1", func() {
							describe("given 2", func() {
								beforeEach(func(t *T, w *World) {
									w.Swap("nums", func(current any) any {
										return append(current.([]int), 4)
									})
								})

								describe("when 2", func() {
									beforeEach(func(t *T, w *World) {
										w.Swap("nums", func(current any) any {
											return append(current.([]int), 5)
										})
									})

									it("then 2", func(t *T, w *World) {
										w.Swap("nums", func(current any) any {
											return append(current.([]int), 6)
										})
										mockAssert.Assert(w.Get("nums"))
									})
								})
							})
						})

						describe("scenario 2", func() {
							describe("given 2", func() {
								beforeEach(func(t *T, w *World) {
									w.Swap("nums", func(current any) any {
										return append(current.([]int), 7)
									})
								})

								describe("when 2", func() {
									beforeEach(func(t *T, w *World) {
										w.Swap("nums", func(current any) any {
											return append(current.([]int), 8)
										})
									})

									it("then 2", func(t *T, w *World) {
										w.Swap("nums", func(current any) any {
											return append(current.([]int), 9)
										})
										mockAssert.Assert(w.Get("nums"))
									})
								})
							})
						})
					})
				})
			})

			describe("Checkout 2", func() {
				describe("given 11", func() {
					beforeEach(func(t *T, w *World) {
						w.Set("nums", []int{11})
					})

					describe("when 11", func() {
						beforeEach(func(t *T, w *World) {
							w.Swap("nums", func(current any) any {
								return append(current.([]int), 12)
							})
							w.Swap("nums", func(current any) any {
								return append(current.([]int), 13)
							})
						})

						describe("scenario 11", func() {
							describe("given 12", func() {
								beforeEach(func(t *T, w *World) {
									w.Swap("nums", func(current any) any {
										return append(current.([]int), 14)
									})
								})

								describe("when 12", func() {
									beforeEach(func(t *T, w *World) {
										w.Swap("nums", func(current any) any {
											return append(current.([]int), 15)
										})
									})

									it("then 12", func(t *T, w *World) {
										w.Swap("nums", func(current any) any {
											return append(current.([]int), 16)
										})
										mockAssert.Assert(w.Get("nums"))
									})
								})
							})
						})

						describe("scenario 12", func() {
							describe("given 12", func() {
								beforeEach(func(t *T, w *World) {
									w.Swap("nums", func(current any) any {
										return append(current.([]int), 17)
									})
								})

								describe("when 12", func() {
									beforeEach(func(t *T, w *World) {
										w.Swap("nums", func(current any) any {
											return append(current.([]int), 18)
										})
									})

									it("then 12", func(t *T, w *World) {
										w.Swap("nums", func(current any) any {
											return append(current.([]int), 19)
										})
										mockAssert.Assert(w.Get("nums"))
									})
								})
							})
						})
					})
				})
			})
		})
	})

	t.Run("assert parallel tests run correctly", func(t *testing.T) {
		t.Parallel()

		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Errorf("test timed out")
		}

		assert.Equal(t, [][]any(nil), testingMock.calls)
		assert.Equal(t, 4, len(mockAssert.calls))

		// sort the calls since they will come out of order because of the tests executing in parallel
		sort.Slice(mockAssert.calls, func(i, j int) bool {
			if len(mockAssert.calls[i][0].([]int)) == len(mockAssert.calls[j][0].([]int)) {
				for index := 0; index < len(mockAssert.calls[i][0].([]int)); index++ {
					if (mockAssert.calls[i][0].([]int))[index] == (mockAssert.calls[j][0].([]int))[index] {
						continue
					}
					return (mockAssert.calls[i][0].([]int))[index] < (mockAssert.calls[j][0].([]int))[index]
				}
			}
			return len(mockAssert.calls[i]) < len(mockAssert.calls[j])
		})

		assert.Equal(t, []any{[]int{1, 2, 3, 4, 5, 6}}, mockAssert.calls[0])
		assert.Equal(t, []any{[]int{1, 2, 3, 7, 8, 9}}, mockAssert.calls[1])
		assert.Equal(t, []any{[]int{11, 12, 13, 14, 15, 16}}, mockAssert.calls[2])
		assert.Equal(t, []any{[]int{11, 12, 13, 17, 18, 19}}, mockAssert.calls[3])

		assert.Equal(t, []string{
			"Checkout 1/given 1/when 1/scenario 1/given 2/when 2/then 2",
			"Checkout 1/given 1/when 1/scenario 2/given 2/when 2/then 2",
			"Checkout 2/given 11/when 11/scenario 11/given 12/when 12/then 12",
			"Checkout 2/given 11/when 11/scenario 12/given 12/when 12/then 12",
		}, testingMock.testTitles)

		assert.Equal(t, strings.Join([]string{
			`Checkout 1`,
			`  given 1`,
			`    when 1`,
			`      scenario 1`,
			`        given 2`,
			`          when 2`,
			`            ✔ then 2`,
			`      scenario 2`,
			`        given 2`,
			`          when 2`,
			`            ✔ then 2`,
			``,
			`Checkout 2`,
			`  given 11`,
			`    when 11`,
			`      scenario 11`,
			`        given 12`,
			`          when 12`,
			`            ✔ then 12`,
			`      scenario 12`,
			`        given 12`,
			`          when 12`,
			`            ✔ then 12`,
			``,
			``,
		}, "\n"), out.String())
	})
}

func TestFailedTestSuitesWithFirstSuiteFailing(t *testing.T) {
	var (
		out  bytes.Buffer
		spec *SpecSuite
		tm   = &mock{t: t}
	)

	func() {
		WithSpecSuite(t, func(s *SpecSuite) {
			spec = s
			s.t = tm
			describe, _, it := s.With(Output(&out)).API()

			describe("describe 1", func() {
				it("it 1", func(t *T) {
					t.Skip()
				})
				it("it 2", func(t *T) {})
			})

			describe("describe 2", func() {
				it("it 3", func(t *T) {})
				it("it 4", func(t *T) {})
			})
		})
	}()

	assert.Equal(t, 4, len(tm.childMocks))
	assert.Equal(t, true, tm.childMocks[0].t.Skipped())
	assert.Equal(t, false, tm.childMocks[1].t.Skipped())
	assert.Equal(t, false, tm.childMocks[2].t.Skipped())
	assert.Equal(t, false, tm.childMocks[3].t.Skipped())
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
		`  s it 1`,
		`  ✔ it 2`,
		``,
		`describe 2`,
		`  ✔ it 3`,
		`  ✔ it 4`,
		``,
		``,
	}, "\n"), out.String())
}

func TestFailedTestSuitesWithLastSuiteFailing(t *testing.T) {
	var (
		out  bytes.Buffer
		spec *SpecSuite
		tm   = &mock{t: t}
	)

	func() {
		WithSpecSuite(t, func(s *SpecSuite) {
			spec = s
			s.t = tm
			describe, _, it := s.With(Output(&out)).API()

			describe("describe 1", func() {
				it("it 1", func(t *T) {})
				it("it 2", func(t *T) {})
			})

			describe("describe 2", func() {
				it("it 3", func(t *T) {})
				it("it 4", func(t *T) {
					t.Skip()
				})
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
		`describe 1`,
		`  ✔ it 1`,
		`  ✔ it 2`,
		``,
		`describe 2`,
		`  ✔ it 3`,
		`  s it 4`,
		``,
		``,
	}, "\n"), out.String())
}
