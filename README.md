# gospec

Gospec is a BDD testing library for the Go programming language.

It's influenced by rspec/mocha and other similar libraries, and also the Cucumber testing library. Gospec provides a similar and familiar developer experience.

The language is a fairly thin layer on top of the standard Go tests. The goal of the project is to use the underlying testing mechanism set by the Go standard library and only provide a means to define tests (specs and features) for a better BDD testing experience.

Gospec has zero dependencies and the intention is for this to not change.

## Usage

The intended usage is to have the standard Go tests, defined like so `func TestMyFeature(t *testing.T) {`, then use the gospec library to initialize a `Suite` or `FeatureSuite` depending on which API is preferred.

<details>
    <summary>Test suite example</summary>

```go
import (
    "github.com/slavsan/gospec"
    ...
)

func TestCartSpec(t *testing.T) {
	gospec.WithSpecSuite(t, func(s *gospec.SpecSuite) {
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
```
</details>

<details>
    <summary>Test suite example with parallel execution</summary>

```go
import (
    "github.com/slavsan/gospec"
    ...
)

func TestCartSpec(t *testing.T) {
	gospec.WithSpecSuite(t, func(s *gospec.SpecSuite) {
        describe, beforeEach, it := s.ParallelAPI()

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
                            assert.Equal(w.T, []string{"Gopher Toy", "Lizard toy"}, w.Get("cart"))
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
                            assert.Equal(w.T, []string{}, w.Get("cart"))
                        })
                    })
                })
            })
        })
    })
}
```
</details>

<details>
    <summary>Feature Test suite example</summary>

```go
import (
	"github.com/slavsan/gospec"
	...
)

func TestCartFeature(t *testing.T) {
	gospec.WithFeatureSuite(t, func(s *gospec.FeatureSuite) {
        feature, background, scenario, given, when, then, _ :=
            s.With(gospec.WithPrintedFilenames()).API()

        feature("Cart", func() {
            var cart []string

            background(func() {
                given("there is a cart with three items", func(t *testing.T) {
                    cart = []string{
                        "Gopher Toy",
                        "Crab Toy",
                    }
                })
            })

            scenario("cart updates", func() {
                given("a new item has already been added", func(t *testing.T) {
                    cart = append(cart, "Lizard toy")
                })
                when("we remove the second item", func(t *testing.T) {
                    cart = []string{cart[0], cart[2]}
                })
                then("the cart should contain the correct two items", func(t *testing.T) {
                    assert.Equal(t, []string{"Gopher Toy", "Lizard toy"}, cart)
                })
            })

            scenario("removing items from the cart", func() {
                given("the second item has already been removed", func(t *testing.T) {
                    cart = cart[:1]
                })
                when("we remove the first item", func(t *testing.T) {
                    cart = cart[:0]
                })
                then("the cart should contain 0 items", func(t *testing.T) {
                    assert.Equal(t, []string{}, cart)
                })
            })
        })
    })
}
```
</details>

<details>
    <summary>Feature Test suite example with parallel execution</summary>

```go
import (
	"github.com/slavsan/gospec"
	...
)

func TestCartFeature(t *testing.T) {
	gospec.WithFeatureSuite(t, func(s *gospec.FeatureSuite) {
        feature, background, scenario, given, when, then := s.ParallelAPI()

        feature("Cart", func() {
            background(func() {
                given("there is a cart with three items", func(t *testing.T, w *gospec.World) {
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
                    assert.Equal(w.T, []string{"Gopher Toy", "Lizard toy"}, w.Get("cart"))
                })
            })

            scenario("removing items from the cart", func() {
                given("the second item has already been removed", func(t *testing.T, w *gospec.World) {
                    w.Swap("cart", func(cart any) any { return cart.([]string)[:1] })
                })
                when("we remove the first item", func(t *testing.T, w *gospec.World) {
                    w.Swap("cart", func(cart any) any { return cart.([]string)[:0] })
                })
                then("the cart should contain 0 items", func(t *testing.T, w *gospec.World) {
                    assert.Equal(w.T, []string{}, w.Get("cart"))
                })
            })
        })
    })
}
```
</details>

## Options

When initializing `SpecSuite` or `FeatureSuite` instances, there are several options to enhance the developer experience.

- `WithOutput` for specifying a custom output writer
- `WithPrintedFilenames` for printing the `filename:line` on each line where there is a `describe`, `beforeEach`, `it`, `feature`, `scenario` (or any other API method) defined. An editor which supports recognition of paths and allows for jumping to source, could improve the developer experience when dealing with maintenance of existing tests/specs.

### Parallel execution

Gospec supports parallel execution of tests.

One important caveat is that in order to do that, your tests would need to use the `ParallelAPI` (instead of the `API`) method, and will need to utilize the `World` struct which gets passed to all of your test functions, e.g. `beforeEach`, `it`, etc., or `given`, `when`, `then` in the case of `FeatureSpec` suites.

This is necessary since the way `gospec` works, it needs to store any contextual data used in one suite separate from the scope of other tests.
