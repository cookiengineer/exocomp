
# Testing Guide

This document captures the unit test conventions used across the Exocomp code
base. It exists so that new tests can be written to the same standard as the
existing suite, with no second-guessing about style, naming or structure.

See [TESTING.md](./TESTING.md) for how to *run* the tests (including the
live-LLM and llama-server workflows). This guide is about how to *write* them.

## Principles

- **White-box tests.** Tests live in the same package as the code under test
  (`package tools`, `package ast`, `package types`, ...), never in a `_test`
  package. This lets tests reach unexported helpers and fields when needed.
- **One `Call()` method, one test.** Every method reachable through a tool's
  `Call()` dispatch must have at least one test. When a new method is added to
  a `Call()` switch, add its test in the same change.
- **Deterministic and offline.** Tests must not need a network, a live LLM or
  external services. The only exception is the `//go:build agents` suite, which
  uses a fake agent binary and (optionally) a local llama-server.
- **Happy path and error path.** Each test exercises both the expected result
  and the failure cases (missing symbol, invalid arguments, sandbox escape,
  not-found, permission denied, ...).

## File Layout and Naming

- One `_test.go` file per source file, named after it:
  `Files.go` -> `Files_test.go`, `GetSymbol.go` -> `GetSymbol_test.go`.
- Larger tools split their tests across several files by concern:
  `Agents_test.go`, `Agents_lifecycle_test.go`, `Agents_quit_test.go`,
  `Agents_main_test.go`, `Agents_list_test.go`.
- Test functions are named `Test<Type>_<Method>`, optionally with a scenario
  suffix:

  ```go
  func TestFiles_Read(t *testing.T)            { ... }
  func TestFiles_ReadSymbol(t *testing.T)      { ... }
  func TestBugs_Add(t *testing.T)              { ... }
  func TestAgents_Inquire(t *testing.T)        { ... }
  func TestPrograms_ExecuteWithoutPermission(t *testing.T) { ... }
  func TestHumans_Call_ArgumentValidation(t *testing.T)    { ... }
  ```

- Shared setup helpers in test files use `camelCase`-free lowercase names,
  e.g. `newTestAgents`, `hireTestAgent`, `buildFakeAgent`, `createTestSkill`.

## Package and Imports

One import per line, sorted alphabetically, with an alias where a package name
would otherwise collide or be ambiguous:

```go
package tools

import "exocomp/types"
import utils_ast "exocomp/utils/ast"
import net_url "net/url"
import "os"
import "path/filepath"
import "strings"
import "testing"
```

`testing` is always the last import.

## Code Style

The code base uses a deliberately spacious, verbose style. Follow it exactly:

- **Blank line after every statement block.** Loops, `if` branches and variable
  declarations are separated by blank lines, inside and outside of functions.
- **Explicit boolean comparisons.** Always write `== true`, `!= true`,
  `== false`, never `!value` or bare `value`:

  ```go
  if ok != true { ... }
  if strings.Contains(result, "foo") == false { ... }
  if err == nil { ... }
  ```

- **`snake_case` variable names**, never camelCase:

  ```go
  from_path, ok1 := arguments["from_path"].(string)
  lines3 := strings.Split(result3, "\n")
  ```

- **Tabs for indentation**, four-space alignment is not used.

## Setup and Teardown

Sandboxed tools are exercised against a real temporary directory. Create it,
derive a sandbox from it, and clean up on completion while preserving it for
debugging on failure:

```go
func TestFiles_Write(t *testing.T) {

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-files-*")
	sandbox       := filepath.Join(playground, "files")
	tool          := NewFiles([]string{"Copy", "List", "Read", "Stat", "Write"}, playground, sandbox)

	if tool != nil {

		result1, err1 := tool.Write("./file.txt", "This is the file content!")

		if result1 != "files.Write: File \"./file.txt\" with 26 B written." {
			t.Errorf("Expected file to be written")
		}

		if err1 != nil {
			t.Errorf("Expected %v to be nil", err1)
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

Conventions used above:

- `os.MkdirTemp("/tmp", "exocomp-test-<tool>-*")` for the playground.
- `filepath.Join(playground, "<tool>")` for the sandbox.
- Guard the whole body with `if tool != nil { ... } else { t.Errorf(...) }`.
- Use the real constructor (`NewFiles`, `NewBugs`, ...) with an explicit
  method allow-list, never a hand-built struct.
- Wrap every result in a `t.Cleanup` that logs the temp folder on failure and
  removes it on success.

For tools that need files on disk (Skills, Requirements, Changelog, Bugs),
build the fixture files inside the playground *before* constructing the tool
(because the constructor reads state at boot time).

## Assertions

- Use `t.Fatalf` only to stop the test when later assertions would be
  meaningless (tool is `nil`, a required fixture is missing, a blocking call
  timed out). Use `t.Errorf` for ordinary checks.
- Prefer `strings.Contains` over exact equality for free-form result text, so
  tests don't break on reworded messages:

  ```go
  if strings.Contains(result, "Attempt to escape sandbox") == false {
  	t.Errorf("Expected %v to detect attempt to escape sandbox", err)
  }
  ```

- Use `%q` for quoted strings, `%v` for errors and values, `%d` for counts.
- Check the error message, not just that an error exists:

  ```go
  if err == nil {
  	t.Errorf("Expected %v to be not nil", err)
  } else if strings.Contains(err.Error(), "has no Symbol") == false {
  	t.Errorf("Expected missing symbol error, got %v", err)
  }
  ```

## Table-Driven Tests

Where the same assertion runs over many inputs (parsers, type resolution),
use a table instead of copy-pasted bodies. The table is an anonymous struct
slice declared right inside the test:

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
`changelog.*`, `requirements.*`) must include escape attempts in its tests.
Use traversal and absolute paths outside the sandbox and assert the error:

```go
result, err := tool.Read("./../../../file.txt")
result, err := tool.Read("/etc/passwd")
result, err := tool.Write("./..\\..\\../file.txt", "content")

if err != nil {
	if strings.Contains(err.Error(), "Attempt to escape sandbox") == false {
		t.Errorf("Expected %v to detect attempt to escape sandbox", err)
	}
} else {
	t.Errorf("Expected %v to be not nil", err)
}
```

## Blocking and Concurrency Tools

Tools whose methods block (`humans.Ask`/`humans.Choose`, `agents.Await`)
are tested with a goroutine plus a polling helper, and bounded by a deadline:

```go
result_ch := make(chan string, 1)
err_ch    := make(chan error, 1)

go func() {
	result, err := tool.Ask("What is your name?")
	result_ch <- result
	err_ch    <- err
}()

id := waitForQuestion(t, tool, 1*time.Second)

select {
case result := <-result_ch:
	t.Fatalf("Expected Ask to block until answered, but it returned early: %s", result)
case <-time.After(50 * time.Millisecond):
	// still blocked, expected
}
```

Always set an explicit timeout so a regression cannot hang the suite forever.

## Fake Binaries and Build Tags

Tests that would otherwise require spawning a real agent use a fake binary:

- `source/tools/testdata/fakeagent/main.go` is a tiny `main` that prints
  `schemas.Message:` JSON lines and exits, driven by the `EXOCOMP_FAKE_SCENARIO`
  environment variable.
- `buildFakeAgent` (in `Agents_lifecycle_test.go`) builds it once via
  `sync.Once` and `newTestAgents` injects it through `t.Setenv("EXOCOMP_AGENT", ...)`.
- Add new scenarios to `fakeagent/main.go` when a test needs a different
  behaviour (e.g. `summarize` for `agents.Inquire`).

Tests that need a live llama-server are fenced behind `//go:build agents` as
the first line of the file, so the default `go test ./...` run stays offline.

## Coverage Checklist

For each tool, confirm every method in its `Call()` switch is tested:

- `files`: `Copy`, `List`, `Read`, `ReadSymbol`, `Search`, `Stat`, `Write`, `WriteSymbol`
- `requirements`: `List`, `DefineFunc`, `DefineInterface`, `DefineStruct`, `DefineType`, `Search`, `Signoff`
- `bugs`: `List`, `Add`, `Fix`, `Search`
- `changelog`: `Add`, `Change`, `Deprecate`, `Fix`, `List`, `Remove`, `Search`
- `agents`: `Await`, `List`, `Roles`, `Hire`, `Fire`, `Inquire`, `Quit`
- `humans`: `Ask`, `Choose`, `Answer`
- `programs`: `List`, `Execute`, `Stat`
- `skills`: `List`, `Load`, `Unload`, `Execute`
- `websites`: `Fetch`, `List`, `Stat`

The `utils/ast` package follows the same rules for its exported API:
`GetSymbol`, `GetSymbolType`, `GetPackageSymbols`, `ReadSymbol`, `WriteSymbol`,
`HasSymbol` and `SearchSymbols` each have tests covering `func`, `interface`, `struct` and the
basic Go types (`uint8`, `string`, `[]byte`, maps, ...).

## Running the Tests

```bash
cd /path/to/exocomp/source;

# Run everything (no LLM required)
go test ./...;

# Run a single package
go test ./tools/;

# Run a specific test
go test ./tools/ -run 'TestFiles_Read' -v;

# Vet before you're done
go vet ./...;
```
