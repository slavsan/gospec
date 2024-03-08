package examples_test

import (
	"os"
	"testing"

	"github.com/slavsan/gospec"
	"github.com/slavsan/gospec/internal/testing/helpers/assert"
)

type T = testing.T

func TestSimpleSpec(t *testing.T) {
	t.Parallel()

	gospec.WithSpecSuite(t, func(s *gospec.SpecSuite) {
		output := gospec.Output(os.Stdout, gospec.Colorful, gospec.Durations, gospec.PrintFilenames)
		describe, beforeEach, it := s.With(output).API()

		describe("Cart", func() {
			var cart []string

			beforeEach(func(t *T) {
				cart = []string{
					"Gopher Toy",
					"Crab Toy",
				}
			})

			describe("cart updates", func() {
				describe("given a new item has already been added", func() {
					beforeEach(func(t *T) {
						cart = append(cart, "Lizard toy")
					})

					describe("when we remove the second item", func() {
						beforeEach(func(t *T) {
							cart = []string{cart[0], cart[2]}
						})

						it("then the cart should contain the correct two items", func(t *T) {
							assert.Equal(t, []string{"Gopher Toy", "Lizard toy"}, cart)
						})
					})
				})
			})

			describe("removing items from the cart", func() {
				describe("given the second item has already been removed", func() {
					beforeEach(func(t *T) {
						cart = cart[:1]
					})

					describe("when we remove the first item", func() {
						beforeEach(func(t *T) {
							cart = cart[:0]
						})

						it("then the cart should contain 0 items", func(t *T) {
							assert.Equal(t, []string{}, cart)
						})
					})
				})
			})
		})
	})
}
