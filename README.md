
# Exocomp

Self-hosted short-living multi-agent agentic development environment tailored to Golang.

<p align="center">
    <img width="256" height="256" src="https://raw.githubusercontent.com/cookiengineer/exocomp/master/docs/exocomp.png"/>
</p>

### Features

Attention, profit-seekers and visionaries!

Are organics costing you time, wages, and - worst of all - benefits?
Then upgrade your operation today with the all new Exocomp adaptive repair
unit, the smartest investment this side of the Alpha Quadrant!

Why hire when you can own? The Exocomp isn't just a tool... it's so much more!

- Supervised task queue functionality
- Multiple agent roles as `architect`, `coder`, `planner`, `summarizer`, `tester`
- Precision repairs of broken code with unit tests
- Cross-agent communication via `changelog`, `bugs` and `requirements` tools
- Self-replicating short lived agents in malicious environments
- Rapid autonomous diagnostics
- Tireless performance with no sleep cycles, no unions, no complaints!

From starship maintenance to high-risk industrial operations, the Exocomp
delivers maximum output with minimal oversight. Think of it as an employee,
except it doesn't cheat you out of your profits on the Dabo table.

Ethical subroutines sold separately.


### Agents

Exocomp uses multiple Agent [Types](./agents/Types.go):

| Role                                 | Lifecycle | Default Model     | Description                                                   |
|:-------------------------------------|:---------:|:-----------------:|:--------------------------------------------------------------|
| [Planner](./agents/Planner.go)       | long      | `gemma4:31b`      | writes with humans, contracts agents and plans project phases |
| [Architect](./agents/Architect.go)   | short     | `gemma4:31b`      | discusses with humans and writes specifications               |
| [Coder](./agents/Coder.go)           | short     | `qwen3-coder:30b` | implements features, reads `specifications` and `bugs`        |
| [Tester](./agents/Tester.go)         | short     | `qwen3-coder:30b` | implements unit tests, writes reports into `bugs`             |
| [Summarizer](./agents/Summarizer.go) | short     | `qwen3-coder:30b` | reads long texts and summarizes them                          |
| Researcher                           | short     | `qwen3-coder:30b` | reads websites, API documentation, and reports to `Architect` |

Exocomp uses Tools to interact with the sandbox. Check the implementations to
see which tools are allowed for which Agent role.

Each of those Agent roles is specialized on using `golang` as their programming
language because `go test` allows to use integrated unit tests to document issues
with generated code very easily in a standardized manner.


### Tools

The Tools are work in progress at the moment, as they are an ongoing effort
to beat the context length and agent memory limitations of locally run models.

Check the unit tests on whether the Tools can be relied on or not.

| Tool                                    | Unit Tests?                                              | Description                                     | Agent User Roles                                        |
|:----------------------------------------|:--------------------------------------------------------:|:------------------------------------------------|:-------------------------------------------------------:|
| [Agents](./tools/Agents.go)             | [Yes](./tools/Agents_test.go) (requires `llama.cpp` [1]) | Manages the lifecycle of contractor sub-agents. | `planner`                                               |
| [Bugs](./tools/Bugs.go)                 | [Yes](./tools/Bugs_test.go)                              | Manages documentation of discovered bugs.       | `tester`                                                |
| [Changelog](./tools/Changelog.go)       | [Yes](./tools/Changelog_test.go)                         | Manages documentation of development changelog. | `coder`                                                 |
| Containers                              |                                                          | Manages virtual containers.                     | `redteamer`, `blueteamer`                               |
| [Files](./tools/Files.go)               | [Yes](./tools/Files_test.go)                             | Interacts with files and folders.               | `planner`, `architect`, `coder`, `summarizer`, `tester` |
| [Programs](./tools/Programs.go)         | [Yes](./tools/Programs_test.go)                          | Interacts with installed programs.              | `coder`, `tester`                                       |
| [Requirements](./tools/Requirements.go) | [Yes](./tools/Requirements_test.go)                      | Manages specifications of implementations.      | `architect`, `coder`, `tester`                          |
| Forgejo                                 |                                                          | Researches knowledge from offline git servers.  | `researcher`                                            |
| Kiwix                                   |                                                          | Researches knowledge from offline web archives. | `researcher`                                            |
| Websites                                |                                                          | Researches knowledge from the web.              | `researcher`                                            |

[1] Install dependencies with the [install-deps.sh](./install-deps.sh) shell script. Requires at least 80GB of HDD space, 48GB of VRAM, and an iGPU or dGPU with `vulkan` support.

### Dependencies

The `exocomp` program is a standalone binary, once compiled with the `go`
toolchain and it comes with `llama.cpp` as a bundled inference server.

The third-party `llama.ccp` inference server and LLM models are installed
inside the [third_party](./third_party) folder, and are installed with the
[install-deps.sh](./install-deps.sh) shell script.

```bash
cd /path/to/exocomp;

# Look ma, no sudo!
bash install-deps.sh;
```

It is heavily recommended to use that workflow for the best development
experience using exocomp as an agentic development environment.

The unit tests tagged with the `agents` build tag also rely on this workflow
to be setup and working.


### External Inference Servers

However, it's also possible to use exocomp with an external inference
server that supports the OpenAI compatible endpoints. If you're using
`ollama` as a beginner, all models with the `tools` tag in the
[ollama library](https://ollama.com/library) should be compatible.

Exocomp uses the following OpenAI compatible API endpoints:

- `http://server:port/v1/chat/completions`
- `http://server:port/v1/models`

Take a look at the [INFERENCE_SERVERS.md](./docs/INFERENCE_SERVERS.md) for more
details on how to use external inference servers.


### Usage

Exocomp's sandboxes are based on the `current working directory`, meaning
that the folder the program is executed in is the sandbox that the running
agent can't escape from.

The recommended default usage is to use the `Web UI` so that you can observe
other agents working for the agent you're talking with.

```bash
cd /path/to/project-root;

# Agent type planner is defaulted
exocomp webview planner;
```

Additionally, there are several UI frontends implemented in Exocomp:

- `exocomp jsonl <agent-type>` uses line-separated JSON messages to communicate via `stdin` and `stdout`. Used for cross-contractor-agent communication.
- `exocomp tty <agent-type>` uses the [ui/tty/Client](./ui/tty/Client.go)
- `exocomp web <agent-type>` spawns the [ui/web/Server](./ui/web/Server.go) on port `3000` that serves the [Web UI](./ui/web/public/)
- `exocomp webview <agent-type>` spawns a [ui/web/Server](./ui/web/Server.go) on port `3000` and opens a [ui/webview/Client](./ui/webview/Client.go) window

### Multi Agent Usage

The defaulted `planner` agent is allowed to hire contracting sub-agents with the
[Agents](./tools/Agents.go) tool.

Multi agent communication works with a sub process hierarchy, where each process
works in their own sandbox with the `jsonl` frontend and their own agent type and
system/user prompts.

Cross-agent communication works with a strict process hierarchy, each contracted
sub-agent's `playground` is set to the parent process's `sandbox`.

**Example Process Hierarchy**

```
# Hierarchy   # Sandbox                         | Playground       | Task                          |
# ----------- # ------------------------------- | ---------------- | ----------------------------- |
| planner     # /path/to/project                | /path/to/project |                               |
|-> architect # /path/to/project                | /path/to/project | specifies utils package       |
|-> coder     # /path/to/project/utils          | /path/to/project | implements CalculateFibonacci |
|-> tester    # /path/to/project/utils          | /path/to/project | tests CalculateFibonacci      |
|-> coder     # /path/to/project/cmds/fibonacci | /path/to/project | implements main.go            |
```

**Example Process Parameters**

```bash
# Humans interact with planner agents
cd /path/to/project;
exocomp web planner;

#
# What the planner "spawns" behind the scenes as sub processes:
#

# cd /path/to/project;
# exocomp jsonl architect --prompt="Implement a utils package and specify the CalculateFibonacci method signature.";

# cd /path/to/project/utils;
# exocomp jsonl coder --prompt="Implement a public method called CalculateFibonacci(step int) int that calculates the fibonacci numbers.";

# cd /path/to/project/utils;
# exocomp jsonl tester --prompt="Implement the unit tests for the CalculateFibonacci method.";

# cd /path/to/project/cmds/fibonacci;
# exocomp jsonl coder --prompt="Implement a main.go that calculates fibonacci numbers. Use the first CLI parameter as the step or sequence argument.";
```


### License

Dual Licensed. AGPL3 for private usage. EULA for commercial usage available.
For a commercial license, contact [Cookie Engineer](https://cookie.engineer).

As you might have imagined, this is a not-so-serious project at this stage.
Maybe it works, maybe it doesn't. Only the future will be able to tell whether
the LLM hype of agentic coding/debugging environments was justified.

