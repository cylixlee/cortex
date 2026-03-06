# Test Planning

## Planning Tests with Features

When planning a new feature or functionality, treat test planning as an integral part of the planning process. Tests are not an afterthought—they are part of the feature definition.

## When to Plan Tests

### During Feature Planning

Before writing any code:
1. Define the feature behavior
2. Identify the test cases that would verify this behavior
3. Write the tests first

This approach ensures:
- The feature is well-understood before implementation
- Edge cases are considered early
- The implementation is testable by design

### During Implementation

As you implement each piece of functionality:
1. Write a test that describes the next piece of behavior
2. Watch it fail
3. Write the minimum code to pass
4. Move to the next behavior

## How to Plan Tests

### Identify Behaviors, Not Functions

Think about what the system should do, not what functions should exist. Tests describe behaviors from the user's perspective.

Example instead of: "Test the calculateTotal function"
Think: "Test that the total is correctly computed from line items"

### List Test Cases First

Before writing code, list all test cases:
- Normal cases: standard inputs and expected outputs
- Edge cases: empty inputs, null values, boundary conditions
- Error cases: invalid inputs, error conditions

### Consider Test Coverage Strategically

Aim for meaningful coverage rather than maximum coverage:
- Cover happy path thoroughly
- Cover critical error paths
- Cover edge cases that are likely to occur
- Do not test trivial getters and setters unless they contain logic

## Test Planning in Agent Workflow

### Feature Addition Workflow

When adding a new feature:

1. Understand the feature requirements
2. Identify all behaviors the feature must have
3. For each behavior, identify the test cases
4. Write tests before writing production code
5. Implement to make tests pass
6. Refactor while keeping tests green

### Bug Fix Workflow

When fixing a bug:

1. Write a test that reproduces the bug
2. Verify the test fails
3. Fix the bug
4. Verify the test passes

This ensures the bug is documented by a test and cannot regress.

### Refactoring Workflow

Before refactoring:
1. Ensure test coverage exists for the code to be refactored
2. Run tests to establish a baseline
3. Refactor in small steps
4. Run tests after each step
5. All tests should continue to pass

## What to Test

### Test Behavior, Not Implementation

Tests should verify what the system does, not how it does it. This allows implementation to change without breaking tests.

### Test Through Public Interfaces

Test the public API, not internal implementation details. This maintains encapsulation and reduces test fragility.

### Test Boundary Conditions

Test the edges:
- Empty collections
- Zero, one, and many items
- Minimum and maximum values
- First and last items in a sequence

### Test Error Handling

Test that errors are handled appropriately:
- Invalid inputs produce appropriate errors
- Required fields are validated
- Exceptions are thrown when expected
