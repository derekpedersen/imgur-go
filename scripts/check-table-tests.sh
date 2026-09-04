#!/usr/bin/env bash
set -euo pipefail

# Enforce table-driven test style on changed Go test files.
# A test file with Test* functions must include both []struct and t.Run usage.

has_pattern() {
  local pattern="$1"
  local file="$2"

  if command -v rg >/dev/null 2>&1; then
    rg -q "${pattern}" "${file}"
    return $?
  fi

  grep -Eq "${pattern}" "${file}"
}

resolve_base_ref() {
  if [[ -n "${TABLE_TEST_BASE_REF:-}" ]]; then
    echo "${TABLE_TEST_BASE_REF}"
    return 0
  fi

  if git rev-parse --verify origin/master >/dev/null 2>&1; then
    git merge-base origin/master HEAD
    return 0
  fi

  if git rev-parse --verify HEAD~1 >/dev/null 2>&1; then
    echo "HEAD~1"
    return 0
  fi

  return 1
}

if ! base_ref="$(resolve_base_ref)"; then
  echo "table-test-style: unable to determine base ref; skipping check"
  exit 0
fi

test_files=()
while IFS= read -r f; do
  test_files+=("${f}")
done < <(git diff --name-only --diff-filter=ACM "${base_ref}"...HEAD -- '*_test.go')

if [[ ${#test_files[@]} -eq 0 ]]; then
  echo "table-test-style: no changed test files"
  exit 0
fi

failures=0

for f in "${test_files[@]}"; do
  if [[ ! -f "${f}" ]]; then
    continue
  fi

  if ! has_pattern 'func[[:space:]]+Test' "${f}"; then
    continue
  fi

  has_table=0
  has_subtest=0

  if has_pattern '\[\][[:space:]]*struct[[:space:]]*\{' "${f}"; then
    has_table=1
  fi

  if has_pattern 't\.Run[[:space:]]*\(' "${f}"; then
    has_subtest=1
  fi

  if [[ ${has_table} -eq 0 || ${has_subtest} -eq 0 ]]; then
    echo "table-test-style: ${f} is not table-driven (requires []struct and t.Run)"
    failures=$((failures + 1))
  fi
done

if [[ ${failures} -gt 0 ]]; then
  echo "table-test-style: failed with ${failures} file(s)"
  exit 1
fi

echo "table-test-style: passed"
