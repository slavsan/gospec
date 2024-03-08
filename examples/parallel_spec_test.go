package examples_test

import (
	"testing"

	"github.com/slavsan/gospec"
	"github.com/slavsan/gospec/internal/testing/helpers/assert"
)

func TestSpecSuiteInParallelExample(t *testing.T) {
	parallelTestsWg.Add(1)
	gospec.WithSpecSuite(t, func(s *gospec.SpecSuite) {
		describe, beforeEach, it := s.ParallelAPI(func() { parallelTestsWg.Done() })

		type world struct {
			cart []string
		}

		h := func(w *gospec.World) *world {
			return w.Get("world").(*world)
		}

		describe("Cart", func() {
			beforeEach(func(t *testing.T, w *gospec.World) {
				w.Set("world", func() interface{} {
					return &world{
						cart: []string{
							"Gopher Toy",
							"Crab Toy",
						},
					}
				}())
			})

			describe("cart updates", func() {
				describe("given a new item has already been added", func() {
					beforeEach(func(t *testing.T, w *gospec.World) {
						h(w).cart = append(h(w).cart, "Lizard toy")
					})

					describe("when we remove the second item", func() {
						beforeEach(func(t *testing.T, w *gospec.World) {
							c := h(w).cart
							h(w).cart = []string{c[0], c[2]}
						})

						it("then the cart should contain the correct two items", func(t *testing.T, w *gospec.World) {
							assert.Equal(t, []string{"Gopher Toy", "Lizard toy"}, h(w).cart)
						})
					})
				})
			})

			describe("removing items from the cart", func() {
				describe("given the second item has already been removed", func() {
					beforeEach(func(t *testing.T, w *gospec.World) {
						h(w).cart = h(w).cart[:1]
					})

					describe("when we remove the first item", func() {
						beforeEach(func(t *testing.T, w *gospec.World) {
							h(w).cart = h(w).cart[:0]
						})

						it("then the cart should contain 0 items", func(t *testing.T, w *gospec.World) {
							assert.Equal(t, []string{}, h(w).cart)
						})
					})
				})
			})
		})
	})
}
