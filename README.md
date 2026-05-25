
# Exocomp

Self-hosted multi-agent environment tailored to Golang.

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

Exocomp uses multiple [Agent Roles](./agents/Roles.go):

*Software Development*

| Role                                   | Lifecycle | Default Model     | Description                                 |
|:---------------------------------------|:---------:|:-----------------:|:--------------------------------------------|
| [planner](./agents/planner.yaml)       | long      | `gemma4:31b`      | writes with humans and plans project phases |
| [architect](./agents/architect.yaml)   | short     | `qwen3-coder:30b` | defines software specifications             |
| archivar                               | short     | `qwen3-coder:30b` | reads git repositories and wikis            |
| [coder](./agents/coder.yaml)           | short     | `qwen3-coder:30b` | implements specifications into code         |
| researcher                             | short     | `qwen3-coder:30b` | reads websites and API documentation        |
| [summarizer](./agents/summarizer.yaml) | short     | `qwen3-coder:30b` | reads long texts and summarizes them        |
| [tester](./agents/tester.yaml)         | short     | `qwen3-coder:30b` | implements unit tests, writes bug reports   |

*Pentesting*

| Role                                   | Lifecycle | Default Model             | Description                                                 |
|:---------------------------------------|:---------:|:-------------------------:|:------------------------------------------------------------|
| [exploiter](./agents/exploiter.yaml)   | short     | `qwen3-coder-heretic:30b` | implements exploits in CGo                                  |
| reverser                               | short     | `qwen3-coder-heretic:30b` | translates binaries or code into Go/CGo code                |
| threathunter                           | short     | `qwen3-coder-heretic:30b` | researches weaknesses and vulnerabilities in infrastructure |
| [webscanner](./agents/webscanner.yaml) | short     | `qwen3-coder-heretic:30b` | discovers vulnerabilities in web applications               |

Exocomp uses Tools to interact with the sandbox. Check the implementations to
see which tools are allowed for which Agent role.

Each of those Agent roles is specialized on using `golang` as their programming
language because `go test` allows to use integrated unit tests to document
issues with generated code very easily in a standardized manner.


### Tools

The Tools are work in progress at the moment, as they are an ongoing effort
to beat the context length and agent memory limitations of locally run models.

Check the unit tests on whether the Tools can be relied on or not.

| Tool                                    | Unit Tests?                         | Description                                         | Agent Roles                                             |
|:----------------------------------------|:-----------------------------------:|:----------------------------------------------------|:-------------------------------------------------------:|
| [Agents](./tools/Agents.go)             | [Yes](./tools/Agents_test.go) [1]   | Manages the lifecycle of contractor sub-agents.     | `planner`                                               |
| [Bugs](./tools/Bugs.go)                 | [Yes](./tools/Bugs_test.go)         | Manages documentation of discovered bugs.           | `tester`                                                |
| [Changelog](./tools/Changelog.go)       | [Yes](./tools/Changelog_test.go)    | Manages documentation of development changelog.     | `coder`                                                 |
| Containers                              |                                     | Manages virtual containers.                         | `redteamer`, `blueteamer`                               |
| Exploits                                |                                     | Manages PoCs for CVEs from local dataset.           | `pentester`, `reverser`                                 |
| [Files](./tools/Files.go)               | [Yes](./tools/Files_test.go)        | Interacts with files and folders.                   | `planner`, `architect`, `coder`, `summarizer`, `tester` |
| Findings                                |                                     | Reports findings of vulnerabilities and weaknesses. | `pentester`                                             |
| [Programs](./tools/Programs.go)         | [Yes](./tools/Programs_test.go)     | Interacts with installed programs.                  | `coder`, `tester`                                       |
| [Requirements](./tools/Requirements.go) | [Yes](./tools/Requirements_test.go) | Manages specifications of implementations.          | `architect`, `coder`, `tester`                          |
| [Skills](./tools/Skills.go)             |                                     | Loads and Unloads Agent Skills. [2]                 | `planner`, `architect`, `coder`, `tester`               |
| Kiwix                                   |                                     | Researches knowledge from offline web archives.     | `researcher`                                            |
| Vulnerabilities                         |                                     | Manages vulnerabilities from local dataset.         | `pentester`, `threathunter`                             |
| Websites                                |                                     | Researches knowledge from the web.                  | `pentester`, `researcher`                               |

- [1] Requires `llama.cpp` with `qwen3-coder:30b` and `Q8_0` quantization and 48GB VRAM GPU with `vulkan` support.
- [2] Implements `SKILL.md` support, in compliance with [agentskills.io/specification](https://agentskills.io/specification).


### Building

The `exocomp` project comes in several variants. All programs support `CGO_ENABLED=0`,
so they can be used without any dynamically linked dependencies.

- [agimus](./source/cmds/agimus/main.go) which is used for testing `assistant` sandboxes.
- [exocomp](./source/cmds/exocomp-web/main.go) which supports all UIs.
- [exocomp-agent](./source/cmds/exocomp-agent/main.go) which supports only the `agent` and `terminal` UI.
- [exocomp-web](./source/cmds/exocomp-web/main.go) which supports only the `agent` and `web` UI.
- [exocomp-installer](./installer/cmds/exocomp-installer/main.go) that bundles all exocomp builds and required agent programs.


```bash
# Build exocomp and exocomp-installer
cd path/to/exocomp/toolchain;
go run build.go;

# Show exocomp usage instructions
cd /path/to/exocomp/build;
./linux/exocomp;
```


### Testing

Testing requires a [llama.cpp](https://github.com/ggml-org/llama.cpp)
`llama-server` instance running with a `qwen3-coder:30b` model and `Q8_0`
quantization. Take a look at the [TESTING.md](./docs/TESTING.md) for more
details.


### Usage

Exocomp's sandboxes are based on the `current working directory`, meaning that
the folder in which the program is executed is the sandbox that the running agent
can't escape from.

The recommended default usage is to use the `Web UI` so that you can observe
other agents working for the agent you're talking with.

```bash
cd /path/to/project-root;

# Agent type planner is defaulted
exocomp webview planner;
```

Take a look at the [USAGE.md](./docs/USAGE.md) for more details.


### Supported Inference Servers

It's also possible to use exocomp with an external inference server
that supports the OpenAI compatible endpoints. Take a look at the
[SERVERS.md](./docs/SERVERS.md) for more details.


### License

Dual Licensed. AGPL3 for private usage. EULA for commercial usage available.
For a commercial license, contact [Cookie Engineer](https://cookie.engineer).

As you might have imagined, this is a not-so-serious project at this stage.
Maybe it works, maybe it doesn't. Only the future will be able to tell whether
the LLM hype of agentic coding/debugging environments was justified.

