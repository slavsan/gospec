package examples_test

import (
	"testing"

	"github.com/slavsan/gospec"
	"github.com/slavsan/gospec/internal/testing/helpers/assert"
)

func TestSimpleFeatureExample(t *testing.T) {
	feature, background, scenario, given, when, then, _ :=
		gospec.NewFeatureSuite(t, gospec.WithPrintedFilenames()).API()

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
}
