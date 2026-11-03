# Cshargextcap Agent Test Execution Rules

These rules apply to all coding agents working in this repository.

## Structured Logging

All `log/slog` calls must pass structured attributes instead of alternating key/value arguments. Use the appropriate `slog` attribute constructor, such as `slog.String("rootpath", rootpath)`, `slog.Int("count", count)`, or `slog.Any("value", value)`. Always pass errors as `xslog.Error(err)` from `github.com/thediveo/nonstd/xslog`.

## Deferred Calls

When a `defer` statement calls a function that returns one or more values (such as an `error`), the return value(s) must be explicitly discarded rather than implicitly dropped, to satisfy linting rules. Wrap the call in a closure and assign the result(s) to `_`, for example use `defer func() { _ = session.Close() }()` instead of `defer session.Close()`.

When registering session cleanup in Ginkgo test specs, use the compact form `DeferCleanup(sess.Close, ctx)`.

## Sets

For sets, use `sets.Set[T]` from `github.com/thediveo/nonstd/sets` instead of `map[T]struct{}` or `map[T]bool`. Create sets with `sets.New(...)`, then use `Add`, `Delete`, `Contains`, and `Elements`. A nil set is safe for `Contains`, so it can mean "no filter" or an empty exclusion list.

## Package Documentation

Every package's doc comment (the `/* Package <name> ... */` comment preceding the `package` clause) must live in its own dedicated `doc.go` file containing nothing but that comment and the `package` clause. Do not attach package doc comments to files that also contain other declarations.

## Code Formatting

When formatting Go code, use the custom formatting tool configured in `gostwire.code-workspace`:

```bash
golangci-lint fmt <path/to/file.go>
```

Do not use plain `gofmt` as `golangci-lint fmt` enforces import grouping (`gci`) and other repository formatting rules.

## Test Framework

Unit tests must be written using Ginkgo v2 (`github.com/onsi/ginkgo/v2`) and Gomega (`github.com/onsi/gomega`), not the plain `testing` package style.

Each package with tests must have a separate `package_test.go` file that contains the Ginkgo/Gomega test suite plumbing (the `Test<Package>` function calling `RegisterFailHandler` and `RunSpecs`), for example:

```go
package foo

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFoo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "cshargextcap/internal/foo package")
}
```

Actual specs (`Describe`/`It` blocks) belong in separate `*_test.go` files, not in `package_test.go`. Test helpers shared by the specs of a package (fixture loaders, builders, constants, and the like) belong in `package_test.go` after the suite plumbing, instead of in separate helper files.

Each `*_test.go` file containing specs must use exactly one top-level `Describe`; put additional `Describe`, `Context`, and `When` groups inside it rather than declaring multiple top-level `Describe` blocks in the same file.

When writing Ginkgo specs (`Describe`, `Context`, `When`, `It`, `DescribeTable`, etc.), always space all specs from each other by a single empty line, including the ends of different spec levels.

Where a value comes from a call returning `(T, error)` and the error is only used to fail the test, prefer wrapping the call inline with `Successful(...)` from `github.com/thediveo/success` instead of assigning the error to a variable and asserting `Expect(err).NotTo(HaveOccurred())` separately. For example, use `Expect(foobar(...)).To(HaveField("Path", Successful(filepath.Abs("..."))))` rather than a separate `path, err := filepath.Abs(...)` plus error check. For calls returning `(T1, T2, error)` or `(T1, T2, T3, error)`, use `Successful2R(...)` or `Successful3R(...)` respectively.

Where a value comes from a function call or map lookup returning `(T, ok bool)` (or `(T, bool)`), prefer wrapping the call inline with `Allright(...)` from `github.com/thediveo/success` instead of receiving `(value, ok)` and asserting `Expect(ok).To(BeTrue())` separately. For example, use `textContent := Allright(dict.Lookup("key"))` rather than a separate `Expect(ok).To(BeTrue())` assertion.

Where a value needs to be type-asserted (typecast) to a target type `T`, prefer wrapping the expression inline with `AssignableTo[T](...)` from `github.com/thediveo/success` instead of performing `v, ok := val.(T)` and asserting `Expect(ok).To(BeTrue())` separately. For example, use `textContent := AssignableTo[*mcp.TextContent](result.Content[0])`.

Where none of the non-error return values are needed and only the error matters, prefer passing the call directly into `Expect(...).Error()` instead of assigning it to a discarded variable; this also covers expected errors from calls with additional return values. For example, use `Expect(server.Connect(ctx, transport, nil)).Error().NotTo(HaveOccurred())` for a successful call or `Expect(iedDetails(ctx, nil, fstest.MapFS{})).Error().To(HaveOccurred())` for an expected failure rather than assigning the other return values and checking the error separately.

In tests, prefer slice and collection helpers from `github.com/thediveo/nonstd` (such as `xslices.Map`, `xslices.Filter`, `xslices.FindFuncOk`, etc.) over manual loops to avoid low-level test assembly noise and maintain a useful level of abstraction.

Never add test code to explicitly remove or kill stale containers, as container cleanup is automatically handled by `morbyd`'s `session.WithAutoCleaning(...)`.

For temporary directories in Ginkgo specs, do not use `GinkgoT().TempDir()`; instead use the plain stdlib together with `DeferCleanup`, for example `tmpdir := Successful(os.MkdirTemp("", "mobydav-test-"))` followed by `DeferCleanup(os.RemoveAll, tmpdir)`.

## End-to-End Tests

End-to-end tests live in the `e2e/` package only and are opt-in: its suite skips itself unless the Ginkgo label filter explicitly selects the `e2e` label. Do not use build tags to gate them. `RunSpecs` already applies the `e2e` label suite-wide; use further labels (such as `k8s` or `tshark`) on specs to allow narrower selections.

Run them explicitly, keeping the usual constraints:

```bash
go test -v -race -p 1 ./e2e/... -ginkgo.label-filter=e2e
```

## Required Test Matrix

When validating changes, run relevant tests in both privilege modes:

1. Non-root execution.
2. Root execution.

Do not treat one as a substitute for the other. Some suites are gated by root checks and only execute meaningful coverage when run as root.

## No Parallel Test Suite Execution

Test suites must not be run in parallel.

1. Disable package-level parallelism with `-p 1`.
2. Avoid commands that execute multiple package test suites concurrently.
3. Prefer sequential package runs when targeting multiple packages.

## Required Command Style

Use these forms unless there is a documented reason not to:

```bash
go test -race -p 1 ./...
go test -race -p 1 -exec 'sudo -n' ./...
```

For focused package runs, keep the same constraints:

```bash
go test -v -race -p 1 ./path/to/package
go test -v -race -p 1 -exec 'sudo -n' ./path/to/package
```

## Validation Expectation

A change is not fully validated until:

1. Non-root tests pass (or expected skips are understood).
2. Root-mode tests pass.
3. Runs were done with the race detector enabled.
4. Runs were done without parallel suite execution.