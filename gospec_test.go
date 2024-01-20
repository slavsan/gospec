package gospec

import (
	"fmt"
	"testing"

	"github.com/slavsan/gospec/internal/testing/helpers/assert"
)

func TestExampleSuite(t *testing.T) {
	s := NewTestSuite(t)
	describe, context, it, beforeEach := s.Describe, s.Describe, s.It, s.BeforeEach

	describe("Checkout", func() {
		context("when shopping cart has 1 item", func() {
			var cart []string
			var appliedCoupons []string

			beforeEach(func() {
				cart = []string{"Gopher toy"}
			})
			it("should have 1 item in the cart", func() {
				assert.Equal(t, 1, len(cart))
			})
			it("should have no coupon applied by default", func() {
				assert.Equal(t, 0, len(appliedCoupons))
			})
			context("when we add one more item to the cart", func() {
				beforeEach(func() {
					cart = append(cart, "Crab toy")
				})
				it("should have 2 items in the cart", func() {
					assert.Equal(t, 2, len(cart))
				})
				context("when the coupon is eligible for this purchase", func() {
					context("and the coupon gets applied", func() {
						context("but the coupon value is higher than the purchase value", func() {
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
				assert.Equal(t, 0, len(cart))
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
				assert.Equal(t, true, signedUp)
			})
		})
	})
}

func TestDescribe(t *testing.T) {
	s := NewTestSuite(t)
	describe, context, it, beforeEach := s.Describe, s.Describe, s.It, s.BeforeEach

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
			})
			it("should have empty stack", func() {
				assert.Equal(t, 0, len(s2.stack))
			})
			it("should have just no suites defined", func() {
				assert.Equal(t, 0, len(suites))
			})
		})
		context("with two sibling describes", func() {
			var s2 *Suite
			beforeEach(func() {
				s2 = NewTestSuite(t)
				describe2 := s2.Describe

				// unit under test
				describe2("describe", func() {
				})
				describe2("sibling describe", func() {
				})
			})
			it("should have empty stack", func() {
				assert.Equal(t, 0, len(s2.stack))
			})
			it("should have 0 suites", func() {
				assert.Equal(t, 0, len(s2.suites))
			})
		})
		context("with two describes, one parent and one child", func() {
			var s2 *Suite
			beforeEach(func() {
				s2 = NewTestSuite(t)
				describe2 := s2.Describe

				// unit under test
				describe2("parent describe", func() {
					describe2("child describe", func() {
					})
				})
			})
			it("should have empty stack", func() {
				assert.Equal(t, 0, len(s2.stack))
			})
			it("should have 0 suites", func() {
				assert.Equal(t, 0, len(s2.suites))
			})
		})
		context("with three nested describes", func() {
			var s2 *Suite
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
			})
			it("should have empty stack", func() {
				assert.Equal(t, 0, len(s2.stack))
			})
			it("should have 0 suites", func() {
				assert.Equal(t, 0, len(s2.suites))
			})
		})
	})
}

func TestIt(t *testing.T) { //nolint:maintidx
	s := NewTestSuite(t)
	describe, context, it, beforeEach := s.Describe, s.Describe, s.It, s.BeforeEach

	describe("it block", func() {
		context("with single describe and single it block", func() {
			var s2 *Suite
			beforeEach(func() {
				s2 = NewTestSuite(t)
				describe2, it2 := s2.Describe, s2.It

				describe2("testing describe", func() {
					it2("testing it", func() {
					})
				})
			})
			it("should have empty stack", func() {
				assert.Equal(t, 0, len(s2.stack))
			})
			it("should have just 1 suite with just 2 steps defined", func() {
				assert.Equal(t, 1, len(s2.suites))
				assert.Equal(t, 2, len(s2.suites[0]))
				assert.Equal(t, "testing describe/testing it", buildSuiteTitle(s2.suites[0]))
			})
		})
		context("with single describe and two it blocks", func() {
			var s2 *Suite
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
			})
			it("should have empty stack", func() {
				assert.Equal(t, 0, len(s2.stack))
			})
			it("should have 2 suites with just 2 steps defined in each", func() {
				assert.Equal(t, 2, len(s2.suites))
				assert.Equal(t, 2, len(s2.suites[0]))
				assert.Equal(t, "testing describe/testing it", buildSuiteTitle(s2.suites[0]))
				assert.Equal(t, 2, len(s2.suites[1]))
				assert.Equal(t, "testing describe/testing another it", buildSuiteTitle(s2.suites[1]))
			})
		})
		context("with two describe blocks and one it block in each", func() {
			var s2 *Suite
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
			})
			it("should have empty stack", func() {
				assert.Equal(t, 0, len(s2.stack))
			})
			it("should have 2 suites with just 2 steps defined in each", func() {
				assert.Equal(t, 2, len(s2.suites))
				assert.Equal(t, 2, len(s2.suites[0]))
				assert.Equal(t, "testing describe/testing it", buildSuiteTitle(s2.suites[0]))
				assert.Equal(t, 2, len(s2.suites[1]))
				assert.Equal(t, "testing another describe/testing another it", buildSuiteTitle(s2.suites[1]))
			})
		})
		context("with two describe blocks and two it blocks in each", func() {
			var s2 *Suite
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
			})
			it("should have empty stack", func() {
				assert.Equal(t, 0, len(s2.stack))
			})
			it("should have 4 suites with just 2 steps defined in each", func() {
				assert.Equal(t, 4, len(s2.suites))
				assert.Equal(t, 2, len(s2.suites[0]))
				assert.Equal(t, "testing describe/first it", buildSuiteTitle(s2.suites[0]))
				assert.Equal(t, 2, len(s2.suites[1]))
				assert.Equal(t, "testing describe/second it", buildSuiteTitle(s2.suites[1]))
				assert.Equal(t, 2, len(s2.suites[2]))
				assert.Equal(t, "testing another describe/third it", buildSuiteTitle(s2.suites[2]))
				assert.Equal(t, 2, len(s2.suites[3]))
				assert.Equal(t, "testing another describe/forth it", buildSuiteTitle(s2.suites[3]))
			})
		})
		context("with a more complex example", func() {
			var s2 *Suite
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
			})
			it("should have empty stack", func() {
				assert.Equal(t, 0, len(s2.stack))
			})
			it("should have 4 suites with just 2 steps defined in each", func() {
				assert.Equal(t, 3, len(s2.suites))
				assert.Equal(t, 2, len(s2.suites[0]))
				assert.Equal(t, "testing describe/first it", buildSuiteTitle(s2.suites[0]))
				assert.Equal(t, 2, len(s2.suites[1]))
				assert.Equal(t, "testing describe/second it", buildSuiteTitle(s2.suites[1]))
				assert.Equal(t, 3, len(s2.suites[2]))
				assert.Equal(t, "testing describe/testing nested context/third it", buildSuiteTitle(s2.suites[2]))
			})
		})
		context("with an even more complex example", func() {
			var s2 *Suite
			beforeEach(func() {
				s2 = NewTestSuite(t)
				describe2, context2, beforeEach2, it2 := s2.Describe, s2.Describe, s2.BeforeEach, s2.It

				// unit under test
				describe2("Checkout", func() {
					context2("when shopping cart has 1 item", func() {
						var cart []string
						var appliedCoupons []string

						beforeEach2(func() {
							cart = []string{"Gopher toy"}
						})
						it2("should have 1 item in the cart", func() {
							assert.Equal(t, 1, len(cart))
						})
						it2("should have no coupon applied by default", func() {
							assert.Equal(t, 0, len(appliedCoupons))
						})
						context2("when we add one more item to the cart", func() {
							beforeEach2(func() {
								cart = append(cart, "Crab toy")
							})
							it2("should have 2 items in the cart", func() {
								assert.Equal(t, 2, len(cart))
							})
							context2("when the coupon is eligible for this purchase", func() {
								context2("and the coupon gets applied", func() {
									context2("but the coupon value is higher than the purchase value", func() {
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
							assert.Equal(t, 0, len(cart))
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
							assert.Equal(t, true, signedUp)
						})
					})
				})
			})
			it("should have an empty stack", func() {
				assert.Equal(t, 0, len(s2.stack))
			})
			it("should have 4 suites with just 2 steps defined in each", func() {
				assert.Equal(t, 5, len(s2.suites))
				assert.Equal(t, 4, len(s2.suites[0]))
				assert.Equal(t, "Checkout/when shopping cart has 1 item/should have 1 item in the cart", buildSuiteTitle(s2.suites[0]))
				assert.Equal(t, 4, len(s2.suites[1]))
				assert.Equal(t, "Checkout/when shopping cart has 1 item/should have no coupon applied by default", buildSuiteTitle(s2.suites[1]))
				assert.Equal(t, 6, len(s2.suites[2]))
				assert.Equal(t, "Checkout/when shopping cart has 1 item/when we add one more item to the cart/should have 2 items in the cart", buildSuiteTitle(s2.suites[2]))
				assert.Equal(t, 3, len(s2.suites[3]))
				assert.Equal(t, "Checkout/when shopping cart is empty/should have 0 items", buildSuiteTitle(s2.suites[3]))
				assert.Equal(t, 4, len(s2.suites[4]))
				assert.Equal(t, "Sign Up/when the user signs up/should be signed up", buildSuiteTitle(s2.suites[4]))
			})
		})

		context("with another even more complex example", func() {
			var s2 *Suite
			beforeEach(func() {
				s2 = NewTestSuite(t)
				describe2, context2, beforeEach2, it2 := s2.Describe, s2.Describe, s2.BeforeEach, s2.It

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
			})
			it("should have empty stack", func() {
				assert.Equal(t, 0, len(s2.stack))
			})
			it("should have 4 suites with just 2 steps defined in each", func() {
				assert.Equal(t, 2, len(s2.suites))
				assert.Equal(t, 4, len(s2.suites[0]))
				assert.Equal(t, "1/1.1/1.1.1", buildSuiteTitle(s2.suites[0]))
				assert.Equal(t, 4, len(s2.suites[1]))
				assert.Equal(t, "1/1.2/1.2.1", buildSuiteTitle(s2.suites[1]))
			})
		})
	})
}

func TestUsingTestTable(t *testing.T) {
	s := NewTestSuite(t)
	describe, it := s.Describe, s.It

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
				assert.Equal(t, tc.expected, tc.left == tc.right)
			})
		}
	})
}
