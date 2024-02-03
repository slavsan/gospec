package examples_test

import (
	"os"
	"testing"

	"github.com/slavsan/gospec"
	"github.com/slavsan/gospec/internal/testing/helpers/assert"
)

func TestSimpleFeatureExample(t *testing.T) {
	gospec.WithFeatureSuite(t, func(s *gospec.FeatureSuite) {
		feature, background, scenario, given, when, then, _ :=
			s.With(gospec.Output(os.Stdout, gospec.PrintFilenames)).API()

		feature("Cart", func() {
			var cart []string

			background(func() {
				given("there is a cart with three items", func(t *T) {
					cart = []string{
						"Gopher Toy",
						"Crab Toy",
					}
				})
			})

			scenario("cart updates", func() {
				given("a new item has already been added", func(t *T) {
					cart = append(cart, "Lizard toy")
				})
				when("we remove the second item", func(t *T) {
					cart = []string{cart[0], cart[2]}
				})
				then("the cart should contain the correct two items", func(t *T) {
					assert.Equal(t, []string{"Gopher Toy", "Lizard toy"}, cart)
				})
			})

			scenario("removing items from the cart", func() {
				given("the second item has already been removed", func(t *T) {
					cart = cart[:1]
				})
				when("we remove the first item", func(t *T) {
					cart = cart[:0]
				})
				then("the cart should contain 0 items", func(t *T) {
					assert.Equal(t, []string{}, cart)
				})
			})
		})
	})
}
