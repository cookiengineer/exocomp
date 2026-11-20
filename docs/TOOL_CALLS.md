# Tool Calls

A tool call is an exchange between the model and
[engine.Session](../source/engine/Session.go). Tool names are always
`namespace.Method` (for example `files.Read` or `agents.Await`).

## 1. Assistant Requests a Tool Call

The model returns an assistant message with one or more `tool_calls`:

```json
{
    "role": "assistant",
    "content": "Let me read the file first",
    "tool_calls": [{
        "id": "call_0",
        "type": "function",
        "function": {
            "name": "files.Read",
            "arguments": {
                "path": "./utils/Parse.go"
            }
        }
    }]
}
```

`arguments` is a JSON object whose keys match the tool schema. `id` links the
result back to this call.

## 2. Session Executes the Call

[Session.ReceiveChatResponse](../source/engine/Session.go) passes each entry to
`CallTool`. `CallTool` resolves the namespace through `Session.GetTool`, checks
`HasMethod`, and calls `Tool.Call(method, arguments)`.

Tool implementations live below [`source/tools`](../source/tools/).

## 3. Session Appends the Tool Result

The result or an `Error: ...` string is appended as a `tool` message with the
same `tool_call_id` and `tool_name`:

```json
{
    "role": "tool",
    "content": "files.Read: File \"./utils/Parse.go\" contents\npackage utils\n\nfunc Parse(...) { ... }",
    "tool_call_id": "call_0",
    "tool_name": "files.Read"
}
```

A failure uses a short, deterministic message:

```text
files.Read: File "./x.go" does not exist.
```

## 4. Assistant Continues

The model receives the tool message and responds with its next action or
answer. A user message can follow at any point.

## Blocking Calls

These calls block inside `CallTool` until an external event resolves them:

| Tool | Waits for |
|:-----|:----------|
| `agents.Await` | The hired agent to finish and return a work report. |
| `humans.Ask` | A free-form answer. |
| `humans.Choose` | One selected choice. |
| `humans.AskPermission` | A capability decision. |

All other tools return immediately.

## Sanitized Tool Errors

Tool errors are sanitized before they reach the model. See
[Tool Error Message Sandboxing](./ARCHITECTURE.md#tool-error-message-sandboxing).
A failing tool uses a relative path and a fixed message instead of the raw Go
error, so absolute paths, temp directories, executable paths and DNS/IP details
are not exposed.
