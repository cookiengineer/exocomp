# Testing

## Unit Tests

Run the deterministic test suite from the repository root:

```bash
make test;
```

`make test` runs `go -C source test ./...` and needs no network or LLM. It
includes the agent lifecycle tests, which build the fake agent in
[`source/tools/testdata/fakeagent`](../source/tools/testdata/fakeagent/main.go)
and inject it through `EXOCOMP_AGENT`.

The toolchain module has its own tests:

```bash
make test-toolchain;
```

Use the targeted commands below only for checks that `make test` does not cover.

The `programs` and `skills` execution tests need Linux with landlock or macOS
with `/usr/bin/sandbox-exec`. They build `exocomp-sandbox` once and inject it
through `EXOCOMP_SANDBOX`. Without a usable platform sandbox those tests are
skipped.

The `containers` tests use the fake podman binary in
[`source/tools/testdata/fakepodman`](../source/tools/testdata/fakepodman/main.go)
and need no real podman. A real integration test is skipped by default:

```bash
EXOCOMP_PODMAN_TEST=1 go -C source/tools test -run TestContainers_Integration -v ./containers;
```

It needs a working rootless podman installation and network access.

## Live Agent Tests

The live suite is behind `//go:build agents` and is not part of `make test`:

```bash
go -C source/tools/agents test -tags=agents -v ./;
```

It is a local development workflow. `Agents_main_test.go` builds the child
binary and starts a llama-server from `source/tools/third_party`, so review the
paths at the top of that file before running it. The expected model is
`huihui_ai/Qwen3.8-abliterated:27b` at `Q4_K_M` quantization. See
[SERVERS.md](./SERVERS.md) for server parameters.

## Writing Tests

See [TESTING_GUIDE.md](./TESTING_GUIDE.md) for file naming, setup, fixtures,
assertions and the offline fake-binary conventions.
