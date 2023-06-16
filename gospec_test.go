package gospec

import (
	"fmt"
	"testing"
)

func TestExampleSuite(t *testing.T) {
	s := NewTestSuite(t)
	describe, context, it, beforeEach, expect, start :=
		s.Describe, s.Describe, s.It, s.BeforeEach, s.Expect, s.Start
	defer start()

	describe("Checkout", func() {
		context("when shopping cart has 1 item", func() {
			var cart []string
			var appliedCoupons []string

			beforeEach(func() {
				cart = []string{"Gopher toy"}
			})
			it("should have 1 item in the cart", func() {
				expect(cart).To.Have.LengthOf(1)
			})
			it("should have no coupon applied by default", func() {
				expect(appliedCoupons).To.Have.LengthOf(0)
			})
			context("when we add one more item to the cart", func() {
				beforeEach(func() {
					cart = append(cart, "Crab toy")
				})
				it("should have 2 items in the cart", func() {
					expect(cart).To.Have.LengthOf(2)
				})
				context("when the coupon is eligible for this purchase", func() {
					context("and the coupon gets applied", func() {
						context("but the coupon value is higher than the purcahse value", func() {
							// ..
						})
						context("and the coupon value is less than the purchase value", func() {
							beforeEach(func() {
								// ..
							})
							describe("when completing the purchase", func() {
								// ..
							})
						})
					})
				})
			})
		})
		context("when shopping cart is empty", func() {
			var cart []string
			it("should have 0 items", func() {
				expect(cart).To.Have.LengthOf(0)
			})
		})
	})

	describe("Sign Up", func() {
		var signedUp bool
		context("when the user signs up", func() {
			beforeEach(func() {
				signedUp = true
			})
			it("should be signed up", func() {
				expect(signedUp).To.Be.True()
			})
		})
	})
}

func TestDescribe(t *testing.T) {
	s := NewTestSuite(t)
	describe, context, it, beforeEach, expect, start :=
		s.Describe, s.Describe, s.It, s.BeforeEach, s.Expect, s.Start
	defer start()

	describe("describe block", func() {
		context("with single describe", func() {
			var s2 *Suite
			var suites [][]*step
			beforeEach(func() {
				s2 = NewTestSuite(t)
				describe2 := s2.Describe

				// unit under test
				describe2("testing describe", func() {
				})

				suites = s2.buildSuites()
			})
			it("should have 1 step defined", func() {
				expect(s2.steps).To.Have.LengthOf(1)
			})
			it("should have just one suite with just one step defined", func() {
				expect(suites).To.Have.LengthOf(1)
				expect(suites[0]).To.Have.LengthOf(1)
				expect(buildSuiteTitle(suites[0])).To.Be.EqualTo("testing describe")
			})
		})
		context("with two sibling describes", func() {
			var s2 *Suite
			var suites [][]*step
			beforeEach(func() {
				s2 = NewTestSuite(t)
				describe2 := s2.Describe

				// unit under test
				describe2("describe", func() {
				})
				describe2("sibling describe", func() {
				})

				suites = s2.buildSuites()
			})
			it("should have 2 steps defined", func() {
				expect(s2.steps).To.Have.LengthOf(2)
			})
			it("should have just two suites with just one step defined", func() {
				expect(suites).To.Have.LengthOf(2)
				expect(suites[0]).To.Have.LengthOf(1)
				expect(buildSuiteTitle(suites[0])).To.Be.EqualTo("describe")
				expect(suites[1]).To.Have.LengthOf(1)
				expect(buildSuiteTitle(suites[1])).To.Be.EqualTo("sibling describe")
			})
		})
		context("with two describes, one parent and one child", func() {
			var s2 *Suite
			var suites [][]*step
			beforeEach(func() {
				s2 = NewTestSuite(t)
				describe2 := s2.Describe

				// unit under test
				describe2("parent describe", func() {
					describe2("child describe", func() {
					})
				})

				suites = s2.buildSuites()
			})
			it("should have 2 steps defined", func() {
				expect(s2.steps).To.Have.LengthOf(2)
			})
			it("should have just one suite with just one step defined", func() {
				expect(suites).To.Have.LengthOf(1)
				expect(suites[0]).To.Have.LengthOf(2)
				expect(buildSuiteTitle(suites[0])).To.Be.EqualTo("parent describe/child describe")
			})
		})
		context("with three nested describes", func() {
			var s2 *Suite
			var suites [][]*step
			beforeEach(func() {
				s2 = NewTestSuite(t)
				describe2 := s2.Describe

				// unit under test
				describe2("top most describe", func() {
					describe2("nested describe", func() {
						describe2("most nested describe", func() {
						})
					})
				})

				suites = s2.buildSuites()
			})
			it("should have 3 steps defined", func() {
				expect(s2.steps).To.Have.LengthOf(3)
			})
			it("should have just one suite with just one step defined", func() {
				expect(suites).To.Have.LengthOf(1)
				expect(suites[0]).To.Have.LengthOf(3)
				expect(buildSuiteTitle(suites[0])).To.Be.EqualTo("top most describe/nested describe/most nested describe")
			})
		})
	})
}

func TestIt(t *testing.T) {
	s := NewTestSuite(t)
	describe, context, it, beforeEach, expect, start :=
		s.Describe, s.Describe, s.It, s.BeforeEach, s.Expect, s.Start
	defer start()

	describe("it block", func() {
		context("with single describe and single it block", func() {
			var s2 *Suite
			var suites [][]*step
			beforeEach(func() {
				s2 = NewTestSuite(t)
				describe2, it2 := s2.Describe, s2.It

				// unit under test
				describe2("testing describe", func() {
					it2("testing it", func() {
					})
				})

				suites = s2.buildSuites()
			})
			it("should have 2 steps defined", func() {
				expect(s2.steps).To.Have.LengthOf(2)
			})
			it("should have just 1 suite with just 2 steps defined", func() {
				expect(suites).To.Have.LengthOf(1)
				expect(suites[0]).To.Have.LengthOf(2)
				expect(buildSuiteTitle(suites[0])).To.Be.EqualTo("testing describe/testing it")
			})
		})
		context("with single describe and two it blocks", func() {
			var s2 *Suite
			var suites [][]*step
			beforeEach(func() {
				s2 = NewTestSuite(t)
				describe2, it2 := s2.Describe, s2.It

				// unit under test
				describe2("testing describe", func() {
					it2("testing it", func() {
					})
					it2("testing another it", func() {
					})
				})

				suites = s2.buildSuites()
			})
			it("should have 3 steps defined", func() {
				expect(s2.steps).To.Have.LengthOf(3)
			})
			it("should have 2 suites with just 2 steps defined in each", func() {
				expect(suites).To.Have.LengthOf(2)
				expect(suites[0]).To.Have.LengthOf(2)
				expect(buildSuiteTitle(suites[0])).To.Be.EqualTo("testing describe/testing it")
				expect(suites[1]).To.Have.LengthOf(2)
				expect(buildSuiteTitle(suites[1])).To.Be.EqualTo("testing describe/testing another it")
			})
		})
		context("with two describe blocks and one it block in each", func() {
			var s2 *Suite
			var suites [][]*step
			beforeEach(func() {
				s2 = NewTestSuite(t)
				describe2, it2 := s2.Describe, s2.It

				// unit under test
				describe2("testing describe", func() {
					it2("testing it", func() {
					})
				})

				// unit under test
				describe2("testing another describe", func() {
					it2("testing another it", func() {
					})
				})

				suites = s2.buildSuites()
			})
			it("should have 4 steps defined", func() {
				expect(s2.steps).To.Have.LengthOf(4)
			})
			it("should have 2 suites with just 2 steps defined in each", func() {
				expect(suites).To.Have.LengthOf(2)
				expect(suites[0]).To.Have.LengthOf(2)
				expect(buildSuiteTitle(suites[0])).To.Be.EqualTo("testing describe/testing it")
				expect(suites[1]).To.Have.LengthOf(2)
				expect(buildSuiteTitle(suites[1])).To.Be.EqualTo("testing another describe/testing another it")
			})
		})
		context("with two describe blocks and two it blocks in each", func() {
			var s2 *Suite
			var suites [][]*step
			beforeEach(func() {
				s2 = NewTestSuite(t)
				describe2, it2 := s2.Describe, s2.It

				// unit under test
				describe2("testing describe", func() {
					it2("first it", func() {
					})
					it2("second it", func() {
					})
				})

				// unit under test
				describe2("testing another describe", func() {
					it2("third it", func() {
					})
					it2("forth it", func() {
					})
				})

				suites = s2.buildSuites()
			})
			it("should have 6 steps defined", func() {
				expect(s2.steps).To.Have.LengthOf(6)
			})
			it("should have 4 suites with just 2 steps defined in each", func() {
				expect(suites).To.Have.LengthOf(4)
				expect(suites[0]).To.Have.LengthOf(2)
				expect(buildSuiteTitle(suites[0])).To.Be.EqualTo("testing describe/first it")
				expect(suites[1]).To.Have.LengthOf(2)
				expect(buildSuiteTitle(suites[1])).To.Be.EqualTo("testing describe/second it")
				expect(suites[2]).To.Have.LengthOf(2)
				expect(buildSuiteTitle(suites[2])).To.Be.EqualTo("testing another describe/third it")
				expect(suites[3]).To.Have.LengthOf(2)
				expect(buildSuiteTitle(suites[3])).To.Be.EqualTo("testing another describe/forth it")
			})
		})
		context("with a more complex example", func() {
			var s2 *Suite
			var suites [][]*step
			beforeEach(func() {
				s2 = NewTestSuite(t)
				describe2, it2 := s2.Describe, s2.It

				// unit under test
				describe2("testing describe", func() {
					it2("first it", func() {
					})
					it2("second it", func() {
					})
					describe2("testing nested context", func() {
						it2("third it", func() {
						})
						describe2("testing nested empty context", func() {
						})
					})
				})

				suites = s2.buildSuites()
			})
			it("should have 6 steps defined", func() {
				expect(s2.steps).To.Have.LengthOf(6)
			})
			it("should have 4 suites with just 2 steps defined in each", func() {
				expect(suites).To.Have.LengthOf(4)
				expect(suites[0]).To.Have.LengthOf(2)
				expect(buildSuiteTitle(suites[0])).To.Be.EqualTo("testing describe/first it")
				expect(suites[1]).To.Have.LengthOf(2)
				expect(buildSuiteTitle(suites[1])).To.Be.EqualTo("testing describe/second it")
				expect(suites[2]).To.Have.LengthOf(3)
				expect(buildSuiteTitle(suites[2])).To.Be.EqualTo("testing describe/testing nested context/third it")
				expect(suites[3]).To.Have.LengthOf(3)
				expect(buildSuiteTitle(suites[3])).To.Be.EqualTo("testing describe/testing nested context/testing nested empty context")
			})
		})
		context("with an even more complex example", func() {
			var s2 *Suite
			var suites [][]*step
			beforeEach(func() {
				s2 = NewTestSuite(t)
				describe2, context2, beforeEach2, it2 :=
					s2.Describe, s2.Describe, s2.BeforeEach, s2.It

				// unit under test
				describe2("Checkout", func() {
					context2("when shopping cart has 1 item", func() {
						var cart []string
						var appliedCoupons []string

						beforeEach2(func() {
							cart = []string{"Gopher toy"}
						})
						it2("should have 1 item in the cart", func() {
							expect(cart).To.Have.LengthOf(1)
						})
						it2("should have no coupon applied by default", func() {
							expect(appliedCoupons).To.Have.LengthOf(0)
						})
						context2("when we add one more item to the cart", func() {
							beforeEach2(func() {
								cart = append(cart, "Crab toy")
							})
							it2("should have 2 items in the cart", func() {
								expect(cart).To.Have.LengthOf(2)
							})
							context2("when the coupon is eligible for this purchase", func() {
								context2("and the coupon gets applied", func() {
									context2("but the coupon value is higher than the purcahse value", func() {
										// ..
									})
									context2("and the coupon value is less than the purchase value", func() {
										beforeEach2(func() {
											// ..
										})
										describe2("when completing the purchase", func() {
											// ..
										})
									})
								})
							})
						})
					})
					context2("when shopping cart is empty", func() {
						var cart []string
						it2("should have 0 items", func() {
							expect(cart).To.Have.LengthOf(0)
						})
					})
				})

				describe2("Sign Up", func() {
					var signedUp bool
					context2("when the user signs up", func() {
						beforeEach2(func() {
							signedUp = true
						})
						it2("should be signed up", func() {
							expect(signedUp).To.Be.True()
						})
					})
				})

				suites = s2.buildSuites()
			})
			it("should have 5 steps defined", func() {
				expect(s2.steps).To.Have.LengthOf(20)
			})
			it("should have 4 suites with just 2 steps defined in each", func() {
				expect(suites).To.Have.LengthOf(7)
				expect(suites[0]).To.Have.LengthOf(4)
				expect(buildSuiteTitle(suites[0])).To.Be.EqualTo("Checkout/when shopping cart has 1 item/should have 1 item in the cart")
				expect(suites[1]).To.Have.LengthOf(4)
				expect(buildSuiteTitle(suites[1])).To.Be.EqualTo("Checkout/when shopping cart has 1 item/should have no coupon applied by default")
				expect(suites[2]).To.Have.LengthOf(6)
				expect(buildSuiteTitle(suites[2])).To.Be.EqualTo("Checkout/when shopping cart has 1 item/when we add one more item to the cart/should have 2 items in the cart")
				expect(suites[3]).To.Have.LengthOf(8)
				expect(buildSuiteTitle(suites[3])).To.Be.EqualTo("Checkout/when shopping cart has 1 item/when we add one more item to the cart/when the coupon is eligible for this purchase/and the coupon gets applied/but the coupon value is higher than the purcahse value")
				expect(suites[4]).To.Have.LengthOf(10)
				expect(buildSuiteTitle(suites[4])).To.Be.EqualTo("Checkout/when shopping cart has 1 item/when we add one more item to the cart/when the coupon is eligible for this purchase/and the coupon gets applied/and the coupon value is less than the purchase value/when completing the purchase")
				expect(suites[5]).To.Have.LengthOf(3)
				expect(buildSuiteTitle(suites[5])).To.Be.EqualTo("Checkout/when shopping cart is empty/should have 0 items")
				expect(suites[6]).To.Have.LengthOf(4)
				expect(buildSuiteTitle(suites[6])).To.Be.EqualTo("Sign Up/when the user signs up/should be signed up")
			})
		})

		context("with another even more complex example", func() {
			var s2 *Suite
			var suites [][]*step
			beforeEach(func() {
				s2 = NewTestSuite(t)
				describe2, context2, beforeEach2, it2 :=
					s2.Describe, s2.Describe, s2.BeforeEach, s2.It

				// unit under test
				describe2("1", func() {
					context2("1.1", func() {
						beforeEach2(func() {
						})
						it2("1.1.1", func() {
						})
					})
					context2("1.2", func() {
						beforeEach2(func() {
						})
						it2("1.2.1", func() {
						})
					})
				})

				suites = s2.buildSuites()
			})
			it("should have 5 steps defined", func() {
				expect(s2.steps).To.Have.LengthOf(7)
			})
			it("should have 4 suites with just 2 steps defined in each", func() {
				expect(suites).To.Have.LengthOf(2)
				expect(suites[0]).To.Have.LengthOf(4)
				expect(buildSuiteTitle(suites[0])).To.Be.EqualTo("1/1.1/1.1.1")
				expect(suites[1]).To.Have.LengthOf(4)
				expect(buildSuiteTitle(suites[1])).To.Be.EqualTo("1/1.2/1.2.1")
			})
		})
	})
}

func TestUsingTestTable(t *testing.T) {
	s := NewTestSuite(t)
	describe, it, expect, start := s.Describe, s.It, s.Expect, s.Start
	defer start()

	describe("using a test table", func() {
		testCases := []struct {
			left     bool
			right    bool
			expected bool
		}{
			{true, true, true},
			{false, false, true},
			{true, false, false},
		}

		for _, tc := range testCases {
			tc := tc
			it(fmt.Sprintf("should be %v when %v == %v", tc.expected, tc.left, tc.right), func() {
				expect(tc.left == tc.right).To.Be.EqualTo(tc.expected)
			})
		}
	})
}
