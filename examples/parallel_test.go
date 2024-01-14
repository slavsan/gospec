package examples_test

import (
	"testing"

	"github.com/slavsan/gospec"
	"github.com/slavsan/gospec/internal/testing/helpers/assert"
)

func TestParallel(t *testing.T) {
	t.Parallel()

	feature, background, scenario, given, when, then, _ :=
		gospec.NewFeatureSuite(t, gospec.WithParallel()).API()

	feature("Cart", func() {
		background("", func() {
			given("there is a cart with three items", func(w *gospec.World) {
				w.Set("cart", []string{
					"Gopher Toy",
					"Crab Toy",
				})
			})
		})

		scenario("cart updates", func() {
			given("a new item has already been added", func(w *gospec.World) {
				w.Swap("cart", func(cart any) any { return append(cart.([]string), "Lizard toy") })
			})
			when("we remove the second item", func(w *gospec.World) {
				w.Swap("cart", func(cart any) any { c := cart.([]string); return []string{c[0], c[2]} })
			})
			then("the cart should contain the correct two items", func(w *gospec.World) {
				assert.Equal(w.T, []string{"Gopher Toy", "Lizard toy"}, w.Get("cart"))
			})
		})

		scenario("removing items from the cart", func() {
			given("the second item has already been removed", func(w *gospec.World) {
				w.Swap("cart", func(cart any) any { return cart.([]string)[:1] })
			})
			when("we remove the first item", func(w *gospec.World) {
				w.Swap("cart", func(cart any) any { return cart.([]string)[:0] })
			})
			then("the cart should contain 0 items", func(w *gospec.World) {
				assert.Equal(w.T, []string{}, w.Get("cart"))
			})
		})
	})
}
