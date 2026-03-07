# Test Anti-Patterns

## Testing Anti-Patterns

Avoid these common mistakes that reduce test effectiveness and maintainability.

## Test Absence

### No Tests for New Features

Adding functionality without tests is the most critical anti-pattern. Every new feature should include tests that verify its behavior.

### Only Testing Happy Path

Tests that only verify normal operation miss edge cases and error handling. Include tests for:
- Empty inputs
- Null or undefined values
- Boundary conditions
- Error conditions

## Test Quality Anti-Patterns

### Meaningless Test Names

Tests named after implementation details or generic names like "test1" provide no documentation value.

Bad: `testCalculate()`
Good: `shouldReturnZeroWhenInputIsEmpty()`

### Testing Implementation, Not Behavior

Tests that verify how code works rather than what it does break when implementation changes. This creates fragile tests that fail during legitimate refactoring.

### Multiple Responsibilities

A test that verifies multiple behaviors is difficult to debug and maintain. When it fails, it is unclear which behavior is broken.

### Hardcoded Test Data

Tests with magic numbers and complex setup that is not self-explanatory. Use meaningful, descriptive test data that explains itself.

## Test Structure Anti-Patterns

### Test Code Duplication

Duplicated test setup or assertion logic across tests indicates a need for test utilities or fixtures. Duplication makes maintenance difficult.

### Tight Coupling to Implementation

Tests that are tightly coupled to implementation details such as internal method calls, private methods, or specific variable names. These tests break easily when implementation changes.

### Ignoring Test Output

Tests that pass but produce warnings or errors. These may indicate incomplete tests or tests that are not actually verifying what they claim to verify.

### Tests That Depend on Each Other

Tests that must run in a specific order or depend on shared state from previous tests. This creates fragile test suites that can fail unpredictably.

## Test Maintenance Anti-Patterns

### Skipping Tests

Using @Ignore or commenting out tests. If a test is not relevant, delete it. If it is temporarily broken, fix it or create a tracking issue.

### Updating Tests Without Understanding

Changing assertions to make tests pass without understanding why they failed. This defeats the purpose of tests as verification.

### Slow Tests

Tests that take too long to run reduce the frequency of test execution. Keep tests fast to maintain the TDD feedback loop.

## Test Philosophy Anti-Patterns

### Tests as an Afterthought

Writing tests after production code is complete. This often results in poorly designed tests that verify the implementation rather than the behavior.

### 100% Coverage as a Goal

Pursuing 100% code coverage as a metric rather than meaningful test coverage. Coverage numbers do not indicate test quality.

### Tests as Documentation Only

Writing tests that serve only as documentation without verifying behavior. Tests must fail when behavior is broken to be valuable.

## What to Avoid

### Do Not Test Third-Party Code

Do not write tests that verify the behavior of external libraries or frameworks. Trust that they work and test your integration with them.

### Do Not Test Trivial Code

Getters, setters, and simple constructors do not need tests unless they contain logic.

### Do Not Over-Specify

Tests should verify behavior, not implementation details. Avoid asserting on internal state unless it is part of the documented contract.
