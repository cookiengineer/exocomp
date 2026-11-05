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

Early returns are still acceptable for impossible states, validation that must
happen before any useful work, or when nesting would be absurd. The default is
positive-first nesting.

## Formatting

- Use **tabs** for indentation in Go and JSON files. YAML uses two spaces.
- Opening braces stay on the same line as the statement.
- Put a blank line after an opening brace and before the matching closing brace
  for non-trivial blocks.
- Separate loops, `if` blocks, variable declarations and logical sections with
  blank lines.
- Do not introduce column alignment unless it already exists in the surrounding
  block.
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

Project imports come first, followed by standard library imports.

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
```

Use an alias when the package name is generic, ambiguous or would collide with
an existing identifier:

```go
import net_http "net/http"
import net_url "net/url"
import utils_fmt "exocomp/utils/fmt"
```

Blank imports are reserved for embedded assets:

```go
import _ "embed"
```

## Naming

- Exported types and methods use `CamelCase`.
- Unexported helper functions use `camelCase`.
- Local variables use `snake_case`.
- Acronyms follow the existing local convention rather than being rewritten
  just for style.

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
- Keep error messages deterministic and short.
- Do not return raw `os`, `exec`, `net/http` or `json` errors through tool
  results. Use the sanitizing helpers instead. See
  [ARCHITECTURE.md](./ARCHITECTURE.md#tool-error-message-sandboxing).
- Do not log an error and silently continue unless that behaviour is explicitly
  intended.
- Prefer `fmt.Errorf` with a stable message over wrapping raw host details.

```go
if err == nil {

	return result, nil

} else {
	return "", fmt.Errorf("agents.Hire: Cannot start agent \"%s\".", name)
}
```

## Functions and Helpers

- Keep functions focused on one responsibility.
- Split large tools into helper files named after the helper, for example
  `sanitizeSandboxPath.go`, `resolveSandboxPath.go`, `getAgentWorkReport.go`.
- File names for primary types match the type name, for example `Agent.go`,
  `Config.go`, `Session.go`.
- Prefer small helper functions over deeply nested inline logic, but keep the
  positive/error branch shape intact when extracting them.
- Return the successful value from the deepest positive branch, and return the
  error from the corresponding `else`.

## Comments

Comments are used sparingly and explain *why*, not *what*.

- `// NOTE:` documents an important design decision, lifecycle rule or
  non-obvious constraint.
- `// TODO:` marks planned work.
- `// XXX:` marks a known workaround or debugging-only section.

Do not add comments that merely restate the code.

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

## Tooling

Before finishing a change:

```bash
cd source;
gofmt -w <changed-go-files>;
go test ./...;
go vet ./...;
```

Do not reformat unrelated files. A diff should remain focused on the change
being made.
