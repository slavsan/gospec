package examples_test

import (
	"testing"

	"github.com/slavsan/gospec"
	"github.com/slavsan/gospec/internal/testing/helpers/assert"
)

type W = gospec.W
type T = testing.T

func TestSimpleSpec(t *testing.T) {
	gospec.WithSpecSuite(t, func(s *gospec.SpecSuite) {
		describe, beforeEach, it := s.API()

		describe("Cart", func() {
			var cart []string

			beforeEach(func(t *T, w *W) {
				cart = []string{
					"Gopher Toy",
					"Crab Toy",
				}
			})

			describe("cart updates", func() {
				describe("given a new item has already been added", func() {
					beforeEach(func(t *T, w *W) {
						cart = append(cart, "Lizard toy")
					})

					describe("when we remove the second item", func() {
						beforeEach(func(t *T, w *W) {
							cart = []string{cart[0], cart[2]}
						})

						it("then the cart should contain the correct two items", func(t *T, w *gospec.World) {
							assert.Equal(t, []string{"Gopher Toy", "Lizard toy"}, cart)
						})
					})
				})
			})

			describe("removing items from the cart", func() {
				describe("given the second item has already been removed", func() {
					beforeEach(func(t *T, w *W) {
						cart = cart[:1]
					})

					describe("when we remove the first item", func() {
						beforeEach(func(t *T, w *W) {
							cart = cart[:0]
						})

						it("then the cart should contain 0 items", func(t *T, w *W) {
							assert.Equal(t, []string{}, cart)
						})
					})
				})
			})
		})
	})
}
