
# TODO

## cmds/agimus

- [ ] Verify that playground and sandbox are set to correct `/tmp/agimus-*` folder
- [ ] Verify that internal session backup is stored in correct `/tmp` folder and doesn't override source folder

## Agents Workflow

- [ ] agents.Hire() needs to be more clear for DeepSeek planner model
- [ ] "user" prompt isn't being injected into the subprocess. Might be a Config issue because config.Prompt isn't user prompt but a system prompt!?

## Questions Workflow

- [ ] Questions should be rendered as Tool Calls in the Web UI
      If they're unanswered (no following `tool` Type Message for the same toolcall id)
      then render a UI for answering them _inside_ the tool call message's article element
      If they're answered, then render a UI with disabled elements for now.


## engine/Session

- [ ] Add better debugging messages when Session is recovered and backed up
- [ ] Print the correct config after Recovery and Backup calls

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

