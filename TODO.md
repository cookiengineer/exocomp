
# TODO

## source/engine/Session

- No retry/backoff. Session.infer_chat_completions does a single POST; transient 429/5xx and connection errors just fail the whole turn.
- No finish_reason handling. schemas.Choice.FinishReason is parsed but never read. A length-truncated response with an incomplete tool_calls block is silently accepted. mini re-prompts with guidance.
- No format-error feedback loop. In ReceiveChatResponse, an assistant message with no tool_calls just returns nil (dead-end), and malformed tool-call args are silently skipped (err0/err1… != nil → ignored).
- Choices need to be correctly handled, currently exocomp handles only Choices[0]

## source/schemas and source/engine/Session

- Context Window Size and Token Budget calculations
- No cost tracking (only token counts). Needs to support typical API providers

## source/tools

- vulnerabilities.* tool to search CVEs
- findings.* tool to store pentest findings in a strictly defined schema to finally generate a PDF report
- requirements.Trace(file, symbol) to trace callers/callees of a Symbol (which would feed the Requirements graph view in the frontend).

## Agents Workflow

- [ ] agents.Hire() needs to be more clear for DeepSeek planner model
- [ ] In manual tests something went wrong with the subprocess' `--prompt` flag. It might
      be that the `--prompt` parameter for a subagent process doesn't render it then as a
      `user` role message sent to the LLM/assistant?

## Tools

- [ ] Implement Tool Call Shrinking, meaning to remove failed tool calls if later ones with
      the same parameters were successful, so that smaller LLM context sizes don't get flooded
      with polling tool call error messages. It's probably best to implement this in types/Session.go
- [ ] Implement `requirements.Trace(file, symbol)` to trace which methods
      call or interact with the named symbol. Returns a list of relative package/file.go paths and symbols
- [ ] Implement `tools/Vulnerabilities` to be able to search CVE dataset
- [ ] Implement `tools/Websites.Search()` to be able to search the internet. Don't use google, find an LLM friendly free search API provider.
- [ ] Validate all properties of new agents in `readAgents()`

## Agents

- [ ] Check the repos in described order to find better prompts
- [ ] Binary Reverse Engineer
- [ ] Web Recon Analyst
- [ ] Web Business Analyst
- [ ] Web Subdomain Analyst
- [ ] Researcher


## Web UI

- [ ] Migrate towards use of a WebSocket so that we don't need to poll /api/session/agents all the time
- [ ] It would be really cool if server had a `ui/web/sockets/session/Agents` kind of wrapper for the `/api/session/agents` route

- [ ] `public/ui/Renderer.mjs` should implement lazy-rendering for `nav` element to avoid setting `innerHTML`
- [ ] `public/ui/Renderer.mjs` should implement lazy-rendering for `main` element to avoid setting `innerHTML`

### Chat View

- [x] Show chats with agents

### Agents View

- [ ] Agents view needs to display tool calls better, and inter-agent communications, and filesystem mutations
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
- [ ] Create Changelog Report dialog (e.g. Add/Change/Deprecate/Fix/Remove)

### Requirements View

- [ ] Show packages map of the codebase
- [ ] Show symbols and how they interact with each other
- [ ] Draw dependency lines between packages, based on imports of that package
- [ ] If possible make this map so that links between methods and "what they're calling" can be shown, too.
- [ ] Show/Hide toggle button for showing stdlib packages?

