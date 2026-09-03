---
description: "Use when creating or editing Go tests; enforce table-driven tests and subtest structure."
applyTo: "**/*_test.go"
---

# Go Test Table Style

When writing or updating tests in this repository, use table-driven tests by default.

## Required Pattern
- Define test cases as a slice of structs.
- Include a `name` field for each case.
- Iterate cases with `for _, tc := range tests` and call `t.Run(tc.name, func(t *testing.T) { ... })`.
- Rebind loop variable inside the loop (`tc := tc`) before `t.Run`.

## Required Coverage
- Include at least one success case and one failure case for behavior with branching.
- Validate error behavior explicitly:
- on expected error: assert error is non-nil and message/type matches intent.
- on expected success: assert error is nil and result fields match expected values.

## Style
- Keep fixtures minimal and local to the test function unless shared setup is significant.
- Prefer `t.Helper()` in setup helpers.
- Keep assertions deterministic and avoid time/network dependence in unit tests.

## Example Template
```go
func TestSomething(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "valid input", input: "abc", want: "ABC", wantErr: false},
		{name: "empty input", input: "", want: "", wantErr: true},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got, err := Something(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("want %q, got %q", tc.want, got)
			}
		})
	}
}
```
