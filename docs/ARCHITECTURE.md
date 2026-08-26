
# Architecture

Exocomp is an agentic environment for LLM agents that is optimized for Go
packages. The core idea is that LLM agents are spawned as subprocesses in
sandboxed parts of the codebase, and communicate through explicit tool calls
instead of freeform markdown documents.

Everything a running agent knows is expressed as a `symbol` of the codebase,
i.e. a `"path/to/package"` plus a `Symbol` (a struct, function, method or
variable with its signature).

Sandboxes and packages are equivalent: an agent reworks one subpackage in
isolation while other agents work on other packages without conflicting.

## Process Hierarchy and Sandboxing

Communication between the human, the planner and the contractor agents is a
strict subprocess tree:

```
planner        (exocomp web|terminal, playground == sandbox == project root)
|-> architect  (exocomp agent, sandbox == project root)
|-> coder      (exocomp agent, sandbox == project/utils)
|-> tester     (exocomp agent, sandbox == project/utils)
|-> coder      (exocomp agent, sandbox == project/cmds/fibonacci)
```

- A child's `--playground` is always its parent's `--sandbox`.
- A child's `--sandbox` is the package subfolder it is allowed to touch.
- Because every agent is a real OS process with a distinct working directory,
  the OS process tree mirrors the agent sandboxes. Firewalls, EDR tooling and
  monitoring tools can observe what each subprocess is accessing by looking at
  its CWD and file descriptors.

## Module Layout

```
source/
  agents/        Role templates (*.yaml) and prompt/message rendering
  parsers/
    html/        HTML to markdown/text conversion
    markdown/    Markdown parser for the website frontend
  tools/         Tool implementations and JSON schemas
  types/         Basic types
  schemas/       Wire types (Message, ToolCall, ChatRequest/Response, ...)
  ui/
    jsonl/       Child-process frontend with JSONL over stdin/stdout
    tty/         Interactive terminal frontend (the human's planner)
    web/         Web frontend + REST API (the human's planner)
  cmds/
    exocomp/         UI dispatcher (agent | terminal | web | webview)
    exocomp-agent/   Child-process entry point (exocomp agent)
    exocomp-web/     Web server entry point
```

## JSONL Protocol

A hired agent runs `exocomp agent` and talks to its parent exclusively through
line-delimited JSON on `stdout`. The `stderr` is free-form logging. Every line
is formatted as `<prefix>:<json>\n`:

| Prefix               | Payload                                               | Producer           | Consumer          |
|:---------------------|:------------------------------------------------------|:-------------------|:------------------|
| `schemas.Message`    | [schemas.Message](../source/schemas/Message.go)       | child renderer     | parent reader     |
| `types.ContextUsage` | [types.ContextUsage](../source/types/ContextUsage.go) | child heartbeat    | parent reader     |

The child emits `schemas.Message:` lines for every system/assistant/tool message
it produces. The parent appends them to the child's `types.Agent.Messages` so
that the planner can inspect and summarize the child's work later.

The `types.ContextUsage:` line doubles as a liveness heartbeat.

## Agent Lifecycle

Agent state is tracked on [types.Agent](../source/types/Agent.go) via the
`Status` field:

| Status     | Meaning                                         |
|:-----------|:------------------------------------------------|
| `working`  | Hired, subprocess running                       |
| `finished` | Exited `0` with an `agents.Quit` work report    |
| `failed`   | Exited without a work report (non-zero, killed) |
| `fired`    | Killed via `agents.Fire`                        |

Lifecycle is owned by the [Agents tool](../source/tools/Agents.go):

**agents.Hire(role, prompt, sandbox)**:

1. generates a name
2. resolves the sandbox
3. spawns `exocomp agent --name ... --prompt ... --sandbox ...`
4. records the `*types.Agent` as `working` and starts a single reader goroutine.

**Reader**:

One goroutine owns the child's stdout.

- reads lines unbounded, no `bufio.Scanner` token cap
- appends `schemas.Message`s and updates `ContextUsage`/liveness
- On `EOF` it calls `cmd.Wait()` correctly after all reads, finalizes `Status` to `finished` or `failed`,
  then closes the process `done` channel and removes the process from the running map.

**agents.Await(name)**:

1. blocks on the process `done` channel until the child finishes
2. returns the `agents.Quit: Work Report` (or a terminal error).

It deliberately does not retry: polling would make the planner fill its
context window with identical "still working" messages.

**agents.Fire(name)**:

1. cancels the child context and waits for `done`
2. marks the agent `fired`.

**agents.Quit(message)**:

This is called by the process child. The report is flushed synchronously
to stdout before the process exits, so the parent reliably observes it.

A liveness watchdog cancels a child that stops producing output for longer
than `IdleTimeout` (default `5m`) and a hard `Timeout` (default `30m`) bounds
the total lifetime. Both are fields on `Agents` with defaults set in `NewAgents`.

In addition to `Status`, `types.Agent` carries `StartedAt` and `FinishedAt`
timestamps so the schedule view can show how long each agent ran.

## Tool Lifecycle

Each Agent role declares `allowed-tools` and `allowed-programs` in its YAML
`agents/*.yaml` files. At startup the frontend builds a `tools.Toolset` containing
only the allowed tools and their JSON schemas, and registers them on the `Session`
keyed by their namespace (the part before the first `.`).

- Tool schema names are `namespace.Method` (e.g. `files.Read`, `requirements.DefineFunc`, `agents.Hire`).
- `Session.CallTool` resolves the namespace, calls `Tool.Call(method, arguments)` and appends the result (or `Error: ...`) as a `tool` message.
- `Tool.Get(id)` is the symbol lookup used by `skills` and other meta tools.
- Sandbox-bound tools (files, programs, requirements, bugs, changelog) resolve all paths through `resolveSandboxPath`/`sanitizeSandboxPath` so a contractor cannot escape its sandbox.

Each tool's LLM-facing contract lives in its embedded JSON schema at `tools/*.json`
with precise parameter descriptions and syntax examples.

Role separation is enforced by *capability*, not by the prompt alone.

An architect has `files.Read` but not `files.Write` and must record its output
with `requirements.Define*`. A coder gets `files.Write`, `changelog.*` and
`programs.Execute`. A tester only writes unit tests and `bugs.*` reports.
... and so on.

## Session Loop

[types.Session](../source/types/Session.go) drives a single agent's
request/response loop:

1. `SendChatRequest` appends the user message and calls
   `infer_chat_completions` (one HTTP POST to `/v1/chat/completions`).
2. `ReceiveChatResponse` appends the assistant message; if it contains
   `tool_calls`, each is executed via `CallTool` (blocking for `agents.Await`)
   and then the loop recurses into `infer_chat_completions`.
3. `Recovery` snapshots the session and agents to
   `<playground>/.exocomp/` so a planner can be restored across restarts.
