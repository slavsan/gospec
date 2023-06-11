package gospec

import (
	"testing"
)

func TestExampleSuite(t *testing.T) {
	s := NewTestSuite(t)
	describe, it, beforeEach, expect, start :=
		s.Describe, s.It, s.BeforeEach, s.Expect, s.Start
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
		describe("when shopping cart has 1 item", func() {
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
			describe("when we add one more item to the cart", func() {
				beforeEach(func() {
					cart = append(cart, "Crab toy")
				})
				it("should have 2 items in the cart", func() {
					expect(cart).To.Have.LengthOf(2)
				})
				describe("when the coupon is eligible for this purchase", func() {
					describe("and the coupon gets applied", func() {
						describe("but the coupon value is higher than the purcahse value", func() {
							// ..
						})
						describe("and the coupon value is less than the purchase value", func() {
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
		describe("when shopping cart is empty", func() {
			var cart []string
			it("should have 0 items", func() {
				expect(cart).To.Have.LengthOf(0)
			})
		})
	})

	describe("Sign Up", func() {
		var signedUp bool
		describe("when the user signs up", func() {
			beforeEach(func() {
				signedUp = true
			})
			it("should be signed up", func() {
				expect(signedUp).To.Be.True()
			})
		})
	})
}
