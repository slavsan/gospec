package examples_test

import (
	"testing"

	"github.com/slavsan/gospec"
	"github.com/slavsan/gospec/internal/testing/helpers/assert"
)

func TestSimple(t *testing.T) {
	feature, background, scenario, given, when, then, _, _ :=
		gospec.NewFeatureSuite(t).API()

	feature("Cart", func() {
		var cart []string

		background("", func() {
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
	})
}
