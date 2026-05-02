
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

| Tool                                    | Unit Tests?                         | Description                                     | Agent User Roles                                        |
|:----------------------------------------|:-----------------------------------:|:------------------------------------------------|:-------------------------------------------------------:|
| [Agents](./tools/Agents.go)             |                                     | Manages the lifecycle of contractor sub-agents. | `manager`                                               |
| [Bugs](./tools/Bugs.go)                 | [Yes](./tools/Bugs_test.go)         | Manages documentation of discovered bugs.       | `tester`                                                |
| [Changelog](./tools/Changelog.go)       | [Yes](./tools/Changelog_test.go)    | Manages documentation of development changelog. | `coder`                                                 |
| Containers                              |                                     | Manages virtual containers.                     | `redteamer`, `blueteamer`                               |
| [Files](./tools/Files.go)               | [Yes](./tools/Files_test.go)        | Interacts with files and folders.               | `manager`, `architect`, `coder`, `summarizer`, `tester` |
| [Programs](./tools/Programs.go)         | [Yes](./tools/Programs_test.go)     | Interacts with installed programs.              | `coder`, `tester`                                       |
| [Requirements](./tools/Requirements.go) | [Yes](./tools/Requirements_test.go) | Manages specifications of implementations.      | `architect`, `coder`, `tester`                          |
| Forgejo                                 |                                     | Researches knowledge from offline git servers.  | `researcher`                                            |
| Kiwix                                   |                                     | Researches knowledge from offline web archives. | `researcher`                                            |
| Websites                                |                                     | Researches knowledge from the web.              | `researcher`                                            |


### Models

If you're using `ollama`, all models with the `tools` tag in the
[ollama library](https://ollama.com/library) should be compatible.

Exocomp uses the following OpenAI compatible endpoints:

- `http://server:port/v1/chat/completions`
- `http://server:port/v1/models`


### Requirements and Usage

The `exocomp` program is a standalone binary, once compiled with the `go`
toolchain. However, the models currently aren't embedded and are called via
an external (locally hostable) `ollama`, `vllm`, or `llama.cpp` server.

**Ollama Usage**:

```bash
# Start ollama server
ollama serve;

# Install qwen3 coder model
ollama pull qwen3-coder:30b

# Run exocomp with ollama
cd /path/to/exocomp;
go run ./cmds/exocomp/main.go web --debug;

# Custom CLI flags usage
# go run ./cmds/exocomp/main.go tty assistant --url="http://localhost:11434/v1" --model="qwen3-coder:30b";
```

**llama.cpp Usage**:

- Download [Qwen3 Coder 30B A3B Q8](https://huggingface.co/unsloth/Qwen3-Coder-30B-A3B-Instruct-GGUF) gguf model file.
- Download [llama.cpp](https://github.com/ggml-org/llama.cpp/releases) build.

If you're unsure about hardware support, use the `vulkan` build on Linux.

```bash
# Start llama.cpp server

# If you have a 48GB VRAM GPU:
# --ctx-size 262144
# --batch-size 2048
# --ubatch-size 512

./llama-server -m ./models/qwen3-coder-30b-a3b-instruct-q8_0.gguf \
	--gpu-layers all \
	--ctx-size 32768 \
	--batch-size 512 \
	--ubatch-size 128 \
	--cache-type-k q8_0 \
	--cache-type-v q8_0 \
	--flash-attn auto \
	--no-slots \
	--no-webui \
	--no-webui-mcp-proxy \
	--jinja \
    --port=11434;

# Run exocomp with llama.cpp
cd /path/to/exocomp;
go run ./cmds/exocomp/main.go web --debug;

# Custom CLI flags usage
# go run ./cmds/exocomp/main.go tty assistant --url="http://localhost:11434/v1" --model="qwen3-coder:30b";
```


### License

Dual Licensed. AGPL3 for private usage. EULA for commercial usage available.
For a commercial license, contact [Cookie Engineer](https://cookie.engineer).

As you might have imagined, this is a not-so-serious project at this stage.
Maybe it works, maybe it doesn't. Only the future will be able to tell whether
the LLM hype of agentic coding/debugging environments was justified.

