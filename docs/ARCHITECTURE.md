# Architecture

Exocomp is a self-hosted multi-agent environment for Go packages. Agents run as
sandboxed OS subprocesses and communicate through explicit tool calls instead of
freeform documents. A running agent works on a `symbol`: a `"path/to/package"`
plus a function, method, struct, interface or named type. Sandboxes map to
packages, so parallel agents can work without conflicting.

## Modules

| Path | Module | Purpose |
|:-----|:-------|:--------|
| `source/` | `exocomp` | Runtime, tools, engine and user interfaces |
| `app/` | `exocomp-app` | WASM web client |
| `installer/` | `installer` | Embedded installer |
| `toolchain/` | `exocomp-toolchain` | Build, test and lint programs |

The `toolchain/` module keeps each command in its own `main` package:
`build/`, `build-webapp/`, `clean/`, `install/`, `test/` and `lint/`. Shared
helpers live under `internal/utils/<domain>/`; web app build actions live under
`internal/utils/webapp/`.

## Source Layout

```text
source/
  actions/       CLI actions (config dump, usage text)
  adapters/      Provider-specific request/response transforms
  agents/        Role templates (*.yaml), prompts and message rendering
  parsers/
    html/        HTML to markdown/text conversion
    markdown/    Shared markdown parser for the web and TTY frontends
    yaml/        YAML parser used by the config and role loaders
  providers/     Provider presets (URL, models, pricing)
  tools/         Toolset wiring and one package per tool
    agents/        agents.*        contractor lifecycle and profile switch
    assets/        assets.*        recon asset inventory
    bugs/          bugs.*          discovered bug reports
    changelog/     changelog.*     development changelog
    containers/    containers.*    podman containers
    files/         files.*         sandboxed file and folder access
    findings/      findings.*      findings, evidence and reports
    humans/        humans.*        questions and capabilities
    programs/      programs.*      installed program execution
    requirements/  requirements.*  implementation specifications
    shell/         shell.*         interactive PTY sessions
    skills/        skills.*        Agent Skills
    websites/      websites.*      Reader Mode fetching
    testdata/      Fake binaries used by tests
    testutil/      Shared test binaries and sandbox helpers
  types/         Shared types (Agent, Config, Asset, Bug, Finding, ...)
  schemas/       Wire types (Message, ToolCall, ChatRequest/Response, ...)
  utils/         Shared helpers (api, ast, bytes, capabilities, chat, cli,
                 findings, fmt, fs, http, landlock, podman, rand, sandbox,
                 seatbelt, shell, slices, strings, time, tty)
  ui/
    jsonl/       Child-process frontend with JSONL over stdin/stdout
    tty/         Interactive terminal frontend
    web/         Web frontend and REST API
  cmds/
    agimus/          Easter egg TTY entry point
    exocomp/         UI dispatcher (config | terminal)
    exocomp-agent/   Child-process entry point (JSONL only)
    exocomp-web/     Web server entry point (config | web)
    exocomp-sandbox/ Landlock / Seatbelt helper spawned by the sandbox
```

`source/tools/Toolset.go` plus `Toolset_test.go` are the only files in the
`tools` package. `Toolset` imports each tool package and returns the tools
allowed for a role. `engine` depends on `tools` and the individual tool
packages it type-asserts.

`parsers/markdown` owns syntax and the `Document`/`Element` AST. The web app
renders that AST to safe HTML in `app/utils/html`; the TTY renders it to styled
transcript lines in `source/ui/tty/components`. Both frontends share the
parser.

## Process Hierarchy and Sandboxing

Communication between the human, the planner and contractor agents is a strict
subprocess tree:

```text
planner        (exocomp-web web | exocomp terminal, playground == sandbox == project root)
|-> architect  (exocomp-agent, sandbox == project root)
|-> coder      (exocomp-agent, sandbox == project/utils)
|-> tester     (exocomp-agent, sandbox == project/utils)
|-> coder      (exocomp-agent, sandbox == project/cmds/fibonacci)
```

- A child's `--playground` is its parent's `--sandbox`.
- A child's `--sandbox` is the package subfolder it may touch.
- Every agent is a real OS process with a distinct working directory, so the
  process tree mirrors the agent sandboxes.

### Symlink Policy

Path-bound tools reject symbolic links that could resolve outside the sandbox.
`utils_sandbox.ResolvePath` and `utils_sandbox.SanitizePath` validate every
existing path component with `os.Lstat` and refuse symlinked components.
Directory listings, searches and folder copies skip symlink entries.
`programs.Execute` applies the same validation to path-like arguments.

### Process Confinement (Landlock / Seatbelt)

`programs.Execute`, `shell.*` and `skills.Execute` run their subprocesses
through the `exocomp-sandbox` helper. The helper applies the platform sandbox
before executing the target:

- Linux: a [landlock](https://docs.kernel.org/userspace-api/landlock.html)
  ruleset is applied before `execve`.
- macOS: an SBPL profile is generated and applied through
  `/usr/bin/sandbox-exec -p`.

Both backends share a fixed filesystem and process policy:

- The sandbox folder is the only writable location.
- A read-only allowlist covers system binaries, libraries, configuration and
  the resolved target program directory.
- Basic device nodes (`/dev/null`, `/dev/zero`, `/dev/urandom`, ...) are
  readable and writable.
- Temporary `humans.AskPermission` grants can add read-only or read-write paths.
- There is no `off` switch and no fallback. If the platform sandbox or the
  helper is unavailable, execution fails with a fixed error.

Network access is denied by default:

- Without a `network.Connect` capability, Linux runs in a private network
  namespace with only loopback raised. macOS emits no network rules at all.
- With a `network.Connect` capability, network access is enabled for the whole
  sandbox. Linux stays in the host network namespace. macOS allows outbound IP
  traffic (`*:*`), the mDNSResponder resolver socket and loopback bind/inbound
  for local proxies and test servers. Unrelated unix-socket connects
  (docker.sock, SSH agent, ...) stay denied.

This is all-or-nothing because neither landlock nor seatbelt can match an
individual host: landlock network rules are port-based, and seatbelt accepts
only `*` or `localhost` as the host token in an IP filter. The domain, IP and
CIDR scopes of `network.Connect` are enforced at the `websites.*` boundary, by
role prompts and by `assets.Scope`, not by the kernel.

On Linux, `PR_SET_NO_NEW_PRIVS` is set before restricting, which also blocks
setuid privilege escalation. The Linux backend refuses to run on other operating
systems. `/proc` and `/sys` are readable for compatibility. The helper also runs
in private user and mount namespaces (plus a private network namespace when
network access is denied), masks `<sandbox>/.exocomp` with an empty read-only
tmpfs, drops all capabilities and runs the target as uid 0 inside the user
namespace, mapped to the invoking host user.

On macOS, the SBPL profile is deny-by-default. Paths are resolved through
`filepath.EvalSymlinks` before being emitted because the kernel matches resolved
vnodes. Both the original and resolved forms are included, which covers
`/tmp` -> `/private/tmp`, `/var` -> `/private/var` and `/etc` -> `/private/etc`.
The profile allows only the configured file, process, sysctl, signal, Mach
lookup and POSIX semaphore operations. Using `sandbox-exec` instead of the cgo
`sandbox_init` API keeps the project buildable with `CGO_ENABLED=0`.

HOME, temp, runtime and Go cache directories are redirected into a private
runtime directory under the system temp directory, keyed by the sandbox path.
The runtime directory is outside the sandbox and `.exocomp`, and is granted to
the helper as a writable path, so tools such as `go`, `gofmt` and `gcc` can
write without touching internal state. Credential and session sockets inherited
from the agent environment (`SSH_AUTH_SOCK`, `GPG_AGENT_INFO`,
`DBUS_SESSION_BUS_ADDRESS`, `KRB5CCNAME`, `DISPLAY`, ...) are removed from the
child environment.

`.exocomp` is reserved for tool-managed state and ignored by the Go toolchain.
All `files.*` methods reject paths inside `.exocomp`; use `findings.*`,
`assets.*`, `findings.BindEvidence` and `shell.*` to access that state. Executed
programs cannot access it either: the Linux helper masks it, the macOS profile
adds an explicit deny rule, and containers get a read-only tmpfs over
`/workspace/.exocomp`.

The helper is discovered next to the running `exocomp` binary, or via
`EXOCOMP_SANDBOX` for development. The installer places it in the same `bin`
folder as `exocomp` and `exocomp-agent`.

`exocomp-agent` is the only JSONL child-process binary. `agents.Hire` and the
session summarizer look for it next to the running binary first, then on `PATH`;
`EXOCOMP_AGENT` overrides both for development. `exocomp` and `exocomp-web`
expose only their own frontend and have no `agent` mode.

### Temporary Capabilities

`humans.AskPermission` lets an agent request temporary capabilities. A request
contains a reason and a list of `types.Capability` entries:

- `files.Read` with an absolute path or folder.
- `files.Write` with an absolute path or folder (implies read access).
- `network.Connect` with an exact domain, exact IP or CIDR.

The human selects capabilities individually. Granted capabilities live in
`engine.Session.Capabilities` with `json:"-"`, so they are never written to the
config, recovery backup or any other on-disk state. They are cleared when the
session ends or the agent role switches, and are not inherited by hired child
agents.

`files.Read` and `files.Write` are enforced by the OS sandbox: paths are added
to the landlock/seatbelt policy for subsequent `files.*`, `programs.Execute`,
`shell.*` and `skills.Execute` calls. The policy of a running process is fixed
when it is spawned. A long-lived `shell.*` session keeps the policy it was
started with; close and reopen it after a capability is granted.

`network.Connect` is enforced in two places:

- `websites.*` matches the exact domain, IP or CIDR on every request and
  redirect.
- The OS sandbox switches network access on or off for `programs.Execute`,
  `shell.*` and `skills.Execute`. A single grant enables outbound network access
  for every sandboxed program in the session. Raw programs such as `curl`,
  `nmap`, `nc` or `ssh` are therefore not domain-filtered; scope is enforced by
  role prompts and `assets.Scope`.

Known limitations: on Linux, `/proc` and `/sys` are readable, and metadata
changes, ptrace, signals and already-open file descriptors are not mediated. In
grant mode the sandbox returns to the host network namespace, so it can reach
host loopback services and use UDP and ICMP. On macOS, the seatbelt profile
denies unrelated unix sockets, but the platform still does not mediate every
operation and pre-opened file descriptors remain usable. Programs that expect
host configuration or caches outside the sandbox may need sandbox copies.

### Container Execution (Podman)

`containers.*` shells out to an externally installed `podman` binary. It does
not use `exocomp-sandbox`: podman needs its rootless socket, runtime
directories, user namespace helpers and network access, which the process
sandbox does not grant. The Go code stays buildable with `CGO_ENABLED=0`
because there is no podman Go binding and no cgo dependency. On Windows the tool
works with a podman machine, even though the host process sandbox is unavailable
there.

Container isolation is enforced by the tool:

- Every container gets the label `io.exocomp.sandbox=<sha256(sandbox)[:16]>`.
  `List` filters on this label and every lifecycle or exec operation verifies it
  through `podman inspect`.
- The sandbox is mounted at `/workspace` and is the default workdir. A read-only
  tmpfs is mounted over `/workspace/.exocomp`. No arbitrary volumes, port
  publishing, capabilities, privileged mode or host namespaces are accepted
  from the model.
- Images must appear in the role's `allowed-images` list. An empty list denies
  all images.
- The default network mode is `none`; `bridge` must be requested explicitly and
  `host` is never allowed.
- Podman output that reaches the model goes through `utils/podman`'s
  `sanitizeError`, so host paths and raw runtime errors are not leaked.

## JSONL Protocol

A hired agent runs `exocomp-agent` and talks to its parent through line-delimited
JSON on `stdout`. `stderr` is free-form logging. Every line is formatted as
`<prefix>:<json>\n`:

| Prefix | Payload | Producer | Consumer |
|:-------|:--------|:---------|:---------|
| `schemas.Message` | [schemas.Message](../source/schemas/Message.go) | child renderer | parent reader |
| `types.ContextUsage` | [types.ContextUsage](../source/types/ContextUsage.go) | child heartbeat | parent reader |

The child emits `schemas.Message:` lines for every system, assistant and tool
message it produces. The parent appends them to the child's
`types.Agent.Messages` so the planner can inspect and summarize the work. The
`types.ContextUsage:` line doubles as a liveness heartbeat.

`stdin` is the reverse channel. It accepts one `schemas.Message` JSON object per
line with `Role: "user"`. The TTY uses this to message a hired agent directly.
Messages are queued behind the child's current request, so this is manual
guidance rather than a preemptive interrupt. Use `agents.Fire` to stop a child.

The TTY renders two sidebars: agents on the left and session state
(`CHANGELOG`/`OPEN BUGS`) on the right. `Ctrl+Up/Down` cycles the selected agent
and `Alt+Up/Down` scrolls the transcript. The composer sends to the selected
agent; the planner remains `Session.Agent` and is never replaced for viewing.

Both UI frontends read agent state through `engine.Session` snapshot methods:
`GetAgents`, `GetAgent`, `GetAgentInfo` and `GetAgentMessages`. The web routes
use the same methods as the TTY, so agent lists, configs and conversation
histories behave identically across frontends.

## Agent Lifecycle

Agent state is tracked on [types.Agent](../source/types/Agent.go) through the
`Status` field:

| Status | Meaning |
|:-------|:--------|
| `working` | Hired, subprocess running |
| `finished` | Exited `0` with an `agents.Quit` work report |
| `failed` | Exited without a work report (non-zero, killed) |
| `fired` | Killed via `agents.Fire` |

Lifecycle is owned by the [Agents tool](../source/tools/agents/Agents.go).

`agents.Hire(role, prompt, sandbox)`:

1. Generates a name.
2. Resolves the sandbox.
3. Spawns `exocomp-agent --name ... --prompt ... --sandbox ...`.
4. Records the `*types.Agent` as `working` and starts one reader goroutine.

The reader goroutine owns the child's stdout. It reads lines without a
`bufio.Scanner` token cap, appends `schemas.Message`s, updates
`ContextUsage`/liveness, and on `EOF` calls `cmd.Wait()` after all reads,
finalizes `Status` to `finished` or `failed`, closes the process `done` channel
and removes the process from the running map.

`agents.Await(name)` blocks on `done` until the child finishes and returns the
`agents.Quit: Work Report` or a terminal error. It deliberately does not retry:
polling would fill the planner's context with identical "still working"
messages.

`agents.Fire(name)` cancels the child context, waits for `done` and marks the
agent `fired`.

`agents.Quit(message)` is called by the child. The report is flushed
synchronously to stdout before the process exits, so the parent reliably
observes it.

`agents.Switch(role)` switches the main TTY session to another preconfigured
profile. It is only wired in the TTY frontend; hired subprocesses expose the tool
but receive an error. Before switching, the previous conversation is summarized
by a summarizer subprocess. The session starts fresh with the target profile's
system prompt and that summary as its first assistant message. Tools, programs
and adapters are rebuilt from the target profile so switching cannot retain
capabilities from the previous role.

A liveness watchdog cancels a child that stops producing output for longer than
`IdleTimeout` (default `5m`). A hard `Timeout` (default `30m`) bounds the total
lifetime. Both are fields on `Agents` and are set in `NewAgents`. `types.Agent`
also carries `StartedAt` and `FinishedAt` timestamps for the schedule view.

## Tool Lifecycle

Each agent role declares `allowed-tools`, `allowed-programs` and
`allowed-images` in its YAML `agents/*.yaml` file. At startup the frontend builds
a `tools.Toolset` containing only the allowed tools and their JSON schemas, and
registers them on the `Session` keyed by namespace (the part before the first
`.`).

- Tool schema names are `namespace.Method` (for example `files.Read`,
  `requirements.DefineFunc`, `agents.Hire`).
- `Session.CallTool` resolves the namespace, invokes the tool and appends the
  result (or `Error: ...`) as a `tool` message carrying the matching
  `tool_call_id`.
- Sandbox-bound tools (files, programs, requirements, bugs, changelog, findings,
  assets) resolve paths through `utils_sandbox.ResolvePath` and
  `utils_sandbox.SanitizePath`.
- `containers.*` mounts only the agent sandbox and scopes every container to the
  sandbox label. `podman` must not be added to `allowed-programs`, otherwise
  `programs.Execute` could bypass these checks.
- `findings.*` stores one JSON file per finding plus an index under
  `<playground>/.exocomp/findings/`. Evidence artifacts live under
  `<playground>/.exocomp/evidence/<finding-id>/` with SHA-256 hashes. Inline
  evidence is redacted and capped; artifact contents are never returned through
  the tool. `findings.Import` parses nuclei JSONL artifacts and creates candidate
  findings. `findings.Export` copies findings, evidence, reports or stored scan
  artifacts into `./work`.
- `assets.*` stores a deduplicated inventory under
  `<playground>/.exocomp/assets/index.json`. `assets.Import` parses httpx,
  naabu, subfinder, dnsx, tlsx, katana, gobuster and urlfinder artifacts.
  `assets.Export` writes host lists, URL lists or JSON into `./work`.
- There is no `scans.*` tool. Agents run the preshipped Go scanner binaries
  through `programs.Execute` or `shell.*` and write structured artifacts under
  `./work/scans/`. The importers consume those artifacts, so the model never
  parses raw scanner output. Tool-managed state stays private in `.exocomp` and
  is accessed only through import/export methods.

Each tool's LLM-facing contract lives in the embedded JSON schema next to the
tool, for example `tools/files/Files.json`.

Role separation is enforced by capability, not by the prompt alone. An architect
has `files.Read` but not `files.Write` and records output with
`requirements.Define*`. A coder gets `files.Write`, `changelog.*` and
`programs.Execute`. A tester writes unit tests and `bugs.*` reports. `planner`,
`architect`, `summarizer` and `webresearcher` have no allowed programs and no
`programs.*` tools, so they cannot execute code.

## Tool Error Message Sandboxing

Tool results and error messages are part of the model's context window. Raw Go
errors can leak host information. For example, `os.Stat`'s `*PathError` embeds
the absolute path:

```text
files.Stat: stat /home/alice/projects/secret/...: no such file or directory
```

`utils_sandbox.ResolvePath` and `utils_sandbox.SanitizePath` resolve every path
to an absolute location internally, so passing `err.Error()` back to the model
leaks the host layout. The same applies to `os.Executable()` paths in
`agents.Hire`, the `/tmp/...` temp dir in `agents.Inquire`, and `net/http` dial
errors in `websites.*`.

### Rules

1. Never return a raw `err.Error()` from `os`, `exec`, `net/http` or `json`
   through a tool's `Call(...)`. Map it to a fixed, path-safe message.
2. Echo the user-supplied or sanitized relative path (`"./x"`), never the
   resolved absolute path.
3. Keep messages short, deterministic and prefixed `namespace.Method:`. Prefer
   `File "./x" does not exist.` over an ENOENT string.

### Helpers

| File | Purpose |
|:-----|:--------|
| `source/tools/files/sanitizeFilesystemError.go` | Maps `os.*` filesystem errors to `does not exist.`, `Permission denied.` or `Cannot access ...`, using the relative path. |
| `source/utils/sandbox/SanitizeExecutionError.go` | Maps subprocess errors from `programs.Execute` and `skills.Execute` to `Permission denied.`, `Program doesn't exist.` or `exited with an error.` |

### Adding a Tool

1. Resolve and sanitize every path through `utils_sandbox.ResolvePath` or
   `utils_sandbox.SanitizePath`.
2. On `os.*` failure, return
   `sanitizeFilesystemError(namespace, method, kind, relativePath, err)`.
3. On subprocess failure, return
   `sanitizeExecutionError(namespace, method, kind, name, result, err)`.
4. On HTTP failure, return a fixed message that omits dial, DNS and IP details.
5. Add a test that asserts the error contains no absolute path, `/tmp/`, IP
   address or resolved path. See `source/tools/*/*_test.go`.

## Session Loop

[engine.Session](../source/engine/Session.go) drives one agent's
request/response loop:

1. `SendChatRequest` appends the user message and calls
   `inferChatCompletions`, which posts to `/v1/chat/completions`. Transient
   failures (transport errors, `408`, `409`, `425`, `429`, `500`, `502`, `503`,
   `504` and `529`) are retried up to four times with exponential backoff and
   jitter. A `Retry-After` header is honored up to 60 seconds, and the complete
   call is bounded by a 10 minute timeout. Authentication, payment, permission,
   schema and model errors are not retried.
2. `ReceiveChatResponse` appends the assistant message. If it contains
   `tool_calls`, each is executed through `CallTool` and the loop recurses into
   `inferChatCompletions`.
   - `agents.Await` blocks until the hired agent finishes.
   - `humans.Ask` and `humans.Choose` block until a human answers. The answer is
     submitted through `/api/session/calltool` as a `humans.Answer` and appends
     the correct message with `Role="tool"`.
3. `Recovery` snapshots the session and agents to `<playground>/.exocomp/`, so a
   planner can be restored across restarts.

## Utility Ownership

| Package | Owns |
|:--------|:-----|
| `types` | Wire and data types such as `Capability`, `Config` and `Question`. Must not import `utils/capabilities`. |
| `utils/capabilities` | Capability validation, normalization, canonical paths, path and network matching, and the capability `Set`. |
| `utils/chat` | Tool-call `Arguments` accessors, tool-call discovery, retry and token-estimation helpers. |
| `utils/findings` | Finding and asset IDs, text limits, redaction, sorting, query matching, evidence storage, index persistence and report rendering. |
| `utils/http` | MIME content types and the HTTP handler helper. |
| `utils/sandbox` | Sandbox path validation and process execution. |
| `utils/shell` | Generic shell parsing, formatting, terminal filtering, process groups, signals and shell environments. |
| `utils/slices` | Order-preserving unique string slices. |
| `tools/<namespace>` | Tool schemas, argument adaptation, result formatting and tool-specific errors. |
| `tools/testutil` | Shared test binaries and sandbox test helpers. |

Dependency rules:

- `types` must not import `utils/capabilities`.
- `utils` must not import `tools`.
- Tool packages must not import each other. Their project imports are limited
  to `types`, `schemas`, `agents`, `parsers/*` and `utils/*`.
- `engine` may import `utils/capabilities` and individual tool packages.
- `utils/sandbox` owns sandbox paths and process execution and must not import
  `utils/findings`, `utils/shell` or `tools`.
- `utils/findings` and `utils/shell` may depend on `utils/sandbox`.
- Shared helpers have exactly one home. Do not create `utils/<name>` packages
  that mirror a top-level package.

The dependency rules are implemented in [`toolchain/lint/deps.go`](../toolchain/lint/deps.go)
and are checked by `make lint-toolchain`.
