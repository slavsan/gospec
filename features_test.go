package gospec

import (
	"testing"

	"github.com/slavsan/gospec/internal/testing/helpers/assert"
)

type mock struct {
	t     *testing.T
	calls [][]any
	//
}

func (m *mock) Helper() {
	// ..
}

func (m *mock) Errorf(format string, args ...interface{}) {
	var call []any
	call = append(call, format)
	call = append(call, args...)
	m.calls = append(m.calls, call)
}

func (m *mock) Run(name string, f func(t *testing.T)) bool {
	f(m.t)
	return false
}

type assertMock struct {
	calls [][]any
}

func (m *assertMock) Assert(args ...any) {
	m.calls = append(m.calls, args)
	// ..
}

func TestFeaturesCanBeSetAtTopLevel(t *testing.T) {
	var (
		tm      = &mock{}
		fs      = NewFeatureSuite(tm)
		feature = fs.Feature
	)

	feature("feature 1", func() {})
	feature("feature 2", func() {})

	assert.Equal(t, [][]any(nil), tm.calls)
}

func TestFeaturesCanNotBeNested(t *testing.T) {
	var (
		testingMock = &mock{}
		fs          = NewFeatureSuite(testingMock)
		feature     = fs.Feature
	)

	feature("Checkout", func() {
		feature("nested feature", func() {
		})
	})

	assert.Equal(t, [][]any{{"invalid position for `Feature` function, it must be at top level"}}, testingMock.calls)
}

func TestFeaturesContainOnlyScenariosAndBackgroundCalls(t *testing.T) {
	var (
		testingMock = &mock{}
		fs          = NewFeatureSuite(testingMock)
		feature     = fs.Feature
		background  = fs.Background
		scenario    = fs.Scenario
	)

	feature("Checkout", func() {
		background("background", func() {})
		scenario("scenario 1", func() {})
		scenario("scenario 2", func() {})
	})

	assert.Equal(t, [][]any(nil), testingMock.calls)
}

func TestScenariosCanNotBeNested(t *testing.T) {
	var (
		testingMock = &mock{}
		spec        = NewFeatureSuite(testingMock)
		feature     = spec.Feature
		scenario    = spec.Scenario
	)

	feature("Checkout", func() {
		scenario("scenario 1", func() {
			scenario("scenario 2", func() {})
		})
	})

	assert.Equal(t, [][]any{{"invalid position for `Scenario` function, it must be inside a `Feature` call"}}, testingMock.calls)
}

func TestScenarioCanContainGivenWhenThen(t *testing.T) {
	var (
		testingMock = &mock{}
		spec        = NewFeatureSuite(testingMock)
		feature     = spec.Feature
		scenario    = spec.Scenario
		given       = spec.Given
		when        = spec.When
		then        = spec.Then
	)

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
}

func TestMultipleScenarioWithGivenWhenThen(t *testing.T) {
	var (
		testingMock = &mock{}
		spec        = NewFeatureSuite(testingMock)
		feature     = spec.Feature
		scenario    = spec.Scenario
		given       = spec.Given
		when        = spec.When
		then        = spec.Then
	)

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
}

func TestScenarioWhichHasBackgroundBlock(t *testing.T) {
	var (
		testingMock = &mock{}
		spec        = NewFeatureSuite(testingMock)
		feature     = spec.Feature
		scenario    = spec.Scenario
		background  = spec.Background
		given       = spec.Given
		when        = spec.When
		then        = spec.Then
	)

	feature("Checkout", func() {
		background("background 1", func() {
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
	assert.Equal(t, "background 1", spec.suites[0][1].title)
	assert.Equal(t, "given 0", spec.suites[0][2].title)
	assert.Equal(t, "when 0", spec.suites[0][3].title)
	assert.Equal(t, "then 0", spec.suites[0][4].title)
	assert.Equal(t, "scenario 1", spec.suites[0][5].title)
	assert.Equal(t, "given 1", spec.suites[0][6].title)
	assert.Equal(t, "when 1", spec.suites[0][7].title)
	assert.Equal(t, "then 1", spec.suites[0][8].title)
}

func TestMultipleScenariosWhichShareTheSameBackgroundBlock(t *testing.T) {
	var (
		testingMock = &mock{}
		spec        = NewFeatureSuite(testingMock)
		feature     = spec.Feature
		scenario    = spec.Scenario
		background  = spec.Background
		given       = spec.Given
		when        = spec.When
		then        = spec.Then
	)

	feature("Checkout", func() {
		background("background 1", func() {
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
	assert.Equal(t, "background 1", firstSuite[1].title)
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
	assert.Equal(t, "background 1", secondSuite[1].title)
	assert.Equal(t, "given 0", secondSuite[2].title)
	assert.Equal(t, "when 0", secondSuite[3].title)
	assert.Equal(t, "then 0", secondSuite[4].title)
	assert.Equal(t, "scenario 2", secondSuite[5].title)
	assert.Equal(t, "given 2", secondSuite[6].title)
	assert.Equal(t, "when 2", secondSuite[7].title)
	assert.Equal(t, "then 2", secondSuite[8].title)
}

func TestFeaturesGetExecutedInCorrectOrder(t *testing.T) {
	var (
		testingMock = &mock{t: t}
		mockAssert  = &assertMock{}
		spec        = NewFeatureSuite(testingMock)
		feature, background, scenario,
		given, when, then, world, table = spec.API()
	)

	_ = world
	_ = table

	feature("Checkout 1", func() {
		var nums []int

		background("background 1", func() {
			given("given 1", func() {
				nums = []int{}
				nums = append(nums, 1)
			})
			given("when 1", func() {
				nums = append(nums, 2)
			})
			given("then 1", func() {
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

		scenario("scenario 1", func() {
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

	feature("Checkout 11", func() {
		var nums []int

		background("background 11", func() {
			given("given 11", func() {
				nums = []int{}
				nums = append(nums, 11)
			})
			given("when 11", func() {
				nums = append(nums, 12)
			})
			given("then 11", func() {
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

		scenario("scenario 11", func() {
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
}

//func TestFeature(t *testing.T) {
//	// ..
//
//	f := gospec.NewFeatureSuite(t)
//	feature, background, scenario, given, when, then, world, table, start := f.API()
//	defer start()
//
//	_ = world
//
//	feature("Checkout", func() {
//		type Product struct {
//			Name  string
//			Price float64
//			Type  int
//			// ..
//		}
//
//		var (
//			customer       string
//			cart           []string
//			appliedCoupons []string
//			items          []Product
//		)
//
//		_ = customer
//		_ = cart
//		_ = appliedCoupons
//
//		background("", func() {
//			given("there is a customer account", func() {
//				customer = "John Doe"
//				// world.Define(func() {
//				// 	world.Var("foo", "bar")
//				// 	world.Var("spam", "eggs")
//				// })
//			})
//			given("the cart has 2 items", func() {
//				items = []Product{
//					{Name: "Gopher toy", Price: 14.99, Type: 2},
//					{Name: "Crab toy", Price: 17.49, Type: 8},
//				}
//				table([]string{"Name", "Price"}, items)
//				// ..
//			})
//			given("no coupons are applied", func() {
//				appliedCoupons = []string{}
//			})
//			// ..
//		})
//
//		scenario("with no coupons", func() {
//			given("has the correct total price of both items", func() {
//				cart = []string{}
//				// ..
//			})
//			when("the customer checks out", func() {
//				cart = []string{}
//				cart = append(cart, "Gopher toy")
//				// ..
//			})
//			then("the cart should be 0 again", func() {
//				// ..
//				assert.Equal(t, []string{"Gopher toy"}, cart)
//			})
//			// ...
//		})
//
//		scenario("with coupon applied", func() {
//			given("has the correct total price of both items", func() {
//				cart = []string{}
//				// ..
//			})
//			when("the customer checks out", func() {
//				cart = []string{}
//				// ..
//			})
//			then("the cart should be 0 again", func() {
//				// ..
//				assert.Equal(t, []string{}, cart)
//			})
//			// ...
//		})
//	})
//
//	// debugFeature(f)
//}
