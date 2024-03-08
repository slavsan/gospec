package examples_test

import (
	"testing"

	"github.com/slavsan/gospec"
	"github.com/slavsan/gospec/internal/testing/helpers/assert"
)

func TestFeatureInParallelExample(t *testing.T) {
	t.Parallel()

	parallelTestsWg.Add(1)
	gospec.WithFeatureSuite(t, func(s *gospec.FeatureSuite) {
		feature, background, scenario, given, when, then := s.ParallelAPI(func() { parallelTestsWg.Done() })

		type world struct {
			cart []string
		}

		h := func(w *gospec.World) *world {
			return w.Get("world").(*world)
		}

		feature("Cart", func() {
			background(func() {
				given("there is a cart with three items", func(t *T, w *gospec.World) {
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
				given("the second item has already been removed", func(t *T, w *gospec.World) {
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
