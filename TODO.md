
# TODO

## types/Config

- [ ] Rewrite Config properties, Tokens should be something like "Models" struct

```json
{
    "providers": {
        "deepseek-v4-pro:cloud": {
            "url": "https://api.deepseek.com",
            "alias": "deepseek-v4-pro",
            "token": "sk-1231231232"
        },
        "huihui_ai/Qwen3.6-abliterated:35b": {
            "url": "http://localhost:11434/v1",
            "token": ""
        }
    }
```

## types/Session API Problems

- [ ] Session should handle transformations and mappings of internal "requirements.Read" syntax vs API expected "requirements_Read" syntax.
- [ ] No idea how to do this, because we don't want to change the Agents or YAML format or anything else, so it MUST be handled by the Session itself.

## types/Session DEBUGGING

- [ ] Add better debugging messages when Session is recovered successfully
- [ ] Reprint the correct config after Recovery and Backup calls
- [ ] Write debug messages that it was backed up to "./relative/path/to/exocomp/*.json"

## Questions Tool

- [ ] Implement a Question tool with `question string`, `options []string`, `multiple bool` parameters that
      shows UI elements to confirm/select options and that allows to type your own answer.

## Files Tool

- [ ] Implement `files.ReadSymbol(file, symbol)`

## Tools

- [ ] Implement Tool Call Shrinking, probably best in types/Session.go

- [ ] Implement `requirements.Trace(file, symbol)` to trace which methods
      call or interact with the named symbol. Returns a list of relative package/file.go paths and symbols

## Agents

- [ ] Check the repos in described order to find better prompts
- [ ] Binary Reverse Engineer
- [ ] Web Recon Analyst
- [ ] Web Business Analyst
- [ ] Web Subdomain Analyst
- [ ] Researcher


## Models

- [ ] ollama has 32k context length limit, even though model supports 256k
- [ ] `env OLLAMA_CONTEXT_LENGTH=256k ollama serve` doesn't change anything, needs a fix for that
- [ ] Use `https://github.com/dianlight/gollama.cpp` bindings

## Tools

- [ ] Implement unit tests for `tools/Agents`
- [ ] Implement `tools/Vulnerabilities` to be able to search CVE dataset
- [ ] Implement `tools/Websites` to be able to use `zimdex`
- [ ] Implement `tools/Skills` to parse `$PWD/skills` directory
- [ ] Validate all properties of new agents in `readAgents()`

## Web UI

- [ ] `public/ui/Renderer.mjs` should implement lazy-rendering for `nav` element to avoid setting `innerHTML`
- [ ] `public/ui/Renderer.mjs` should implement lazy-rendering for `main` element to avoid setting `innerHTML`

### Chat View

- [x] Show chats with agents

### Agents View

- [ ] Render Schedule View of agents and their gantt-like workchart
- [ ] Show inter-agent communications
- [ ] Show filesystem mutations
- [ ] Show failures (with work reports)

### Bugs View

- [ ] Show Packages sorted alphabetically in a grid
- [ ] Each grid tile shows the list of bugs for that package
- [ ] Show/Hide toggle button for showing packages with no bugs
- [ ] Create button in Footer
- [ ] Create Bug Report dialog

### Changelog View

- [ ] Show Packages sorted alphabetically in a grid
- [ ] Each grid tile shows the list of changelog entries for that package
- [ ] Show/Hide toggle button for showing packages with no changelog entries
- [ ] Create button in Footer
- [ ] Create Bug Report dialog

### Requirements View

- [ ] Show packages map of the codebase
- [ ] Show symbols and how they interact with each other
- [ ] Draw dependency lines between packages, based on imports of that package
- [ ] If possible make this map so that links between methods and "what they're calling" can be shown, too.
- [ ] Show/Hide toggle button for showing stdlib packages?

