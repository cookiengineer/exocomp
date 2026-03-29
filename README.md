
# Exocomp

Attention, profit-seekers and visionaries!

Are organics costing you time, wages, and - worst of all - benefits?

Then upgrade your operation today with the all new Exocomp adaptive repair
unit, the smartest investment this side of the Alpha Quadrant!

<img width="256" height="256" src="https://raw.githubusercontent.com/cookiengineer/exocomp/master/docs/exocomp.png"/>

## Why hire when you can own?

The Exocomp isn't just a tool... it's so much more!

- Supervised task queue functionality
- Multiple agent roles as `manager`, `coder`, `tester`
- Precision repairs of broken code with unit tests
- Cross-agent communication via `NOTES.md` and `TODO.md`
- Self-replicating in malicious environments
- Rapid autonomous diagnostics
- Tireless performance with no sleep cycles, no unions, no complaints!

From starship maintenance to high-risk industrial operations, the Exocomp
delivers maximum output with minimal oversight. Think of it as an employee,
except it doesn't cheat you out of your profits on the Dabo table.

Ethical subroutines sold separately.

## Architecture

Exocomp uses multiple agent roles:

- `manager` role writes features into backlog
- `coder` implements features, reads `feature` backlog and `bug` reports
- `tester` implements unit tests, writes `bug` reports

Exocomp uses Tools to interact with the sandbox:

- [ ] [Agents](./tools/Agents.go) tool to manage the lifecycle of sub-agents.
- [ ] [Bugs](./tools/Bugs.go) tool to manage documentation of bugs.
- [x] [Changelog](./tools/Changelog.go) tool to manage development notes.
- [ ] [Features](./tools/Features.go) tool to manage documentation of features.
- [x] [Files](./tools/Files.go) tool to read/write/list files and folders.
- [x] [Programs](./tools/Programs.go) tool to execute programs.
- [ ] [Web](./tools/Web.go) tool to research things on the web.


## Requirements and Usage

The `exocomp` program is a standalone binary, once compiled with the `go` toolchain.
However, the models currently aren't embedded and are called via an external
(locally hostable) ollama server.

```bash
sudo pacman -S go ollama;

# Start ollama server
ollama serve;

# Install qwen3 coder model
ollama pull qwen3-coder:30b

# Run exocomp with ollama
cd /path/to/exocomp;
go run exocomp --url="http://localhost:11434/api" --model="qwen3-coder:30b";
```


## License

Dual Licensed. AGPL3 for Open Source usage. EULA for Commercial usage available.
For a commercial license, contact [Cookie Engineer](https://cookie.engineer).

As you might have imagined, this is a not-so-serious project at this stage.
Maybe it works, maybe it doesn't. Only the future will be able to tell whether
the LLM hype of agentic coding/debugging environments was justified.

