
# Tool Calls in Practice

A tool call is a multi-step exchange between the LLM and the [Session](../source/types/Session.go).
Tool names are always `namespace.Method` (e.g. `files.Read`, `agents.Await`).

## 1. Assistant requests a Tool Call

The model returns an assistant message carrying one or more `tool_calls`:

```json
{
    "role":    "assistant",
    "content": "Let me read the file first",
    "tool_calls": [{
        "id":   "call_0",
        "type": "function",
        "function": {
            "name":      "files.Read",
            "arguments": {
                "path": "./utils/Parse.go"
            }
        }
    }]
}
```

`arguments` is a JSON object whose keys match the tool's parameter schema. The
`id` links the tool result back to this call.

## 2. Session executes the Tool Call

[Session.ReceiveChatResponse](../source/types/Session.go) parses each
`tool_calls` entry, resolves the namespace via `Session.GetTool`, and calls
`Tool.Call(method, arguments)`. Tool implementations live in
[source/tools](../source/tools/).

## 3. Session appends the Tool result

The result (or an `Error: ...` string) is appended as a `tool` message carrying
the same `tool_call_id` and `tool_name`:

```json
{
    "role":         "tool",
    "content":      "files.Read: ./utils/Parse.go\npackage utils\n\nfunc Parse(...) { ... }",
    "tool_call_id": "call_0",
    "tool_name":    "files.Read"
}
```

## 4. Assistant responds with a reasoned message

```json
{
    "role":    "assistant",
    "content": "There is a bug at line 123. I propose a fix."
}
```

## 5. User responds

```json
{
    "role":    "user",
    "content": "Yes, please fix it for me."
}
```

## Blocking calls

`agents.Await` blocks inside `CallTool` until the hired agent finishes, so the
next inference happens only after the agent's work report is available. Other
tools return immediately.
