package gospec

import (
	"bytes"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/slavsan/gospec/internal/testing/helpers/assert"
)

func TestFeaturesCanBeSetAtTopLevel(t *testing.T) {
	var (
		out     bytes.Buffer
		tm      = &mock{t: nil}
		fs      = NewFeatureSuite(t, WithOutput(&out))
		feature = fs.Feature
	)

	fs.t = tm

	feature("feature 1", func() {})
	feature("feature 2", func() {})

	assert.Equal(t, [][]any(nil), tm.calls)

	assert.Equal(t, strings.Join([]string{
		`Feature: feature 1`,
		``,
		`Feature: feature 2`,
		``,
		``,
	}, "\n"), out.String())
}

func TestFeaturesCanNotBeNested(t *testing.T) {
	var (
		out         bytes.Buffer
		testingMock = &mock{t: t}
		fs          = NewFeatureSuite(t, WithOutput(&out))
		feature     = fs.Feature
	)

	fs.t = testingMock

	feature("Checkout", func() {
		feature("nested feature", func() {
		})
	})

	assert.Equal(t, [][]any{{"invalid position for `Feature` function, it must be at top level"}}, testingMock.calls)
	assert.Equal(t, "\n", out.String())
}

func TestFeaturesContainOnlyScenariosAndBackgroundCalls(t *testing.T) {
	var (
		out         bytes.Buffer
		testingMock = &mock{t: t}
		fs          = NewFeatureSuite(t, WithOutput(&out))
		feature     = fs.Feature
		background  = fs.Background
		scenario    = fs.Scenario
	)

	fs.t = testingMock

	feature("Checkout", func() {
		background(func() {})
		scenario("scenario 1", func() {})
		scenario("scenario 2", func() {})
	})

	assert.Equal(t, [][]any(nil), testingMock.calls)
	assert.Equal(t, []string{
		"Checkout/scenario 1",
		"Checkout/scenario 2",
	}, testingMock.testTitles)

	assert.Equal(t, strings.Join([]string{
		`Feature: Checkout`,
		``,
		`	Background:`,
		``,
		`	Scenario: scenario 1`,
		``,
		`	Scenario: scenario 2`,
		``,
		``,
	}, "\n"), out.String())
}

func TestScenariosCanNotBeNested(t *testing.T) {
	var (
		out         bytes.Buffer
		testingMock = &mock{t: t}
		spec        = NewFeatureSuite(t, WithOutput(&out))
		feature     = spec.Feature
		scenario    = spec.Scenario
	)

	spec.t = testingMock

	feature("Checkout", func() {
		scenario("scenario 1", func() {
			scenario("scenario 2", func() {})
		})
	})

	assert.Equal(t, [][]any{{"invalid position for `Scenario` function, it must be inside a `Feature` call"}}, testingMock.calls)
	assert.Equal(t, []string{"Checkout/scenario 1"}, testingMock.testTitles)

	assert.Equal(t, strings.Join([]string{
		`Feature: Checkout`,
		``,
		`	Scenario: scenario 1`,
		``,
		``,
	}, "\n"), out.String())
}

func TestScenarioCanContainGivenWhenThen(t *testing.T) {
	var (
		out         bytes.Buffer
		testingMock = &mock{t: t}
		spec        = NewFeatureSuite(t, WithOutput(&out))
		feature     = spec.Feature
		scenario    = spec.Scenario
		given       = spec.Given
		when        = spec.When
		then        = spec.Then
	)

	spec.t = testingMock

	feature("Checkout", func() {
		scenario("scenario 1", func() {
			given("given 1", func() {})
			when("when 1", func() {})
			then("then 1", func() {})
		})
	})

	assert.Equal(t, [][]any(nil), testingMock.calls)
	assert.Equal(t, 0, len(spec.stack))
	assert.Equal(t, 1, len(spec.suites))
	assert.Equal(t, 5, len(spec.suites[0]))
	assert.Equal(t, "Checkout", spec.suites[0][0].title)
	assert.Equal(t, "scenario 1", spec.suites[0][1].title)
	assert.Equal(t, "given 1", spec.suites[0][2].title)
	assert.Equal(t, "when 1", spec.suites[0][3].title)
	assert.Equal(t, "then 1", spec.suites[0][4].title)

	assert.Equal(t, []string{
		"Checkout/scenario 1",
	}, testingMock.testTitles)

	assert.Equal(t, strings.Join([]string{
		`Feature: Checkout`,
		``,
		`	Scenario: scenario 1`,
		`		Given given 1`,
		`		When when 1`,
		`		Then then 1`,
		``,
		``,
	}, "\n"), out.String())
}

func TestMultipleScenarioWithGivenWhenThen(t *testing.T) {
	var (
		out         bytes.Buffer
		testingMock = &mock{t: t}
		spec        = NewFeatureSuite(t, WithOutput(&out))
		feature     = spec.Feature
		scenario    = spec.Scenario
		given       = spec.Given
		when        = spec.When
		then        = spec.Then
	)

	spec.t = testingMock

	feature("Checkout", func() {
		scenario("scenario 1", func() {
			given("given 1", func() {})
			when("when 1", func() {})
			then("then 1", func() {})
		})

		scenario("scenario 2", func() {
			given("given 2", func() {})
			when("when 2", func() {})
			then("then 2", func() {})
		})
	})

	assert.Equal(t, [][]any(nil), testingMock.calls)
	assert.Equal(t, 0, len(spec.stack))
	assert.Equal(t, 2, len(spec.suites))
	assert.Equal(t, 5, len(spec.suites[0]))
	assert.Equal(t, "Checkout", spec.suites[0][0].title)
	assert.Equal(t, "scenario 1", spec.suites[0][1].title)
	assert.Equal(t, "given 1", spec.suites[0][2].title)
	assert.Equal(t, "when 1", spec.suites[0][3].title)
	assert.Equal(t, "then 1", spec.suites[0][4].title)
	assert.Equal(t, 5, len(spec.suites[1]))
	assert.Equal(t, "Checkout", spec.suites[1][0].title)
	assert.Equal(t, "scenario 2", spec.suites[1][1].title)
	assert.Equal(t, "given 2", spec.suites[1][2].title)
	assert.Equal(t, "when 2", spec.suites[1][3].title)
	assert.Equal(t, "then 2", spec.suites[1][4].title)
	assert.Equal(t, []string{
		"Checkout/scenario 1",
		"Checkout/scenario 2",
	}, testingMock.testTitles)
	assert.Equal(t, strings.Join([]string{
		`Feature: Checkout`,
		``,
		`	Scenario: scenario 1`,
		`		Given given 1`,
		`		When when 1`,
		`		Then then 1`,
		``,
		`	Scenario: scenario 2`,
		`		Given given 2`,
		`		When when 2`,
		`		Then then 2`,
		``,
		``,
	}, "\n"), out.String())
}

func TestScenarioWhichHasBackgroundBlock(t *testing.T) {
	var (
		out         bytes.Buffer
		testingMock = &mock{t: t}
		spec        = NewFeatureSuite(t, WithOutput(&out))
		feature     = spec.Feature
		scenario    = spec.Scenario
		background  = spec.Background
		given       = spec.Given
		when        = spec.When
		then        = spec.Then
	)

	spec.t = testingMock

	feature("Checkout", func() {
		background(func() {
			given("given 0", func() {})
			when("when 0", func() {})
			then("then 0", func() {})
		})

		scenario("scenario 1", func() {
			given("given 1", func() {})
			when("when 1", func() {})
			then("then 1", func() {})
		})
	})

	assert.Equal(t, [][]any(nil), testingMock.calls)
	assert.Equal(t, 0, len(spec.stack))
	assert.Equal(t, 1, len(spec.suites))
	assert.Equal(t, 9, len(spec.suites[0]))
	assert.Equal(t, "Checkout", spec.suites[0][0].title)
	assert.Equal(t, "", spec.suites[0][1].title)
	assert.Equal(t, "given 0", spec.suites[0][2].title)
	assert.Equal(t, "when 0", spec.suites[0][3].title)
	assert.Equal(t, "then 0", spec.suites[0][4].title)
	assert.Equal(t, "scenario 1", spec.suites[0][5].title)
	assert.Equal(t, "given 1", spec.suites[0][6].title)
	assert.Equal(t, "when 1", spec.suites[0][7].title)
	assert.Equal(t, "then 1", spec.suites[0][8].title)
	assert.Equal(t, []string{"Checkout/scenario 1"}, testingMock.testTitles)
	assert.Equal(t, strings.Join([]string{
		`Feature: Checkout`,
		``,
		`	Background:`,
		`		Given given 0`,
		`		When when 0`,
		`		Then then 0`,
		``,
		`	Scenario: scenario 1`,
		`		Given given 1`,
		`		When when 1`,
		`		Then then 1`,
		``,
		``,
	}, "\n"), out.String())
}

func TestMultipleScenariosWhichShareTheSameBackgroundBlock(t *testing.T) {
	var (
		out         bytes.Buffer
		testingMock = &mock{t: t}
		spec        = NewFeatureSuite(t, WithOutput(&out))
		feature     = spec.Feature
		scenario    = spec.Scenario
		background  = spec.Background
		given       = spec.Given
		when        = spec.When
		then        = spec.Then
	)

	spec.t = testingMock

	feature("Checkout", func() {
		background(func() {
			given("given 0", func() {})
			when("when 0", func() {})
			then("then 0", func() {})
		})

		scenario("scenario 1", func() {
			given("given 1", func() {})
			when("when 1", func() {})
			then("then 1", func() {})
		})

		scenario("scenario 2", func() {
			given("given 2", func() {})
			when("when 2", func() {})
			then("then 2", func() {})
		})
	})

	assert.Equal(t, [][]any(nil), testingMock.calls)
	assert.Equal(t, 0, len(spec.stack))
	assert.Equal(t, 2, len(spec.suites))

	firstSuite := spec.suites[0]
	assert.Equal(t, 9, len(firstSuite))
	assert.Equal(t, "Checkout", firstSuite[0].title)
	assert.Equal(t, "", firstSuite[1].title)
	assert.Equal(t, "given 0", firstSuite[2].title)
	assert.Equal(t, "when 0", firstSuite[3].title)
	assert.Equal(t, "then 0", firstSuite[4].title)
	assert.Equal(t, "scenario 1", firstSuite[5].title)
	assert.Equal(t, "given 1", firstSuite[6].title)
	assert.Equal(t, "when 1", firstSuite[7].title)
	assert.Equal(t, "then 1", firstSuite[8].title)

	secondSuite := spec.suites[1]
	assert.Equal(t, 9, len(secondSuite))
	assert.Equal(t, "Checkout", secondSuite[0].title)
	assert.Equal(t, "", secondSuite[1].title)
	assert.Equal(t, "given 0", secondSuite[2].title)
	assert.Equal(t, "when 0", secondSuite[3].title)
	assert.Equal(t, "then 0", secondSuite[4].title)
	assert.Equal(t, "scenario 2", secondSuite[5].title)
	assert.Equal(t, "given 2", secondSuite[6].title)
	assert.Equal(t, "when 2", secondSuite[7].title)
	assert.Equal(t, "then 2", secondSuite[8].title)

	assert.Equal(t, []string{
		"Checkout/scenario 1",
		"Checkout/scenario 2",
	}, testingMock.testTitles)

	assert.Equal(t, strings.Join([]string{
		`Feature: Checkout`,
		``,
		`	Background:`,
		`		Given given 0`,
		`		When when 0`,
		`		Then then 0`,
		``,
		`	Scenario: scenario 1`,
		`		Given given 1`,
		`		When when 1`,
		`		Then then 1`,
		``,
		`	Scenario: scenario 2`,
		`		Given given 2`,
		`		When when 2`,
		`		Then then 2`,
		``,
		``,
	}, "\n"), out.String())
}

func TestFeaturesGetExecutedInCorrectOrder(t *testing.T) {
	var (
		out         bytes.Buffer
		testingMock = &mock{t: t}
		mockAssert  = &assertMock{}
		spec        = NewFeatureSuite(t, WithOutput(&out))
		feature, background, scenario,
		given, when, then, table = spec.API()
	)

	spec.t = testingMock

	_ = table

	feature("Checkout 1", func() {
		var nums []int

		background(func() {
			given("given 1", func() {
				nums = []int{}
				nums = append(nums, 1)
			})
			when("when 1", func() {
				nums = append(nums, 2)
			})
			then("then 1", func() {
				nums = append(nums, 3)
			})
		})

		scenario("scenario 1", func() {
			given("given 2", func() {
				nums = append(nums, 4)
			})
			when("when 2", func() {
				nums = append(nums, 5)
			})
			then("then 2", func() {
				nums = append(nums, 6)
				mockAssert.Assert(nums)
			})
		})

		scenario("scenario 2", func() {
			given("given 2", func() {
				nums = append(nums, 7)
			})
			when("when 2", func() {
				nums = append(nums, 8)
			})
			then("then 2", func() {
				nums = append(nums, 9)
				mockAssert.Assert(nums)
			})
		})
	})

	feature("Checkout 2", func() {
		var nums []int

		background(func() {
			given("given 11", func() {
				nums = []int{}
				nums = append(nums, 11)
			})
			when("when 11", func() {
				nums = append(nums, 12)
			})
			then("then 11", func() {
				nums = append(nums, 13)
			})
		})

		scenario("scenario 11", func() {
			given("given 12", func() {
				nums = append(nums, 14)
			})
			when("when 12", func() {
				nums = append(nums, 15)
			})
			then("then 12", func() {
				nums = append(nums, 16)
				mockAssert.Assert(nums)
			})
		})

		scenario("scenario 12", func() {
			given("given 12", func() {
				nums = append(nums, 17)
			})
			when("when 12", func() {
				nums = append(nums, 18)
			})
			then("then 12", func() {
				nums = append(nums, 19)
				mockAssert.Assert(nums)
			})
		})
	})

	assert.Equal(t, [][]any(nil), testingMock.calls)
	assert.Equal(t, 4, len(mockAssert.calls))
	assert.Equal(t, []any{[]int{1, 2, 3, 4, 5, 6}}, mockAssert.calls[0])
	assert.Equal(t, []any{[]int{1, 2, 3, 7, 8, 9}}, mockAssert.calls[1])
	assert.Equal(t, []any{[]int{11, 12, 13, 14, 15, 16}}, mockAssert.calls[2])
	assert.Equal(t, []any{[]int{11, 12, 13, 17, 18, 19}}, mockAssert.calls[3])

	assert.Equal(t, []string{
		"Checkout 1/scenario 1",
		"Checkout 1/scenario 2",
		"Checkout 2/scenario 11",
		"Checkout 2/scenario 12",
	}, testingMock.testTitles)

	assert.Equal(t, strings.Join([]string{
		`Feature: Checkout 1`,
		``,
		`	Background:`,
		`		Given given 1`,
		`		When when 1`,
		`		Then then 1`,
		``,
		`	Scenario: scenario 1`,
		`		Given given 2`,
		`		When when 2`,
		`		Then then 2`,
		``,
		`	Scenario: scenario 2`,
		`		Given given 2`,
		`		When when 2`,
		`		Then then 2`,
		``,
		`Feature: Checkout 2`,
		``,
		`	Background:`,
		`		Given given 11`,
		`		When when 11`,
		`		Then then 11`,
		``,
		`	Scenario: scenario 11`,
		`		Given given 12`,
		`		When when 12`,
		`		Then then 12`,
		``,
		`	Scenario: scenario 12`,
		`		Given given 12`,
		`		When when 12`,
		`		Then then 12`,
		``,
		``,
	}, "\n"), out.String())
}

func TestFeaturesGetExecutedInParallel(t *testing.T) {
	t.Parallel()

	var (
		out         bytes.Buffer
		testingMock = &mock{t: t}
		mockAssert  = &assertMock{}
		done        = make(chan bool, 1)
		spec        = NewFeatureSuite(t, WithOutput(&out), WithParallel(func() { close(done) }))
		feature, background, scenario,
		given, when, then, table = spec.API()
	)

	spec.t = testingMock

	var wg sync.WaitGroup
	wg.Add(4)

	_ = background
	_ = table

	t.Run("run parallel tests", func(t *testing.T) {
		feature("Checkout 1", func() {
			background(func() {
				given("given 1", func(w *World) {
					w.Set("nums", []int{1})
				})
				when("when 1", func(w *World) {
					w.Swap("nums", func(current any) any {
						return append(current.([]int), 2)
					})
				})
				then("then 1", func(w *World) {
					w.Swap("nums", func(current any) any {
						return append(current.([]int), 3)
					})
				})
			})

			scenario("scenario 1", func() {
				given("given 2", func(w *World) {
					w.Swap("nums", func(current any) any {
						return append(current.([]int), 4)
					})
				})
				when("when 2", func(w *World) {
					w.Swap("nums", func(current any) any {
						return append(current.([]int), 5)
					})
				})
				then("then 2", func(w *World) {
					defer wg.Done()
					w.Swap("nums", func(current any) any {
						return append(current.([]int), 6)
					})
					mockAssert.Assert(w.Get("nums"))
				})
			})

			scenario("scenario 2", func() {
				given("given 2", func(w *World) {
					w.Swap("nums", func(current any) any {
						return append(current.([]int), 7)
					})
				})
				when("when 2", func(w *World) {
					w.Swap("nums", func(current any) any {
						return append(current.([]int), 8)
					})
				})
				then("then 2", func(w *World) {
					defer wg.Done()
					w.Swap("nums", func(current any) any {
						return append(current.([]int), 9)
					})
					mockAssert.Assert(w.Get("nums"))
				})
			})
		})

		feature("Checkout 2", func() {
			background(func() {
				given("given 11", func(w *World) {
					w.Set("nums", []int{11})
				})
				when("when 11", func(w *World) {
					w.Swap("nums", func(current any) any {
						return append(current.([]int), 12)
					})
				})
				then("then 11", func(w *World) {
					w.Swap("nums", func(current any) any {
						return append(current.([]int), 13)
					})
				})
			})

			scenario("scenario 11", func() {
				given("given 12", func(w *World) {
					w.Swap("nums", func(current any) any {
						return append(current.([]int), 14)
					})
				})
				when("when 12", func(w *World) {
					w.Swap("nums", func(current any) any {
						return append(current.([]int), 15)
					})
				})
				then("then 12", func(w *World) {
					defer wg.Done()
					w.Swap("nums", func(current any) any {
						return append(current.([]int), 16)
					})
					mockAssert.Assert(w.Get("nums"))
				})
			})

			scenario("scenario 12", func() {
				given("given 12", func(w *World) {
					w.Swap("nums", func(current any) any {
						return append(current.([]int), 17)
					})
				})
				when("when 12", func(w *World) {
					w.Swap("nums", func(current any) any {
						return append(current.([]int), 18)
					})
				})
				then("then 12", func(w *World) {
					defer wg.Done()
					w.Swap("nums", func(current any) any {
						return append(current.([]int), 19)
					})
					mockAssert.Assert(w.Get("nums"))
				})
			})
		})
	})

	t.Run("assert parallel tests run correctly", func(t *testing.T) {
		t.Parallel()

		go func() {
			wg.Wait()
			close(done) // TODO: remove this when the API gets updated
		}()

		select {
		case <-done:
		case <-time.After(2 * time.Second):
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
			"Checkout 1/scenario 1",
			"Checkout 1/scenario 2",
			"Checkout 2/scenario 11",
			"Checkout 2/scenario 12",
		}, testingMock.testTitles)

		assert.Equal(t, strings.Join([]string{
			`Feature: Checkout 1`,
			``,
			`	Background:`,
			`		Given given 1`,
			`		When when 1`,
			`		Then then 1`,
			``,
			`	Scenario: scenario 1`,
			`		Given given 2`,
			`		When when 2`,
			`		Then then 2`,
			``,
			`	Scenario: scenario 2`,
			`		Given given 2`,
			`		When when 2`,
			`		Then then 2`,
			``,
			`Feature: Checkout 2`,
			``,
			`	Background:`,
			`		Given given 11`,
			`		When when 11`,
			`		Then then 11`,
			``,
			`	Scenario: scenario 11`,
			`		Given given 12`,
			`		When when 12`,
			`		Then then 12`,
			``,
			`	Scenario: scenario 12`,
			`		Given given 12`,
			`		When when 12`,
			`		Then then 12`,
			``,
			``,
		}, "\n"), out.String())
	})
}

func TestTableOutput(t *testing.T) {
	var (
		out         bytes.Buffer
		testingMock = &mock{t: t}
		spec        = NewFeatureSuite(t, WithOutput(&out))
		feature     = spec.Feature
		scenario    = spec.Scenario
		given       = spec.Given
		when        = spec.When
		then        = spec.Then
		table       = spec.Table
	)

	spec.t = testingMock

	feature("Checkout", func() {
		type Product struct {
			Name  string
			Price float64
			Type  int
		}

		var items []Product

		scenario("scenario 1", func() {
			given("given 1", func() {
				// ..
			})
			when("when 1", func() {})
			then("then 1", func() {
				items = []Product{
					{Name: "Gopher toy", Price: 14.99, Type: 2},
					{Name: "Crab toy", Price: 17.49, Type: 8},
				}
				table([]string{"Name", "Price"}, items)
			})
		})
	})

	assert.Equal(t, [][]any(nil), testingMock.calls)
	assert.Equal(t, 0, len(spec.stack))
	assert.Equal(t, 1, len(spec.suites))
	assert.Equal(t, 5, len(spec.suites[0]))
	assert.Equal(t, "Checkout", spec.suites[0][0].title)
	assert.Equal(t, "scenario 1", spec.suites[0][1].title)
	assert.Equal(t, "given 1", spec.suites[0][2].title)
	assert.Equal(t, "when 1", spec.suites[0][3].title)
	assert.Equal(t, "then 1", spec.suites[0][4].title)

	assert.Equal(t, []string{
		"Checkout/scenario 1",
	}, testingMock.testTitles)

	assert.Equal(t, strings.Join([]string{
		`Feature: Checkout`,
		``,
		`	Scenario: scenario 1`,
		`		Given given 1`,
		`		When when 1`,
		`		Then then 1`,
		`			| Name       | Price |`,
		`			| Gopher toy | 14.99 |`,
		`			| Crab toy   | 17.49 |`,
		``,
		``,
	}, "\n"), out.String())
}
