
# Exocomp

Attention, profit-seekers and visionaries!

Are organics costing you time, wages, and - worst of all - benefits?

Then upgrade your operation today with the all new Exocomp adaptive repair
unit, the smartest investment this side of the Alpha Quadrant!

<img width="256" height="256" src="https://raw.githubusercontent.com/cookiengineer/exocomp/master/docs/exocomp.png"/>


## Why hire when you can own?

The Exocomp isn't just a tool... it's so much more!

- Supervised task queue functionality
- Multiple agent roles as `architect`, `coder`, `tester`, `manager`
- Precision repairs of broken code with unit tests
- Cross-agent communication via `NOTES.md` and `TODO.md`
- Self-replicating in malicious environments
- Rapid autonomous diagnostics
- Tireless performance with no sleep cycles, no unions, no complaints!

From starship maintenance to high-risk industrial operations, the Exocomp
delivers maximum output with minimal oversight. Think of it as an employee,
except it doesn't cheat you out of your profits on the Dabo table.

Ethical subroutines sold separately.


## Agents

Exocomp uses multiple Agent roles:

| Role                                        | Default Model     | Description                                                   |
|:--------------------------------------------|:-----------------:|:--------------------------------------------------------------|
| [Manager](./agents/Agent.Manager.txt)       | `gemma4:31b`      | writes with humans, contracts agents and plans project phases |
| [Architect](./agents/Agent.Architect.txt)   | `gemma4:31b`      | discusses with humans and writes specifications               |
| [Coder](./agents/Agent.Coder.txt)           | `qwen3-coder:30b` | implements features, reads `specifications` and `bugs`        |
| [Tester](./agents/Agent.Tester.txt)         | `qwen3-coder:30b` | implements unit tests, writes reports into `bugs`             |
| [Summarizer](./agents/Agent.Summarizer.txt) | `qwen3-coder:30b` | reads long texts and summarizes them                          |
| Researcher                                  | TBD               | reads websites, API documentation, and reports to `Architect` |

Exocomp uses Tools to interact with the sandbox, so they're available differently
for each Agent role. The `manager` starts sub-agents that work on specified tasks,
so the `architect`, `coder` and `tester` roles are meant for short agent lifecycles.

Each of those Agent roles is specialized on using `golang` as their programming
language because `go test` allows to use integrated unit tests to document issues
with generated code very easily in a standardized manner.


## Models

If you're using `ollama`, all models with the `tools` tag in the [ollama library](https://ollama.com/library)
should be compatible. Use the API endpoint `http://ollama_instance:port/api/chat`.

If you're using `vllm`, all models with `tools` support should be compatible.
Use the API endpoint `http://vllm_instance:port/v1/chat/completions`.


## Tools

The Tools are work in progress at the moment, some might not be stable and/or not
work with LLMs at all.

Cross-Agent Tools:

- [x] [Agents](./tools/Agents.go) tool to manage the lifecycle of sub-agents.
- [x] [Bugs](./tools/Bugs.go) tool to manage documentation of bugs.
- [x] [Changelog](./tools/Changelog.go) tool to manage documentation of development notes.
- [ ] [Requirements](./tools/Requirements.go) tool to manage documentation of specifications.

Operating System Tools:

- [x] [Files](./tools/Files.go) tool to interact with files and folders.
- [x] [Programs](./tools/Programs.go) tool to interact with programs.

Knowledge-related Tools:

- [ ] [Gogs](./tools/Gogs.go) tool to search things in offline Gogs/Gitea instances.
- [ ] [Kiwix](./tools/Kiwix.go) tool to search things in offline knowledge bases.
- [ ] [Websites](./tools/Websites.go) tool to research things on the web.


## Requirements and Usage

The `exocomp` program is a standalone binary, once compiled with the `go` toolchain.
However, the models currently aren't embedded and are called via an external
(locally hostable) `ollama` or `vllm` server.

```bash
sudo pacman -S go ollama;

# Start ollama server
ollama serve;

# Install qwen3 coder model
ollama pull qwen3-coder:30b

# Run exocomp with ollama
cd /path/to/exocomp;
go run ./cmds/exocomp/main.go tty architect;

# custom CLI flags usage
# go run ./cmds/exocomp/main.go tty architect --url="http://localhost:11434/api" --model="qwen3-coder:30b";
```


## License

Dual Licensed. AGPL3 for private usage. EULA for commercial usage available.
For a commercial license, contact [Cookie Engineer](https://cookie.engineer).

As you might have imagined, this is a not-so-serious project at this stage.
Maybe it works, maybe it doesn't. Only the future will be able to tell whether
the LLM hype of agentic coding/debugging environments was justified.

