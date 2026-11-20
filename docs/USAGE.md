# Usage

The user-facing entry points are:

| Binary | Modes | Interface |
|:-------|:------|:----------|
| `exocomp` | `config`, `terminal` (aliases `tty`, `term`) | TTY frontend |
| `exocomp-web` | `config`, `web` | Web frontend on port `3000` |
| `exocomp-agent` | flags only | Line-delimited JSON over `stdin`/`stdout`; no UI |

Common flags are parsed by [`cli.ParseConfig`](../source/utils/cli/ParseConfig.go):
`--name`, `--role`, `--model`, `--temperature`, `--prompt`, `--playground`,
`--sandbox`, `--url` and `--debug`.

`--sandbox` must be equal to or below `--playground`. Both default to the
current working directory. `--temperature` accepts `0.1` through `1.0`.
`--debug` works as a bare flag or as `--debug=true`.

## Single-Agent Usage

```bash
exocomp terminal --role=architect --model="codestral:22b" --temperature=0.5;
exocomp-web web --role=planner --model="codestral:22b" --temperature=0.7;
```

The default role is `planner`.

## Multi-Agent Usage

The planner can hire contractor subprocesses through the
[Agents](../source/tools/agents/Agents.go) tool. Every child runs
`exocomp-agent` with its own sandbox and inherits the parent's playground.

```text
planner        sandbox /path/to/project                 playground /path/to/project
|-> architect  sandbox /path/to/project                 playground /path/to/project
|-> coder      sandbox /path/to/project/utils           playground /path/to/project
|-> tester     sandbox /path/to/project/utils           playground /path/to/project
|-> coder      sandbox /path/to/project/cmds/fibonacci  playground /path/to/project
```

The planner spawns the children behind the scenes:

```bash
cd /path/to/project;
exocomp-web web --role=planner;

# child processes
exocomp-agent --role="architect" --playground="/path/to/project" --sandbox="/path/to/project" --prompt="Implement a utils package and specify the CalculateFibonacci method signature.";
exocomp-agent --role="coder" --playground="/path/to/project" --sandbox="/path/to/project/utils" --prompt="Implement a public method called CalculateFibonacci(step int) int that calculates the fibonacci numbers.";
exocomp-agent --role="tester" --playground="/path/to/project" --sandbox="/path/to/project/utils" --prompt="Implement the unit tests for the CalculateFibonacci method.";
exocomp-agent --role="coder" --playground="/path/to/project" --sandbox="/path/to/project/cmds/fibonacci" --prompt="Implement a main.go that calculates fibonacci numbers. Use the first CLI parameter as the step or sequence argument.";
```

In the TTY frontend, the left sidebar lists the planner and every hired agent.
`Ctrl+Up/Down` switches the displayed history and `Alt+Up/Down` scrolls the
transcript. The composer follows the selected agent, so a hired agent can be
helped manually. Messages sent while it is busy are queued behind its current
request. `agents.Fire` is the hard stop.

The planner looks for `exocomp-agent` next to its own binary first, then on
`PATH`. Set `EXOCOMP_AGENT` to override the lookup during development.
