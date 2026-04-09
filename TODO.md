
# TODO

- [ ] Changelog needs to be stored via `playground_path`
- [ ] Requirements need to be stored via `playground_path`
- [ ] sandbox is now a nested folder, which means that path is actually relative to playground
- [ ] List() and Search() methods should also show only bugs and requirements for the current sandbox

Maybe this can be easily fixed with a helper method like this:

```go
resolved, err0 := resolveSandboxPath(tool.Playground, path)

if isWithinSandbox(tool.Sandbox, resolved) {
    // old code
}
```

- [ ] Bug in Time usage. Seconds in time leads to wrong serialization, and multiple headlines
- [ ] Probably time.Time needs to be serialized to string first, and then iterated over the strings?

## Planner Agent

- [ ] Watches the completeness of the implementation phases
- [ ] Phases should be Discover, Plan, Specify, Implement, Test
- [ ] Should be a state machine
- [ ] Implement a compression between LLM agents to communicate more efficient
- [ ] https://github.com/LLM-Coding/Semantic-Anchors
- [ ] https://github.com/LLM-Coding/Semantic-Anchors/blob/main/docs/spec-driven-workflow.adoc
- [ ] https://github.com/basicmachines-co/basic-memory
- [ ] https://github.com/mrsimpson/responsible-vibe-mcp/blob/main/packages/core/src/plan-manager.ts


