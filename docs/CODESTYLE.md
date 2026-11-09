# Code Style

This document describes the coding conventions used throughout Exocomp. It
covers the general Go code style. Test-specific conventions live in
[TESTING_GUIDE.md](./TESTING_GUIDE.md).

The goal is not novelty. New code should look like it was written by the same
person who wrote the surrounding file.

## Positive Action First

The most important rule is branch ordering:

- The **positive branch** in terms of the intended logical action goes first
  and is nested.
- The `} else {` branches contain the error, failure or fallback cases.
- Do not flatten a success path into a chain of guard clauses when the same
  logic can be expressed as a nested positive branch with error branches in
  the `else` blocks.

In other words, write the thing the function is supposed to do first, and put
the failure handling below it.

```go
func (client *Client) OnSwitch(role string) (string, error) {

	if client.Session != nil {

		tool := client.Session.GetTool("agents.Switch")

		if tool != nil {

			agents_tool, ok := tool.(*tools.Agents)

			if ok == true {

				// Positive path continues here.

				return "agents.Switch: Switched to agent role.", nil

			} else {
				return "", fmt.Errorf("agents.Switch: Tool does not exist.")
			}

		} else {
			return "", fmt.Errorf("agents.Switch: Tool does not exist.")
		}

	} else {
		return "", fmt.Errorf("agents.Switch: No session available.")
	}

}
```

Avoid writing the same function as a series of early error returns:

```go
// Do not do this when a positive nested branch is possible.

if client.Session == nil {
	return "", fmt.Errorf("agents.Switch: No session available.")
}

tool := client.Session.GetTool("agents.Switch")

if tool == nil {
	return "", fmt.Errorf("agents.Switch: Tool does not exist.")
}

// ...
```

The success value is returned from the deepest positive branch. The matching
error is returned from the corresponding `else`.

`} else if {` chains are part of this shape and are allowed.

Early returns are permitted only for impossible states (nil receivers, missing
map entries that cannot occur), validation that must happen before any useful
work, or when nesting would exceed three levels. They are not a substitute for
the positive-first shape.

## Formatting

- Use **tabs** for indentation in Go and JSON files. YAML uses two spaces.
- Opening braces stay on the same line as the statement.
- A block is **non-trivial** when it contains more than one statement. For
  non-trivial blocks, put a blank line after the opening brace and before the
  matching closing brace. Single-statement blocks may stay compact.
- Separate loops, `if` blocks, variable declarations and logical sections with
  blank lines.
- Do not introduce manual column alignment. Struct field alignment is whatever
  `gofmt` produces.
- Keep function signatures and calls readable; put long argument lists one
  argument per line.

```go
if ok == true {

	result := doSomething(
		first_argument,
		second_argument,
	)

	return result, nil

} else {
	return "", fmt.Errorf("package.Method: Nothing to do.")
}
```

## Imports

Use one import per line. Do not use parenthesized import blocks.

Project imports come first, followed by standard library imports, in one
block. Third-party imports form a separate block after one blank line. Sort
alphabetically within each group.

```go
package tools

import "exocomp/agents"
import "exocomp/schemas"
import "exocomp/types"
import utils_chat "exocomp/utils/chat"
import utils_fmt "exocomp/utils/fmt"
import "fmt"
import net_url "net/url"
import "os"
import "strings"

import html_tree "golang.org/x/net/html"
```

Use an alias when the package name is generic, ambiguous or would collide with
an existing identifier. Aliases use the package's parent path prefix:

```go
import net_http "net/http"
import net_url "net/url"
import utils_fmt "exocomp/utils/fmt"
import utils_ast "exocomp/utils/ast"
import html_tree "golang.org/x/net/html"
```

In test files, `testing` is always the last standard-library import.

Blank imports are reserved for embedded assets:

```go
import _ "embed"
```

For non-test code this document is authoritative over
[TESTING_GUIDE.md](./TESTING_GUIDE.md).

## Naming

- Exported types and methods use `CamelCase`.
- Unexported helper functions use `camelCase`. Do not use `snake_case`.
- Local variables use `snake_case`.
- Acronyms follow the existing local convention rather than being rewritten
  just for style.
- Go source file names use the same spelling as the identifier they contain.
  Underscores are reserved for test files (`_test`), platform suffixes
  (`_GOOS`, `_GOARCH`) and generated assets.

Constructors and parsers are named `New<Type>` and `Parse<Type>`. Converters
that return a type are named `<Type>From<Source>`, never bare
`From<Source>`:

```go
NewConfig(...)
ParseConfig(...)
ColorFromInts(...)   // not FromInts(...)
AgentFromMessage(...)
```

Type names are spelled-out words. Do not collapse a sequence of words into
initials: use `type Color struct`, not `type RGBA struct`. Established
initialisms are allowed when they are conventional and unambiguous: `URL`,
`HTTP`, `ID`, `IP`, `API`, `JSON`, `YAML`, `ASCII`, `ANSI`, `TTY`, `UI`,
`OS`, `CPU`.

Use one receiver name per type across all of its methods.

```go
func (tool *Agents) GetAgent(name string) *types.Agent {

	sanitized_name := utils_fmt.FormatAgentName(name)
	agents_tool, ok := tool.contents[sanitized_name]

	if ok == true {
		return agents_tool
	} else {
		return nil
	}

}
```

Use numbered suffixes when a function returns multiple related values:

```go
result1, err1 := tool.First()
result2, err2 := tool.Second()

if err1 == nil && err2 == nil {
	// ...
}
```

Boolean variables should read as predicates:

```go
is_running := false
has_backup := recovery.HasBackup()
```

## Booleans and Conditions

Always compare booleans explicitly. Do not use `!value` or rely on a bare
boolean value.

```go
if ok == true { ... }
if ok != true { ... }
if has_backup == false { ... }
```

Use `err == nil` for the success path:

```go
result, err := tool.Call(method, arguments)

if err == nil {
	return result, nil
} else {
	return "", fmt.Errorf("tools.Call: %s", err.Error())
}
```

For type assertions and map lookups, keep the `ok` value explicit:

```go
value, ok := arguments["name"].(string)

if ok == true {
	// ...
} else {
	// ...
}
```

## Error Handling

- Error cases belong in `} else {` branches.
- Error messages must be prefixed with the tool or component name, for example
  `agents.Switch: ...`, `files.Read: ...`, `Session.SwitchAgent: ...`.
- The component prefix is mandatory, including for shared helpers. When the
  component is a parameter, build the prefix from it instead of hard-coding a
  different name:

  ```go
  return "", fmt.Errorf("%s.%s: Invalid path \"%s\".", tool, method, path)
  ```

- Never return a bare lowercase message such as `invalid command ...`.
- Keep error messages deterministic and short.
- Do not return raw `os`, `exec`, `net/http` or `json` errors through tool
  results. Use the sanitizing helpers instead. See
  [ARCHITECTURE.md](./ARCHITECTURE.md#tool-error-message-sandboxing).
- Do not log an error and silently continue unless that behaviour is explicitly
  intended.
- Prefer `fmt.Errorf` with a stable message over wrapping raw host details.
- Use `%w` only when a caller must inspect the cause with `errors.Is` or
  `errors.As`. Otherwise use `%s` with a stable message.

```go
if err == nil {

	return result, nil

} else {
	return "", fmt.Errorf("agents.Hire: Cannot start agent \"%s\".", name)
}
```

## Package and File Layout

- One primary type per file, named after that type: `type Config struct` lives
  in `Config.go`.
- All methods of a type live in that type's file, whether the receiver is a
  pointer or a value.
- Constructors and parsers are part of the type they build and live in the
  returned type's file: `NewConfig`, `ParseConfig` and `ConfigFromYAML` all
  live in `Config.go`.
- Every other free function lives in its own file named exactly after the
  function: `func FormatAgentName(...)` lives in `FormatAgentName.go`.
- Every other type lives in its own file named exactly after the type:
  `type Pricing struct` lives in `Pricing.go`.
- Unexported implementation-detail types may live in the file of the type that
  owns them, for example a private `agentProcess` struct in `Agents.go`. If a
  helper type is shared by more than one primary type, give it its own
  lowercase file.
- Unexported helpers that operate on types from the same package use a
  lowercase name and a matching lowercase file, for example
  `getAgentWorkReport.go`.
- Helpers that operate on standard library types belong in `./source/utils`
  (for example `utils/fmt`, `utils/strings`, `utils/fs`), not next to the
  caller.
- When a package grows multiple related implementations (`Foo`, `Bar`, ...),
  group them into a subpackage, each in its own `Foo.go` and `Bar.go`, not as
  loose files beside each other.
- Test fixtures under `testdata/` and small helpers inside `cmds/*/main.go`
  are exempt from the one-implementation-per-file rule.
- Platform-specific files keep their `_GOOS` or `_GOARCH` suffix, for example
  `CheckRoot_linux.go`.
- Do not create a `utils/<name>` package that mirrors or duplicates an existing
  top-level package. Shared helpers have exactly one home.
- `init()` is allowed for package-level tables and registries. Keep it in the
  file that owns the data, ideally one `init()` per file.

```text
source/types/Config.go        // type Config, NewConfig, ParseConfig, methods
source/types/Pricing.go       // type Pricing
source/utils/fmt/FormatAgentName.go
source/utils/tty/Color.go     // type Color, NewColor, ColorFromInts, methods
source/tools/Agents.go        // type Agents, methods, private agentProcess
```

## Functions and Helpers

- Keep functions focused on one responsibility.
- Prefer small helper functions over deeply nested inline logic, but keep the
  positive/error branch shape intact when extracting them.
- Return the successful value from the deepest positive branch, and return the
  error from the corresponding `else`.
- Do not leave a helper in a file merely because that is where it was first
  used. Place it according to the rules below.
- Do not duplicate a helper. Extend the existing one instead.

When adding a helper, decide where it belongs:

1. If it is a method, keep it in the receiver's type file.
2. If it is `New<Type>`, `Parse<Type>` or `<Type>From<Source>`, keep it in the
   returned type's file.
3. If it operates on standard library types or is reusable across packages, put
   it in the matching `utils` package.
4. Otherwise give it its own lowercase file named after the function.

## Comments

Comments are used sparingly and explain *why*, not *what*.

- Doc comments on exported identifiers are allowed and should describe
  behaviour.
- `// NOTE:` documents an important design decision, lifecycle rule or
  non-obvious constraint.
- `// TODO:` marks planned work.
- `// XXX:` marks a known workaround or debugging-only section.

Do not add comments that merely restate the code, such as `// key`,
`// Skip '='` or `// body`. When a structural section marker carries design
intent, convert it into a `// NOTE:`.

## Data and Serialization

- JSON uses tabs for indentation.
- Keep struct tags consistent with the existing type. Some fields use
  `snake_case`, others use hyphenated keys; match the surrounding file.
- Timestamps use the existing `schemas.Datetime` helpers.
- Prefer `json.Marshal` / `json.Unmarshal` and the project parsers over ad-hoc
  string formatting.

## Concurrency

- Guard shared maps and slices with the existing mutexes.
- Keep lock scopes small.
- Do not call user callbacks or spawn subprocesses while holding a mutex.
- Treat `GetMessages`, `GetTools` and similar accessors as views whose locking
  rules are defined by the owning type.

## Tests

Test-specific conventions are documented in
[TESTING_GUIDE.md](./TESTING_GUIDE.md). The same branch-order rule applies to
tests: put the expected case first and keep failure assertions in the `else`
branch.

```go
if err == nil {
	// expected path
} else {
	t.Errorf("Expected %v to be nil", err)
}
```

Shared test helpers use `camelCase` and live in the relevant `_test.go` file.
They are exempt from the one-implementation-per-file rule.

## Tooling

Before finishing a change:

```bash
cd source;
gofmt -w <changed-go-files>;
go build ./...;
go test ./...;
go vet ./...;
```

Do not reformat unrelated files. A diff should remain focused on the change
being made.
