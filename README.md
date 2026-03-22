
# Exocomp

Attention, profit-seekers and visionaries!

Are organics costing you time, wages, and - worst of all - benefits?

Then upgrade your operation today with the all new Exocomp (tm) adaptive repair
unit, the smartest investment this side of the Alpha Quadrant!

<img width="256" height="256" src="https://raw.githubusercontent.com/cookiengineer/exocomp/master/docs/exocomp.png"/>

Why hire when you can own? The Exocomp (tm) isn't just a tool... it's so much more!

- Rapid autonomous diagnostics
- On-the-fly gadget fabrication
- Precision repairs of broken code with unit tests
- Supervised task queue functionality
- Self-replicating in hazardous environments
- Tireless performance, no sleep cycles, no unions, no complaints!

From starship maintenance to high-risk industrial operations, the Exocomp (tm)
delivers maximum output with minimal oversight. Think of it as an employee,
except it doesn't cheat you out of your profits on the Dabo table.

Ethical subroutines sold separately.


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

