package examples_test

import (
	"testing"

	"github.com/slavsan/gospec"
	"github.com/slavsan/gospec/internal/testing/helpers/assert"
)

func TestSimpleSpec(t *testing.T) {
	gospec.TestSuite(t, func(s *gospec.Suite) {
		describe, beforeEach, it := s.API()

		describe("Cart", func() {
			var cart []string

			beforeEach(func() {
				cart = []string{
					"Gopher Toy",
					"Crab Toy",
				}
			})

			describe("cart updates", func() {
				describe("given a new item has already been added", func() {
					beforeEach(func() {
						cart = append(cart, "Lizard toy")
					})

					describe("when we remove the second item", func() {
						beforeEach(func() {
							cart = []string{cart[0], cart[2]}
						})

						it("then the cart should contain the correct two items", func() {
							assert.Equal(t, []string{"Gopher Toy", "Lizard toy"}, cart)
						})
					})
				})
			})

			describe("removing items from the cart", func() {
				describe("given the second item has already been removed", func() {
					beforeEach(func() {
						cart = cart[:1]
					})

					describe("when we remove the first item", func() {
						beforeEach(func() {
							cart = cart[:0]
						})

						it("then the cart should contain 0 items", func() {
							assert.Equal(t, []string{}, cart)
						})
					})
				})
			})
		})
	})
}
