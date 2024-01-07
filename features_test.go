package gospec

import "testing"

func TestFeature(t *testing.T) {
	// ..

	f := NewFeatureSuite(t)
	feature, background, scenario, given, when, then, world, table, expect, start :=
		f.Feature, f.Background, f.Scenario, f.Given, f.When, f.Then, f.World, f.Table, f.Expect, f.Start
	defer start()

	_ = world

	feature("Checkout", func() {
		type Product struct {
			Name  string
			Price float64
			Type  int
			// ..
		}

		var (
			customer       string
			cart           []string
			appliedCoupons []string
			items          []Product
		)

		_ = customer
		_ = cart
		_ = appliedCoupons

		background("", func() {
			given("there is a customer account", func() {
				customer = "John Doe"
				// world.Define(func() {
				// 	world.Var("foo", "bar")
				// 	world.Var("spam", "eggs")
				// })
			})
			given("the cart has 2 items", func() {
				items = []Product{
					{Name: "Gopher toy", Price: 14.99, Type: 2},
					{Name: "Crab toy", Price: 17.49, Type: 8},
				}
				table([]string{"Name", "Price"}, items)
				// ..
			})
			given("no coupons are applied", func() {
				appliedCoupons = []string{}
			})
			// ..
		})

		scenario("with no coupons", func() {
			given("has the correct total price of both items", func() {
				cart = []string{}
				// ..
			})
			when("the customer checks out", func() {
				cart = []string{}
				// ..
			})
			then("the cart should be 0 again", func() {
				// ..
				expect(cart).To.Be.EqualTo(0)
			})
			// ...
		})

		scenario("with coupon applied", func() {
			given("has the correct total price of both items", func() {
				cart = []string{}
				// ..
			})
			when("the customer checks out", func() {
				cart = []string{}
				// ..
			})
			then("the cart should be 0 again", func() {
				// ..
				expect(cart).To.Be.EqualTo(0)
			})
			// ...
		})
	})

	// debugFeature(f)
}
