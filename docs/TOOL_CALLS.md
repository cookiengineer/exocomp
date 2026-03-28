
# Tool Calls in Practice

## 1. Assistant requests Tool Calls

```json
{
    "role":       "assistant",
    "message":    "Let me read the file first",
    "tool_calls": [{
        "type": "function",
        "function": {
            "name": "files.read",
            "arguments": {
                "path": "main.go"
            }
        }
    }]
}
```

## 2. Client parses Tool Calls and executes them

// TODO

## 3. Client responds with Tool Calls, confirming execution

```json
{
    "role":       "tool",
    "message":    "package main\n\nfunc main() { ... }",
    "tool_name":  "files.read",

    // MUST BE the same tool calls as before
    "tool_calls": [{
        "type": "function",
        "function": {
            "name": "files.read",
            "arguments": {
                "path": "main.go"
            }
        }
    }]
}
```

## 4. Assistant responds with reasoned message

```json
{
    "role":    "assistant",
    "message": "There is a bug somewhere at line 123. I propose a fix."
}
```

## 5. User responds with message

```json
{
    "role":    "user",
    "message": "Yes, please fix it for me."
}
```
