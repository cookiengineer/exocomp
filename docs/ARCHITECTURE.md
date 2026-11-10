
# Architecture

Exocomp is an agentic environment for LLM agents that is optimized for Go
packages. The core idea is that LLM agents are spawned as subprocesses in
sandboxed parts of the codebase, and communicate through explicit tool calls
instead of freeform markdown documents.

Everything a running agent knows is expressed as a `symbol` of the codebase,
i.e. a `"path/to/package"` plus a `Symbol` (a function, method, struct,
interface or named type such as `string`, `uint8` or `[]byte`).

Sandboxes and packages are equivalent: an agent reworks one subpackage in
isolation while other agents work on other packages without conflicting.

## Process Hierarchy and Sandboxing

Communication between the human, the planner and the contractor agents is a
strict subprocess tree:

```
planner        (exocomp-web web | exocomp terminal, playground == sandbox == project root)
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

### Symlink Policy

Path-bound tools reject symbolic links that could resolve outside the sandbox.
`resolveSandboxPath`/`sanitizeSandboxPath` validate every existing path
component with `os.Lstat` and refuse symlinked components. Directory listings,
searches and folder copies skip symlink entries entirely.

`programs.Execute` applies the same validation to path-like arguments.

### Process Confinement (Landlock / Seatbelt)

`programs.Execute` and `skills.Execute` run their subprocesses through the
`exocomp-sandbox` helper. The helper is spawned with the target program and
applies the platform sandbox before executing it:

- On Linux, a [landlock](https://docs.kernel.org/userspace-api/landlock.html)
  ruleset is applied before `execve`.
- On macOS, an SBPL profile is generated and applied through
  `/usr/bin/sandbox-exec -p`.

Both backends share the same fixed filesystem and process policy:

- The sandbox folder is the only writable location.
- A read-only allowlist covers system binaries, libraries, configuration and
  the resolved target program directory.
- Basic device nodes (`/dev/null`, `/dev/zero`, `/dev/urandom`, ...) are
  readable and writable.
- The read/write policy is fixed. There is no `off` switch and no fallback: if
  the platform sandbox or the helper is unavailable, execution fails with a
  fixed error.

On Linux, `PR_SET_NO_NEW_PRIVS` is set before restricting, which also blocks
setuid privilege escalation. Because landlock is Linux-only, the Linux backend
refuses to run on other operating systems. Network access is not restricted by
landlock, and `/proc` and `/sys` are readable for compatibility.

On macOS, the generated SBPL profile is deny-by-default and denies network
access. Paths are resolved through `filepath.EvalSymlinks` before being emitted
because the kernel matches resolved vnodes. Both the original and resolved forms
are included, which covers `/tmp` -> `/private/tmp`, `/var` -> `/private/var`
and `/etc` -> `/private/etc`. The profile allows only the configured file,
process, sysctl, signal, Mach lookup and POSIX semaphore operations.
`sandbox-exec` is deprecated since macOS 10.8 but remains available and is used
by Bazel, Nix, Homebrew, Claude Code and Codex. Using `sandbox-exec` instead of
the cgo `sandbox_init` API keeps the whole project buildable with
`CGO_ENABLED=0`.

HOME, temp, runtime and Go cache directories are redirected into
`<sandbox>/.exocomp/landlock/` on Linux and `<sandbox>/.exocomp/seatbelt/` on
macOS so ordinary tools such as `go`, `gofmt` and `gcc` can write without
leaving the sandbox. Credential and session sockets inherited from the agent
environment (`SSH_AUTH_SOCK`, `GPG_AGENT_INFO`, `DBUS_SESSION_BUS_ADDRESS`,
`KRB5CCNAME`, `DISPLAY`, ...) are removed from the child environment.
`.exocomp` is hidden from `files.List` and ignored by the Go toolchain.

The helper is discovered next to the running `exocomp` binary (or via
`EXOCOMP_SANDBOX` for development). The installer places it in the same `bin`
folder as `exocomp` and `exocomp-agent`.

Known limitations: on Linux, landlock does not restrict network access, `/proc`
and `/sys` are readable for compatibility, and metadata changes, ptrace, signals
and already-open file descriptors are not mediated. On macOS, the seatbelt
profile denies network access by default, but the platform still does not
mediate every operation and pre-opened file descriptors remain usable. On both
platforms, programs that expect host configuration or caches outside the sandbox
may need sandbox copies.

### Container Execution (Podman)

`containers.*` shells out to an externally installed `podman` binary. It does
not use `exocomp-sandbox`: podman needs its rootless socket, runtime directories,
user namespace helpers and network access, which the process sandbox does not
grant. The go code stays buildable with `CGO_ENABLED=0` because there is no
podman Go binding and no cgo dependency. On Windows the tool works with a podman
machine installation, even though the host process sandbox is unavailable there.

Container isolation is enforced by the tool itself:

- Every container gets the label `io.exocomp.sandbox=<sha256(sandbox)[:16]>`.
  `List` filters on this label and every lifecycle or exec operation verifies it
  through `podman inspect`, so an agent cannot reach unrelated user containers.
- The sandbox is mounted at `/workspace` and that is the default workdir. No
  arbitrary volumes, port publishing, capabilities, privileged mode or host
  namespaces are accepted from the model.
- Images must appear in the role's `allowed-images` list. An empty list denies
  all images.
- The default network mode is `bridge`; `none` is opt-in and `host` is never
  allowed.
- Podman output that reaches the model goes through
  `utils/podman`'s `sanitizeError`, so host paths and raw runtime errors are not
  leaked.

## Module Layout

```
source/
  agents/        Role templates (*.yaml) and prompt/message rendering
  parsers/
    html/        HTML to markdown/text conversion
    markdown/    Markdown parser for the website frontend
  tools/         Tool implementations and JSON schemas
  types/         Shared types (Agent, Config, Bug, ChangelogEntry, Requirement, ...)
  schemas/       Wire types (Message, ToolCall, ChatRequest/Response, ...)
  utils/         Shared helpers (ast, chat, cli, fmt, fs, http, api adapters, ...)
  ui/
    jsonl/       Child-process frontend with JSONL over stdin/stdout
    tty/         Interactive terminal frontend (the human's planner)
    web/         Web frontend + REST API (the human's planner)
  cmds/
    exocomp/         UI dispatcher (agent | config | terminal)
    exocomp-agent/   Child-process entry point (exocomp agent)
    exocomp-web/     Web server entry point (agent | config | web)
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

**agents.Switch(role)**:

Switches the main TTY session to a different preconfigured agent profile. It is
only wired in the TTY frontend; hired contractor subprocesses expose the tool
but receive an error when they call it. Before switching, the previous
conversation is summarized by a summarizer subprocess. The session then starts
fresh with the target profile's system prompt and that summary as its first
assistant message. Tools, programs and adapters are rebuilt from the target
profile so switching cannot retain capabilities from the previous role.

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
- `Session.CallTool` resolves the namespace, invokes the tool and appends the result (or `Error: ...`) as a `tool` message carrying the matching `tool_call_id`.
- Sandbox-bound tools (files, programs, requirements, bugs, changelog) resolve all paths through `resolveSandboxPath`/`sanitizeSandboxPath` so a contractor cannot escape its sandbox.
- `containers.*` mounts only the agent sandbox into podman containers and scopes every container to the sandbox label. `podman` must not be added to `allowed-programs`, otherwise `programs.Execute` could bypass these checks.

Each tool's LLM-facing contract lives in its embedded JSON schema at `tools/*.json`
with precise parameter descriptions and syntax examples.

Role separation is enforced by *capability*, not by the prompt alone.

An architect has `files.Read` but not `files.Write` and must record its output
with `requirements.Define*`. A coder gets `files.Write`, `changelog.*` and
`programs.Execute`. A tester only writes unit tests and `bugs.*` reports.
... and so on.

## Tool Error Message Sandboxing

Tool results and error messages are part of the LLM's context window. A
low-level Go error string can leak host information that a smaller "abliterated"
model (with fewer safety guardrails) might use to probe or escape the sandbox.
For example, `os.Stat`'s `*PathError` embeds the **absolute** path it was given:

```
files.Stat: stat /home/alice/projects/secret/...: no such file or directory
```

Because `resolveSandboxPath`/`sanitizeSandboxPath` resolve every path to an
absolute location internally, passing `err.Error()` straight back to the model
leaks the host's absolute directory layout. The same applies to
`os.Executable()` paths in `agents.Hire`, the `/tmp/...` temp dir in
`agents.Inquire`, and `net/http` dial errors in `websites.*` (which can expose
DNS names and IPs).

### Rules

1. Never return a raw `err.Error()` from `os`, `exec`, `net/http` or `json`
   through a tool's `Call(...)`. Map it to a fixed, path-safe message.
2. Always echo the **user-supplied or sanitized relative path** (`"./x"`),
   never the resolved absolute path.
3. Keep messages short, deterministic and prefixed `namespace.Method:` so a
   smaller LLM can parse them. Prefer `File "./x" does not exist.` over an
   ENOENT string.

### Helpers

| File | Purpose |
|:-----|:--------|
| `source/tools/sanitizeFilesystemError.go` | Maps `os.*` filesystem errors to `does not exist.` / `Permission denied.` / `Cannot access ...`, using the relative path. |
| `source/tools/sanitizeExecutionError.go` | Maps subprocess errors from `programs.Execute` / `skills.Execute` to `Permission denied.` / `Program doesn't exist.` / `exited with an error.` |

### How to add a tool

1. Resolve/sanitize every path via `resolveSandboxPath` / `sanitizeSandboxPath`.
2. On any `os.*` call failure, return
   `sanitizeFilesystemError(namespace, method, kind, relativePath, err)`.
3. On any subprocess failure, return
   `sanitizeExecutionError(namespace, method, kind, name, result, err)`.
4. On any HTTP failure, return a fixed message that omits dial/DNS/IP details.
5. Add a test that asserts the error contains no absolute path, `/tmp/`, IP
   address or `resolved` path (see `source/tools/*_test.go`).

## Session Loop

[engine.Session](../source/engine/Session.go) drives a single agent's
request/response loop:

1. `SendChatRequest` appends the user message and calls
   `inferChatCompletions` (HTTP POSTs to `/v1/chat/completions`).
   Transient failures (transport errors, `408`, `409`, `425`, `429`, `500`,
   `502`, `503`, `504` and `529`) are retried up to four times with exponential
   backoff and jitter. A `Retry-After` header is honored up to 60 seconds, and a
   full call is bounded by a 10 minute timeout. Authentication, payment,
   permission, schema and model errors are not retried.
2. `ReceiveChatResponse` appends the assistant message; if it contains
   `tool_calls`, each is executed via `CallTool` and then the loop recurses into
   `inferChatCompletions`.
   - `agents.Await` blocks until the hired agent finishes.
   - `humans.Ask` and `humans.Choose` block until a human answers. The answer is
     submitted through `/api/session/calltool` as a `humans.Answer` and appends
     the correct message with `Role="tool"`.
3. `Recovery` snapshots the session and agents to
   `<playground>/.exocomp/` so a planner can be restored across restarts.
