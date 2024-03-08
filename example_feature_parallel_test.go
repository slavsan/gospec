package gospec_test

import (
	"testing"

	"github.com/slavsan/gospec"
	"github.com/slavsan/gospec/internal/testing/helpers/assert"
)

// Do not execute this example because it will fail. It's not a real
// test and its purpose is purely presentational.

// ExampleParallelFeatureSuite shows an example gospec.FeatureSuite which runs tests in parallel.
func Example_parallelFeatureSuite() {
	parallelTestsWg.Add(1)
	gospec.WithFeatureSuite(t, func(s *gospec.FeatureSuite) {
		feature, background, scenario, given, when, then := s.ParallelAPI(func() { parallelTestsWg.Done() })

		// As opposed to just declaring variables and using them, as we'd do
		// in a non-parallel/sequential Spec test, in parallel tests we need
		// to store our state in a struct which we can be passed between the
		// separate steps (functions exposed by the ParallelAPI).
		type world struct {
			cart []string
		}

		// This is a simple helper function which type asserts the world struct
		// which you define for your test. This is to avoid unnecessary type
		// assertions in the test steps.
		h := func(w *gospec.World) *world {
			return w.Get("world").(*world)
		}

		feature("Cart", func() {
			background(func() {
				given("there is a cart with three items", func(t *testing.T, w *gospec.World) {
					// Here we're setting the `world` struct by
					// storing it in the gospec's World, using the
					// "world" key, since that is used in the `h`
					// helper function.
					w.Set("world", func() interface{} {
						return &world{
							cart: []string{
								"Gopher Toy",
								"Crab Toy",
							},
						}
					}())
				})
			})

			scenario("cart updates", func() {
				given("a new item has already been added", func(t *testing.T, w *gospec.World) {
					// using `h(w).cart` we get access to the `cart` in a concurrently-safe way
					// because this way all the `cart` instances are bound (scoped) to the current
					// `w` (*gospec.World) instance, i.e. to the current test.
					h(w).cart = append(h(w).cart, "Lizard toy")
				})
				when("we remove the second item", func(t *testing.T, w *gospec.World) {
					c := h(w).cart
					h(w).cart = []string{c[0], c[2]}
				})
				then("the cart should contain the correct two items", func(t *testing.T, w *gospec.World) {
					assert.Equal(t, []string{"Gopher Toy", "Lizard toy"}, h(w).cart)
				})
			})

			scenario("removing items from the cart", func() {
				given("the second item has already been removed", func(t *testing.T, w *gospec.World) {
					h(w).cart = h(w).cart[:1]
				})
				when("we remove the first item", func(t *testing.T, w *gospec.World) {
					h(w).cart = h(w).cart[:0]
				})
				then("the cart should contain 0 items", func(t *testing.T, w *gospec.World) {
					assert.Equal(t, []string{}, h(w).cart)
				})
			})
		})
	})
}
