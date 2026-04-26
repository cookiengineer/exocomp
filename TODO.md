
# TODO

## Agents

- [ ] `redteamer` who writes malware and targets victim sandboxes
- [ ] `blueteamer` who analyzes metrics and logs of victim sandboxes

## Tools

- [ ] `requirements.DefineStruct()` needs to be implemented
- [ ] `requirements.DefineTest()` needs to be implemented
- [ ] `requirements.List()` should show only requirements within current sandbox
- [ ] `requirements.Search()` should show only requirements within current sandbox

## JSONL UI

- [ ] Implement `agents` unit tests
- [ ] Implement `agents.Inquire(name)` which should summarize the status reports by using a separate summarizer agent and return the result.

## Web UI

- [ ] Implement `routes.Tools` to execute/simulate tool calls with UI elements
- [ ] `public/ui/Renderer` should implement lazy-rendering and avoid setting `nav.innerHTML = ""`
- [ ] `public/ui/Renderer` should implement lazy-rendering and avoid setting `main.innerHTML = ""`


## Planner Agent

- [ ] Watches the completeness of the implementation phases
- [ ] Phases should be Discover, Plan, Specify, Implement, Test
- [ ] Should be a state machine
- [ ] Implement a compression between LLM agents to communicate more efficient
- [ ] https://github.com/LLM-Coding/Semantic-Anchors
- [ ] https://github.com/LLM-Coding/Semantic-Anchors/blob/main/docs/spec-driven-workflow.adoc
- [ ] https://github.com/basicmachines-co/basic-memory
- [ ] https://github.com/mrsimpson/responsible-vibe-mcp/blob/main/packages/core/src/plan-manager.ts


