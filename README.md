# gospec

Gospec is a BDD testing library for the Go programming language.

It's influenced by rspec/mocha and other similar libraries, and also the Cucumber testing library. Gospec provides a similar and familiar developer experience.

The language is a thin layer on top of the standard Go tests.

Gospec has zero dependencies and the intention is for this to not change.

## Usage

The intended usage is to have the standard Go tests, defined like so `func TestMyFeature(t *testing.T) {`, then use the gospec library to initialize a `Suite` or `FeatureSuite` depending on which API is preferred.

### Test suites

```go
import (
	"github.com/slavsan/gospec"
	...
)

func TestCartSpec() {
	describe, beforeEach, it := gospec.NewTestSuite(t).API()

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
}
```

### Feature suites

```go
import (
	"github.com/slavsan/gospec"
)

func TestCartFeature() {
	feature, background, scenario, given, when, then, _ :=
		gospec.NewFeatureSuite(t, gospec.WithPrintedFilenames()).API()

	feature("Cart", func() {
		var cart []string

		background(func() {
			given("there is a cart with three items", func() {
				cart = []string{
					"Gopher Toy",
					"Crab Toy",
				}
			})
		})

		scenario("cart updates", func() {
			given("a new item has already been added", func() {
				cart = append(cart, "Lizard toy")
			})
			when("we remove the second item", func() {
				cart = []string{cart[0], cart[2]}
			})
			then("the cart should contain the correct two items", func() {
				assert.Equal(t, []string{"Gopher Toy", "Lizard toy"}, cart)
			})
		})

		scenario("removing items from the cart", func() {
			given("the second item has already been removed", func() {
				cart = cart[:1]
			})
			when("we remove the first item", func() {
				cart = cart[:0]
			})
			then("the cart should contain 0 items", func() {
				assert.Equal(t, []string{}, cart)
			})
		})
	})
}
```

## Options

When initializing `Suite` or `FeatureSuite` instances, there are several options to enhance the developer experience.

- `WithParallel` for parallel execution of tests
- `WithOutput` for specifying a custom output writer
- `WithPrintedFilenames` for printing the `filename:line` on each line where there is a `describe`, `beforeEach`, `it`, `feature`, `scenario` (or any other API method) defined. An editor which supports recognition of paths and allows for jumping to source, could improve the developer experience when dealing with maintenance of existing tests/specs.

### Parallel execution

Gospec supports parallel execution of tests.

One important caveat is that in order to do that, your tests would need to utilize the `World` struct which gets passed to all of your test functions, e.g. `beforeEach`, `it`, etc., or `given`, `when`, `then` in the case of `FeatureSpec` suites.

This is necessary since the way `gospec` works, it needs to store any contextual data used in one suite separate from the scope of other tests.
