package examples_test

import (
	"testing"

	"github.com/slavsan/gospec"
	"github.com/slavsan/gospec/internal/testing/helpers/assert"
)

func TestFeatureInParallelExample(t *testing.T) {
	t.Parallel()

	gospec.WithFeatureSuite(t, func(s *gospec.FeatureSuite) {
		feature, background, scenario, given, when, then, _ :=
			s.With(gospec.Parallel(func() {})).API()

		feature("Cart", func() {
			background(func() {
				given("there is a cart with three items", func(t *T, w *gospec.World) {
					w.Set("cart", []string{
						"Gopher Toy",
						"Crab Toy",
					})
				})
			})

			scenario("cart updates", func() {
				given("a new item has already been added", func(t *testing.T, w *gospec.World) {
					w.Swap("cart", func(cart any) any { return append(cart.([]string), "Lizard toy") })
				})
				when("we remove the second item", func(t *testing.T, w *gospec.World) {
					w.Swap("cart", func(cart any) any { c := cart.([]string); return []string{c[0], c[2]} })
				})
				then("the cart should contain the correct two items", func(t *testing.T, w *gospec.World) {
					assert.Equal(t, []string{"Gopher Toy", "Lizard toy"}, w.Get("cart"))
				})
			})

			scenario("removing items from the cart", func() {
				given("the second item has already been removed", func(t *T, w *gospec.World) {
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
