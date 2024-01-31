package gospec_test

import (
	"testing"

	"github.com/slavsan/gospec"
	"github.com/slavsan/gospec/internal/testing/helpers/assert"
)

// Do not execute this example because it will fail. It's not a real
// test and its purpose is purely presentational.

// ExampleParallelSpecSuite shows an example gospec.SpecSuite which runs tests in parallel.
func Example_parallelSpecSuite() {
	gospec.WithSpecSuite(t, func(s *gospec.SpecSuite) {
		describe, beforeEach, it := s.ParallelAPI(func() { /* optionally, execute this once all parallel tests have finished */ })

		describe("Cart", func() {
			beforeEach(func(t *testing.T, w *gospec.World) {
				w.Set("cart", []string{
					"Gopher Toy",
					"Crab Toy",
				})
			})

			describe("cart updates", func() {
				describe("given a new item has already been added", func() {
					beforeEach(func(t *testing.T, w *gospec.World) {
						w.Swap("cart", func(cart any) any { return append(cart.([]string), "Lizard toy") })
					})

					describe("when we remove the second item", func() {
						beforeEach(func(t *testing.T, w *gospec.World) {
							w.Swap("cart", func(cart any) any { c := cart.([]string); return []string{c[0], c[2]} })
						})

						it("then the cart should contain the correct two items", func(t *testing.T, w *gospec.World) {
							assert.Equal(t, []string{"Gopher Toy", "Lizard toy"}, w.Get("cart"))
						})
					})
				})
			})

			describe("removing items from the cart", func() {
				describe("given the second item has already been removed", func() {
					beforeEach(func(t *testing.T, w *gospec.World) {
						w.Swap("cart", func(cart any) any { return cart.([]string)[:1] })
					})

					describe("when we remove the first item", func() {
						beforeEach(func(t *testing.T, w *gospec.World) {
							w.Swap("cart", func(cart any) any { return cart.([]string)[:0] })
						})

						it("then the cart should contain 0 items", func(t *testing.T, w *gospec.World) {
							assert.Equal(t, []string{}, w.Get("cart"))
						})
					})
				})
			})
		})
	})
}
