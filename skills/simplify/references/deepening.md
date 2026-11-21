# Deepening

Read this when looking for shallow modules or reshaping one in fab:simplify.

## Words

- **Module:** anything with an interface and an implementation, at any scale: a function, a type, a package.
- **Interface:** everything a caller must know to use it correctly: signature, invariants, ordering, errors,
  required setup, performance.
- **Deep / shallow:** deep means a lot of behavior behind a small interface; shallow means the interface is nearly as
  complex as what it hides.
- **Seam:** where an interface lives, and where behavior can be swapped without editing that place.
- **Adapter:** a concrete thing that fills a seam (production code, an in-memory fake).
- **Leverage / locality:** what depth buys: callers get more per thing learned; changes and bugs stay in one place.

## Tests for a candidate

- **Deletion test:** imagine deleting the module. If its complexity vanishes, it was a pass-through (inline it). If
  it reappears across its callers, it earns its keep.
- **The interface is the test surface:** tests that reach past the interface to check internals mean the module is
  the wrong shape.
- **One adapter is a hypothetical seam, two is a real one:** don't add a port or interface type unless something
  actually varies across it (production plus a test fake counts).

## Dependencies decide the shape

1. **In-process** (pure logic, in-memory state): merge and test through the new interface directly.
2. **Local stand-in exists** (an in-memory database, a temp filesystem): keep the stand-in in the tests; the seam
   stays internal.
3. **Your own remote service:** a port at the seam, a network adapter in production, an in-memory one in tests.
4. **Third-party service:** inject it as a port; tests use a fake.

## Moving tests

Replace, don't layer: write tests at the new interface, then delete the old tests of the shallow parts once the new
ones cover their behavior. Tests should survive internal refactors; one that must change when only the inside
changes was testing past the interface.

## Sketching alternatives

When the shape isn't obvious, give two or three designers (fresh sub-agents if available) the same brief (files,
dependencies and their category, what sits behind the seam) and a different aim each: the smallest interface (one
to three entry points); the most common caller's case made trivial; ports and adapters for remote dependencies.
Each returns the interface, a caller's usage, what it hides, and its trade-offs. Compare them on depth, locality,
and seam placement, then recommend one (or a mix).

Adapted from Matt Pocock's `codebase-design` and `improve-codebase-architecture` skills
([mattpocock/skills](https://github.com/mattpocock/skills), MIT), which draw on John Ousterhout's *A Philosophy of
Software Design* and Michael Feathers' seams.
