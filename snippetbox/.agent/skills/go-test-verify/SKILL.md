---
name: go-test-verify
description: Automatically run and verify tests before finalizing any code changes.
---
# Testing & Verification Protocol

Before submitting any code changes for review:
1. **Execute Tests**: Run `go test -v ./...` to ensure no regressions were introduced.
2. **Auto-Fix**: If tests fail, analyze the output, fix the root cause, and re-run until all tests pass.
3. **Lint & Format**: Run `go fmt ./...` and `go vet ./...` to maintain clean, idiomatic code.
4. **Final Confirmation**: Only present the solution to the user after confirming that the code is both functional and verified by the test suite.