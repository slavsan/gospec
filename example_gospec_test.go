package gospec_test

import (
	"testing"

	"github.com/slavsan/gospec"
	"github.com/slavsan/gospec/internal/testing/helpers/assert"
)

// Do not execute this example because it will fail. It's not a real
// test and its purpose is purely presentational.

// ExampleTestSuite shows an example gospec.SpecSuite.
func Example_suite() {
	gospec.WithSpecSuite(t, func(s *gospec.SpecSuite) {
		describe, beforeEach, it := s.API()

		describe("Cart", func() {
			var cart []string

			beforeEach(func(t *testing.T) {
				cart = []string{
					"Gopher Toy",
					"Crab Toy",
				}
			})

			describe("cart updates", func() {
				describe("given a new item has already been added", func() {
					beforeEach(func(t *testing.T) {
						cart = append(cart, "Lizard toy")
					})

					describe("when we remove the second item", func() {
						beforeEach(func(t *testing.T) {
							cart = []string{cart[0], cart[2]}
						})

						it("then the cart should contain the correct two items", func(t *testing.T) {
							assert.Equal(t, []string{"Gopher Toy", "Lizard toy"}, cart)
						})
					})
				})
			})

			describe("removing items from the cart", func() {
				describe("given the second item has already been removed", func() {
					beforeEach(func(t *testing.T) {
						cart = cart[:1]
					})

					describe("when we remove the first item", func() {
						beforeEach(func(t *testing.T) {
							cart = cart[:0]
						})

						it("then the cart should contain 0 items", func(t *testing.T) {
							assert.Equal(t, []string{}, cart)
						})
					})
				})
			})
		})
	})
}
