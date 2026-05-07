
# TODO

## Agents

- [ ] Check the repos in described order to find better prompts
- [ ] Web Pentester
- [ ] Binary Reverse Engineer
- [ ] Web Recon Analyst
- [ ] Web Business Analyst
- [ ] Web Subdomain Analyst
- [ ] Researcher


## Models

- [ ] ollama has 32k context length limit, even though model supports 256k
- [ ] `env OLLAMA_CONTEXT_LENGTH=256k ollama serve`
- [ ] Use `https://github.com/dianlight/gollama.cpp` bindings
- [ ] Implement `models/LLM.go` for inference tasks
- [ ] Implement templating for both qwen3-coder and gemma4 models
- [ ] Figure out what format gollama.cpp needs and if it can load GGUF files


## Agents

- [ ] `redteamer` who writes malware and targets victim sandboxes
- [ ] `blueteamer` who analyzes metrics and logs of victim sandboxes

## Tools

- [ ] Implement unit tests for `tools/Agents`
- [ ] Implement `tools/Containers` to be able to use `podman`
- [ ] Implement `tools/Websites` to be able to use `zimdex`
- [ ] Implement `tools/Skills` to parse `$PWD/skills` directory

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


