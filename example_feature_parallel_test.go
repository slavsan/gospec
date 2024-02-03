package gospec_test

import (
	"os"
	"testing"

	"github.com/slavsan/gospec"
	"github.com/slavsan/gospec/internal/testing/helpers/assert"
)

// Do not execute this example because it will fail. It's not a real
// test and its purpose is purely presentational.

// ExampleParallelFeatureSuite shows an example gospec.FeatureSuite which runs tests in parallel.
func Example_parallelFeatureSuite() {
	gospec.WithFeatureSuite(t, func(s *gospec.FeatureSuite) {
		feature, background, scenario, given, when, then := s.With(gospec.Output(os.Stdout, gospec.PrintFilenames)).ParallelAPI(func() {
			/* optionally, execute this once all parallel tests have finished */
		})

		feature("Cart", func() {
			background(func() {
				given("there is a cart with three items", func(t *testing.T, w *gospec.World) {
					// use the `World.Set` method to set the initial state of some variable (in this case named `cart`)
					w.Set("cart", []string{
						"Gopher Toy",
						"Crab Toy",
					})
				})
			})

			scenario("cart updates", func() {
				given("a new item has already been added", func(t *testing.T, w *gospec.World) {
					// Update the state of the cart variable
					w.Swap("cart", func(cart any) any { return append(cart.([]string), "Lizard toy") })
				})
				when("we remove the second item", func(t *testing.T, w *gospec.World) {
					// update the state of the cart variable again
					w.Swap("cart", func(cart any) any { c := cart.([]string); return []string{c[0], c[2]} })
				})
				then("the cart should contain the correct two items", func(t *testing.T, w *gospec.World) {
					// assert the state is as expected
					assert.Equal(t, []string{"Gopher Toy", "Lizard toy"}, w.Get("cart"))
				})
			})

			scenario("removing items from the cart", func() {
				given("the second item has already been removed", func(t *testing.T, w *gospec.World) {
					w.Swap("cart", func(cart any) any { return cart.([]string)[:1] })
				})
				when("we remove the first item", func(t *testing.T, w *gospec.World) {
					w.Swap("cart", func(cart any) any { return cart.([]string)[:0] })
				})
				then("the cart should contain 0 items", func(t *testing.T, w *gospec.World) {
					assert.Equal(t, []string{}, w.Get("cart"))
				})
			})
		})
	})
}
