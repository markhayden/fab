# Writing tests

Read this when build writes a task's checks, and when test judges weak tests.

## A test worth keeping

- Its name says what the caller gets ("checkout with a valid cart is confirmed"), not how the code does it.
- One behavior per test, checked the way a caller would, so it survives any change to the inside.

## Anti-patterns

- **Coupled to the inside:** mocks the module's own collaborators, calls private functions, asserts call counts or
  order, or checks a side channel (querying the database) instead of the interface. Tell: it breaks when the code is
  refactored without any change in behavior.
- **Can't fail:** the expected value is derived the same way the code derives it (`add(a, b) == a + b`), a snapshot
  produced by the code under test, a constant compared with itself.

## Mocking

Mock only at the system's edges: third-party services (payments, email), the clock, randomness, and sometimes the
network or file system. Prefer a real local stand-in (an in-memory or temporary database, a temp directory) to a
mock. Never mock code the project owns.

Adapted from Matt Pocock's `tdd` skill and its `tests.md` and `mocking.md`
([mattpocock/skills](https://github.com/mattpocock/skills), MIT).
