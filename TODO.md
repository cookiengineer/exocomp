
# TODO

## Tools

- [ ] Implement `tools/Symbols` which uses functions, classes, symbols.
      Behind the scenes a code graph and "who calls what" would complement
      the requirements tool.
- [ ] Should have a Get(package name, symbol name) method
- [ ] Should have a Search(package name, symbol name) method

## Agents

- [ ] Check the repos in described order to find better prompts
- [ ] Binary Reverse Engineer
- [ ] Web Recon Analyst
- [ ] Web Business Analyst
- [ ] Web Subdomain Analyst
- [ ] Researcher


## Models

- [ ] ollama has 32k context length limit, even though model supports 256k
- [ ] `env OLLAMA_CONTEXT_LENGTH=256k ollama serve` doesn't change shit
- [ ] Use `https://github.com/dianlight/gollama.cpp` bindings
- [ ] Implement correct llama.cpp templates for `qwen3-coder:30b` (tool calls don't work)
- [ ] Implement correct llama.cpp templates for `gemma4:31b` (tool calls don't work)
- [ ] Implement correct llama.cpp templates for `qwen3.6:35b-heretic` (tool calls don't work)

## Tools

- [ ] Implement unit tests for `tools/Agents`
- [ ] Implement `tools/Containers` to be able to use `podman`
- [ ] Implement `tools/Websites` to be able to use `zimdex`
- [ ] Implement `tools/Skills` to parse `$PWD/skills` directory

## JSONL UI

- [ ] Implement `agents` unit tests

## Web UI

- [ ] `public/ui/Renderer.mjs` should implement lazy-rendering for `nav` element to avoid setting `innerHTML`
- [ ] `public/ui/Renderer.mjs` should implement lazy-rendering for `main` element to avoid setting `innerHTML`


## Agent History Compression

- [ ] Implement a compression between LLM agents to communicate more efficient
- [ ] Might be optimum use case for summarizer inventing its own language
- [ ] https://github.com/LLM-Coding/Semantic-Anchors
- [ ] https://github.com/LLM-Coding/Semantic-Anchors/blob/main/docs/spec-driven-workflow.adoc
- [ ] https://github.com/juliusbrussee/caveman
- [ ] https://github.com/mrsimpson/responsible-vibe-mcp/blob/main/packages/core/src/plan-manager.ts


