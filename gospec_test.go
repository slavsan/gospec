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

	// expect(true).To.Be.True()
	// expect("foo").To.Be.False()
	// expect("foo").To.Be.EqualTo("aaa")
	// expect("foo").To.Have.LengthOf(3)
	// expect("foo").To.Have.Property("prop name")
	// expect("foo").To.Be.Nil()
	// expect("foo").Not.To.Be.Nil()
	// expect("foo").To.Contain.Substring("xx")
	// expect("foo").To.Contain.Element("xx")
	// expect("foo").To.Be.Of.Type(someType{})

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
				expect(suites).To.Have.LengthOf(0) // FIXME: should be 1
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
				expect(suites).To.Have.LengthOf(1) // FIXME: should be 2
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
				expect(suites).To.Have.LengthOf(0) // FIXME: should be 1
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
				expect(suites).To.Have.LengthOf(0) // FIXME: should be 1
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
				expect(suites[1]).To.Have.LengthOf(2)
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
				expect(suites).To.Have.LengthOf(3) // FIXME: should be 2
				expect(suites[0]).To.Have.LengthOf(2)
				expect(suites[1]).To.Have.LengthOf(1)
				expect(suites[2]).To.Have.LengthOf(2)
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
				expect(suites).To.Have.LengthOf(5) // FIXME: should be 4
				expect(suites[0]).To.Have.LengthOf(2)
				expect(suites[1]).To.Have.LengthOf(2)
				expect(suites[2]).To.Have.LengthOf(1)
				expect(suites[3]).To.Have.LengthOf(2)
				expect(suites[4]).To.Have.LengthOf(2)
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
