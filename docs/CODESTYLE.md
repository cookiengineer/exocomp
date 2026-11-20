# Code Style

The codestyle linter is the executable source of truth. This document is an
overview; the rules live in
[`toolchain/internal/utils/lint`](../toolchain/internal/utils/lint/).
[`toolchain/lint/main.go`](../toolchain/lint/main.go) is the CLI wrapper and owns
flag parsing for the linter.

Run the checks from the repository root. The targets are defined in the
[`Makefile`](../Makefile):

```bash
make lint;
make lint-toolchain;
make test;
make test-toolchain;
```

`make lint` checks non-test files below `source/`. `make lint-toolchain`
checks the `toolchain/` module. Use the linter directly for individual rules,
test files or rule listings:

```bash
go -C toolchain run ./lint -list;
go -C toolchain run ./lint -rules layout/method-file;
go -C toolchain run ./lint -tests;
go -C toolchain run ./lint -root .;
```

`-tests` currently reports existing drift in the source test tree. `make lint`
and `make lint-toolchain` are the required gates. `-root .` targets the
toolchain module.

Flag parsing for every toolchain command lives in that command's `main.go`.
Internal packages receive already-parsed values and do not import `flag` or
read `os.Args`. Command implementations live under
`toolchain/internal/utils/<name>/`; command folders contain only `main.go` and
CLI-only files.

A diagnostic is suppressed with `//nolint:<rule>` on the same line or the line
above. Comma-separated rules and `//nolint:all` are supported.

## Rule Overview

Rule IDs and summaries are registered in
[`rules.go`](../toolchain/internal/utils/lint/rules.go). The table below links
to the implementation files.

| Area | Implementation | Intent |
|:-----|:---------------|:-------|
| Branches | [`checkGuard.go`](../toolchain/internal/utils/lint/checkGuard.go) | Put the positive action first and keep error or fallback cases in `else`. Consecutive guards become one `else if` chain. |
| Conditions | [`checkBareCondition.go`](../toolchain/internal/utils/lint/checkBareCondition.go) | Compare booleans explicitly (`== true`, `== false`, `err == nil`); do not use `!` negation. |
| Formatting | [`checkGofmt.go`](../toolchain/internal/utils/lint/checkGofmt.go) | Follow `gofmt`; separate multi-statement blocks with blank lines; put a blank line above every `case` and `default`; use tabs in JSON. |
| Imports | [`checkImportParen.go`](../toolchain/internal/utils/lint/checkImportParen.go) | One import per line. Project imports first, then the standard library in the same block, then a blank line and third-party imports. Sort each group. |
| Naming | [`checkFileName.go`](../toolchain/internal/utils/lint/checkFileName.go) | Exported identifiers are CamelCase, unexported functions camelCase, locals and parameters snake_case, booleans read as predicates, and a type uses one receiver name. |
| Errors | [`checkErrorPrefix.go`](../toolchain/internal/utils/lint/checkErrorPrefix.go) | Prefix every error with `component.Method:` and start the message with a capital letter. |
| Layout | [`checkOneType.go`](../toolchain/internal/utils/lint/checkOneType.go) | One exported type per file; methods and constructors in the owner's file; every other free function in its own file named after it; at most one `init()` per file. |
| Dependencies | [`checkToolIsolation.go`](../toolchain/internal/utils/lint/checkToolIsolation.go) | Tool packages do not import each other and only import allowed project packages. `types` does not import `utils/capabilities`; `utils` does not import `tools`. |
| Tests | [`checkTestFile.go`](../toolchain/internal/utils/lint/checkTestFile.go) | Test files are named after a source file; test functions follow `Test<Type>_<Method>`. |

The linter only checks what is implemented in those files. Run
`go -C toolchain run ./lint -list` for the current rule IDs and summaries.

## Positive Branch Example

```go
func (client *Client) OnSwitch(role string) (string, error) {

	if client.Session != nil {

		// ...

		return "agents.Switch: Switched to agent role.", nil

	} else {
		return "", fmt.Errorf("agents.Switch: No session available.")
	}

}
```

## Conventions Not Enforced by the Linter

- Comments explain why, not what. Use `// NOTE:`, `// TODO:` and `// XXX:`.
- Keep struct tags consistent with the surrounding type. Use
  `schemas.Datetime` for timestamps and the project parsers instead of ad-hoc
  string formatting.
- Guard shared maps and slices with the existing mutexes. Keep lock scopes
  small and do not call user callbacks or spawn subprocesses while holding a
  lock.
- Put shared helpers in the package that owns the behavior. See
  [Utility Ownership](./ARCHITECTURE.md#utility-ownership).
- Do not reformat unrelated files. Keep a diff focused on the change.

Test-specific conventions live in [TESTING_GUIDE.md](./TESTING_GUIDE.md).
