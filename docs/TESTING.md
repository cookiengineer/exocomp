
# Testing

### Requirements

Testing the live `agents` tool requires a [llama.cpp](https://github.com/ggml-org/llama.cpp)
`llama-server` instance running. All other tests run without a server.

```bash
# For 48GB of VRAM use
# --ctx-size 262144

# Start llama.cpp server
llama-server \
	--model "Huihui-CyberStrike-OffSec-35B-abliterated-Q8_0.gguf" \
	--alias "huihui_ai/Qwen3.6-abliterated:35b" \
	--gpu-layers all \
	--ctx-size 32768 \
	--batch-size 512 \
	--ubatch-size 128 \
	--cache-type-k q8_0 \
	--cache-type-v q8_0 \
	--flash-attn auto \
	--no-slots \
	--no-webui \
	--no-webui-mcp-proxy \
	--jinja \
	--port 11434;
```

### Unit Tests

```bash
cd /path/to/exocomp/source;

# Test everything (no LLM required)
go test -v ./...;
```

This includes the deterministic agent lifecycle tests, which build a fake agent
binary ([testdata/fakeagent](../source/tools/testdata/fakeagent/main.go)) and
inject it via the `EXOCOMP_AGENT` environment variable, so the parent-side
`agents.Hire` / `agents.Await` / `agents.Fire` lifecycle is exercised without a
live LLM.

### Live Agents Tool

Testing the live `agents` tool comes with certain limitations because of how the
`go test` workflow is designed. The `tools/Agents_main_test.go` builds a
standalone binary of `cmds/exocomp/main.go` that is injected into `tools/Agents.go`
via environment variable. This is unavoidable due to otherwise cyclic
dependencies.

```bash
# Run multi-agent unit tests against a live llama-server
cd /path/to/exocomp/source/tools;

go clean -testcache;
go test -tags=agents -v ./
```
