# TDD Principles

## The Three Laws of TDD

1. **Write no production code except to make a failing test pass.**
2. **Write no more of a test than is sufficient to fail.**
3. **Write no more production code than is sufficient to make the failing test pass.**

These laws form a tight feedback loop: test first, watch it fail, write minimal code to pass, then refactor.

## Red-Green-Refactor Cycle

### Red: Write a Failing Test

Write a test that describes the behavior you want. The test should fail because the behavior does not yet exist. This step forces you to think about the interface before implementation.

### Green: Make the Test Pass

Write the minimum production code necessary to make the test pass. Do not optimize, do not add features beyond what the test requires. Speed matters more than elegance in this phase.

### Refactor: Clean Up

Once tests pass, clean up the code while keeping tests green. Remove duplication, improve naming, simplify structure. Refactoring is only safe because tests verify behavior is preserved.

## Test Structure Principles

### Arrange-Act-Assert

Structure each test in three clear phases:
- **Arrange**: Set up the objects and conditions for the test
- **Act**: Execute the behavior being tested
- **Assert**: Verify the outcome matches expectations

### One Assertion Per Test (Recommended)

Each test should verify one behavior. Multiple assertions are acceptable when they test different aspects of the same behavior, but avoid testing unrelated behaviors in a single test.

### Test Names Describe Behavior

Name tests after the behavior they verify, not after implementation details. The test name should read as a specification: "should do X when Y".

## Test Scope

### Unit Tests

Test individual functions, methods, or classes in isolation. The scope is small and the test should run fast.

### Integration Tests

Test how components work together. This includes testing interactions with databases, file systems, or external services.

### When to Use Each

- Use unit tests for business logic, algorithms, and data transformations
- Use integration tests for workflows that span multiple components
- Prefer unit tests for speed and isolation
- Use integration tests sparingly and only when necessary

## What Makes a Good Test

### Fast

Tests should run in milliseconds. Slow tests reduce the feedback loop benefit of TDD.

### Independent

Tests should not depend on each other. Each test can run in isolation and in any order.

### Repeatable

Tests should produce the same result every time. No randomness, no dependency on external state.

### Self-Verifying

Tests should clearly pass or fail. No manual interpretation required.

### Timely

Write tests before production code. This ensures testability and drives cleaner design.
