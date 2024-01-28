package gospec_test

import (
	"github.com/slavsan/gospec"
	"github.com/slavsan/gospec/internal/testing/helpers/assert"
)

// Do not execute this example because it will fail. It's not a real
// test and its purpose is purely presentational.

// ExampleFeatureSuite shows an example gospec.FeatureSuite.
func Example_featureSuite() {
	gospec.FeatureSuite2(t, func(s *gospec.FeatureSuite) {
		feature, background, scenario, given, when, then, _ :=
			s.With(gospec.WithPrintedFilenames()).API()

		feature("Cart", func() {
			var cart []string

			background(func() {
				given("there is a cart with three items", func() {
					cart = []string{
						"Gopher Toy",
						"Crab Toy",
					}
				})
			})

			scenario("cart updates", func() {
				given("a new item has already been added", func() {
					cart = append(cart, "Lizard toy")
				})
				when("we remove the second item", func() {
					cart = []string{cart[0], cart[2]}
				})
				then("the cart should contain the correct two items", func() {
					assert.Equal(t, []string{"Gopher Toy", "Lizard toy"}, cart)
				})
			})

			scenario("removing items from the cart", func() {
				given("the second item has already been removed", func() {
					cart = cart[:1]
				})
				when("we remove the first item", func() {
					cart = cart[:0]
				})
				then("the cart should contain 0 items", func() {
					assert.Equal(t, []string{}, cart)
				})
			})
		})
	})
}
