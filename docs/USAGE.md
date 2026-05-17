
# Usage

### Single-Agent Usage

There are several User Interfaces implemented in Exocomp:

- `exocomp agent` uses line-separated JSON messages to communicate via `stdin` and `stdout`. Used for cross-agent communication.
- `exocomp terminal` uses the [tty/Client](../source/ui/tty/Client.go)
- `exocomp web` spawns the [web/Server](../source/ui/web/Server.go) on port `3000` that serves the [Web UI](../source/ui/web/public/).
- `exocomp webview` spawns a [web/Server](../source/ui/web/Server.go) on port `3000` and opens a [webview/Client](../source/ui/webview/Client.go) window.


### Multi-Agent Usage

The defaulted `planner` agent is allowed to hire contracting sub-agents with the
[Agents](../source/tools/Agents.go) tool.

Multi agent communication works with a sub process hierarchy, where each process
works in their own sandbox with the `jsonl` frontend and their own agent type and
system/user prompts.

Cross-agent communication works with a strict process hierarchy, each contracted
sub-agent's `playground` is set to the parent process's `sandbox`.


#### Example Process Hierarchy

```
| Hierarchy   | Sandbox                         | Playground       | Task                          |
|:------------|:--------------------------------|:-----------------|:------------------------------|
| planner     | /path/to/project                | /path/to/project |                               |
|-> architect | /path/to/project                | /path/to/project | specifies utils package       |
|-> coder     | /path/to/project/utils          | /path/to/project | implements CalculateFibonacci |
|-> tester    | /path/to/project/utils          | /path/to/project | tests CalculateFibonacci      |
|-> coder     | /path/to/project/cmds/fibonacci | /path/to/project | implements main.go            |
```

#### Example Process Parameters

This is what the Human interacts with:

```bash
cd /path/to/project;
exocomp web planner;
```

This is what the Planner agent spawns behind the scenes as sub processes:

```bash
# cd /path/to/project;
# exocomp agent --agent=architect --prompt="Implement a utils package and specify the CalculateFibonacci method signature.";

# cd /path/to/project/utils;
# exocomp agent --agent=coder --prompt="Implement a public method called CalculateFibonacci(step int) int that calculates the fibonacci numbers.";

# cd /path/to/project/utils;
# exocomp agent --agent=tester --prompt="Implement the unit tests for the CalculateFibonacci method.";

# cd /path/to/project/cmds/fibonacci;
# exocomp agent --agent=coder --prompt="Implement a main.go that calculates fibonacci numbers. Use the first CLI parameter as the step or sequence argument.";
```

