# Testing Guide

How to write tests. See [TESTING.md](./TESTING.md) for prerequisites and live
server workflows, and [CODESTYLE.md](./CODESTYLE.md) for the linter-enforced
style.

Run the unit tests from the repository root:

```bash
make test;
```

`make test` runs `go -C source test ./...`. The toolchain module tests run
with `make test-toolchain`. Use the same Go command directly for a single
package or test:

```bash
go -C source test ./tools/files -run 'TestFiles_Read' -v;
```

`make lint` excludes `_test.go` files. `go -C toolchain run ./lint -tests` is
available, but the source test suite is not fully lint-clean yet.

## Principles

- White-box tests: use `package tools`, `package ast`, `package types`, and so
  on, never an external `_test` package.
- Every method in a tool's `Call()` dispatch has at least one test. Add the
  test in the same change that adds the method.
- Tests are deterministic and offline. The `//go:build agents` suite is the
  only exception; it starts a local llama-server and a real child process.
- Cover the expected path and the failure cases: invalid arguments, missing
  symbols, sandbox escapes, not-found and permission denied.

## File and Function Naming

- One `_test.go` file per source file: `Files.go` -> `Files_test.go`.
- Larger tool packages split tests by concern:
  `Agents_test.go`, `Agents_lifecycle_test.go`, `Agents_quit_test.go`,
  `Agents_main_test.go`, `Agents_list_test.go`.
- Test functions are named `Test<Type>_<Method>`, optionally with a scenario
  suffix:

  ```go
  func TestFiles_Read(t *testing.T)             { ... }
  func TestFiles_ReadSymbol(t *testing.T)       { ... }
  func TestBugs_Add(t *testing.T)               { ... }
  func TestAgents_Inquire(t *testing.T)         { ... }
  func TestPrograms_ExecuteWithoutPermission(t *testing.T) { ... }
  func TestHumans_Call_ArgumentValidation(t *testing.T)    { ... }
  ```

- Shared test helpers use camelCase: `newTestAgents`, `hireTestAgent`,
  `buildFakeAgent`, `createTestSkill`.

## Imports

Follow [CODESTYLE.md](./CODESTYLE.md). In test files, `testing` is the last
standard-library import.

## Setup and Teardown

Use the real constructor with an explicit method allow-list, a temporary
playground and a sandbox below it. Preserve the playground on failure:

```go
func TestFiles_Write(t *testing.T) {

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-files-*")
	sandbox := filepath.Join(playground, "files")
	tool := NewFiles([]string{"Copy", "List", "Read", "Stat", "Write"}, playground, sandbox)

	if tool != nil {

		result, err := tool.Write("./file.txt", "This is the file content!")

		if strings.Contains(result, "files.Write: File \"./file.txt\"") == false {
			t.Errorf("Expected file to be written")
		}

		if err != nil {
			t.Errorf("Expected %v to be nil", err)
		}

	} else {
		t.Errorf("Expected tool to be not nil")
	}

	t.Cleanup(func() {

		if t.Failed() == true {
			t.Logf("Preserving folder %s for debugging.", playground)
		} else {
			os.RemoveAll(playground)
		}

	})

}
```

Conventions:

- `os.MkdirTemp("/tmp", "exocomp-test-<tool>-*")` for the playground.
- `filepath.Join(playground, "<tool>")` for the sandbox.
- Guard the body with `if tool != nil { ... } else { t.Errorf(...) }`.
- Use `t.Cleanup` to preserve the playground on failure and remove it on
  success.
- Create fixture files before calling the constructor for tools that read state
  at boot time (`skills`, `requirements`).

## Assertions

- Use `t.Fatalf` only when later assertions would be meaningless (the tool is
  `nil`, a required fixture is missing, a blocking call timed out). Use
  `t.Errorf` for ordinary checks.
- Prefer `strings.Contains` over exact equality for free-form result text, so
  tests do not break on reworded messages:

```go
if strings.Contains(result, "Attempt to escape sandbox") == false {
	t.Errorf("Expected sandbox escape to be detected")
}
```

- Use `%q` for quoted strings, `%v` for errors and values, and `%d` for counts.
- Check the error message, not just that an error exists:

```go
if err == nil {
	t.Errorf("Expected %v to be not nil", err)
} else if strings.Contains(err.Error(), "has no Symbol") == false {
	t.Errorf("Expected missing symbol error, got %v", err)
}
```

## Table-Driven Tests

Use a table when the same assertion runs over many inputs (parsers, type
resolution). Declare the table inside the test:

```go
func TestGetSymbolType_BasicTypes(t *testing.T) {

	tests := []struct {
		source string
		symbol string
		want   string
	}{
		{"package dummy\ntype MyEnum string", "MyEnum", "string"},
		{"package dummy\ntype IPv4 [4]byte", "IPv4", "[4]byte"},
	}

	for _, test := range tests {

		result, err := GetSymbolType([]byte(test.source), test.symbol)

		if err != nil {
			t.Errorf("Expected %v to be nil, got %v", nil, err)
		}

		if result != test.want {
			t.Errorf("Expected %q, got %q", test.want, result)
		}

	}

}
```

## Sandbox Escape Tests

Any tool that resolves paths (`files.*`, `programs.*`, `skills.*`, `bugs.*`,
`changelog.*`, `requirements.*`) must test traversal and absolute paths outside
the sandbox:

```go
for _, path := range []string{"./../../../file.txt", "/etc/passwd", "./..\\..\\../file.txt"} {

	result, err := tool.Read(path)

	if err != nil {
		if strings.Contains(err.Error(), "Attempt to escape sandbox") == false {
			t.Errorf("Expected sandbox escape to be detected, got %v", err)
		}
	} else {
		t.Errorf("Expected %v to be not nil, got %q", err, result)
	}

}
```

Repeat the same check for `Write` with the same paths.

## Blocking and Concurrency Tools

Tools whose methods block (`humans.Ask`, `humans.Choose`, `agents.Await`) use a
goroutine plus a polling helper, bounded by a deadline:

```go
result_ch := make(chan string, 1)
err_ch := make(chan error, 1)

go func() {
	result, err := tool.Ask("What is your name?")
	result_ch <- result
	err_ch <- err
}()

waitForQuestion(t, tool, 1*time.Second)

select {
case result := <-result_ch:
	t.Fatalf("Expected Ask to block until answered, got %s", result)
case <-time.After(50 * time.Millisecond):
	// Still blocked, as expected.
}
```

Always set an explicit timeout so a regression cannot hang the suite.

## Fake Binaries and Build Tags

Tests that would otherwise need a real external binary use a fake:

- [`source/tools/testdata/fakeagent/main.go`](../source/tools/testdata/fakeagent/main.go)
  prints `schemas.Message:` JSON lines and exits. It is driven by
  `EXOCOMP_FAKE_SCENARIO` (`quit`, `large`, `noquit`, `fail`, `hang`,
  `summarize`, `input`, `args`).
- [`source/tools/testdata/fakepodman/main.go`](../source/tools/testdata/fakepodman/main.go)
  records arguments through `EXOCOMP_FAKE_PODMAN_ARGS`.
- [`source/tools/testdata/fakeshell/main.go`](../source/tools/testdata/fakeshell/main.go)
  provides a deterministic shell for PTY tests.

Helpers live in [`source/tools/testutil`](../source/tools/testutil):
`BuildFakeAgent`, `BuildFakePodman`, `BuildFakeShell`, `BuildSandboxExec`,
`RequireSandbox` and `CreateTestSymlink`. The `Build*` helpers build once
through `sync.Once` and return the binary path; callers set the matching
`EXOCOMP_*` variable. `RequireSandbox` also sets `EXOCOMP_SANDBOX`.

Tests that need a live llama-server are fenced behind `//go:build agents` as the
first line of the file, so `make test` stays offline. See
[TESTING.md](./TESTING.md) for that suite.
