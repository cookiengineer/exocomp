
# Testing

```bash
cd /path/to/exocomp;

# Test everything
go test -v ./...;
```

## tools/Agents

Testing the `agents` tool comes with certain limitations because of how the
`go test` workflow is designed. The `tools/Agents_main_test.go` builds a
standalone binary of `cmds/exocomp/main.go` that is injected into `tools/Agents.go`
via environment variable.

This is unavoidable due to otherwise existing cyclic dependencies and/or
limitations of how the `TestMain(m *testing.M)` method is only called by
the test compilat once and not the actual test packages, therefore the
`os.Args` and `os.Executable()` will always be the wrong parameters.

```bash
cd /path/to/exocomp/tools;

# Test agents Tool
go clean -testcache;
go test -tags=agents -v ./;
```
