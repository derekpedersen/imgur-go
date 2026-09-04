# Copilot Instructions for imgur-go

## Goal
This repository follows idiomatic Go patterns. Keep code simple, explicit, and package-focused.

## Core Rules
- Prefer concrete types over interface-plus-implementation pairs when there is only one implementation.
- Accept interfaces at boundaries only when there is a clear consumer-side need.
- Return errors instead of logging from service packages; callers own logging policy.
- Keep constructors explicit and validate required dependencies at creation time.
- Use one request path via the shared imgur client for API calls.
- Avoid pointer-to-scalar return values unless nil is semantically meaningful.

## Testing Rules
- Use table-driven tests as the default style in Go test files.
- Use subtests with `t.Run` for each case.
- Name test cases clearly and concisely.
- Keep one assertion focus per test case when practical.
- Integration tests must remain environment-gated.

## Test Naming and Structure
- Test function names should describe behavior, not implementation details.
- Arrange test cases as slices of structs with named fields.
- Include both success and failure cases for non-trivial behavior.
- Prefer helper functions for repeated setup.

## Review Checklist for Copilot Changes
- Is the API surface idiomatic and minimal?
- Are errors wrapped with enough context?
- Are tests table-driven with clear case names?
- Are integration tests gated and unit tests deterministic?
