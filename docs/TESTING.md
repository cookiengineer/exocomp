
# Testing

### Requirements

Testing the `agents` tool requires a [llama.cpp](https://github.com/ggml-org/llama.cpp)
`llama-server` instance running with a `qwen3-coder:30b` model and `Q8_0`
quantization.

For all tests, use the [unsloth Qwen3 Coder A3B Instruct](https://huggingface.co/unsloth/Qwen3-Coder-30B-A3B-Instruct-GGUF)
`Q8_0` (8-bit) quantization of the Qwen3 Coder model.

```bash
# For 48GB of VRAM use
# --ctx-size 262144

# Start llama.cpp server
llama-server \
	--model "./models/qwen3-coder-30b-a3b-instruct-q8_0.gguf" \
	--alias "qwen3-coder:30b" \
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
cd /path/to/exocomp;

# Test everything
go test -v ./...;
```


### Agents Tool

Testing the `agents` tool comes with certain limitations because of how the
`go test` workflow is designed. The `tools/Agents_main_test.go` builds a
standalone binary of `cmds/exocomp/main.go` that is injected into `tools/Agents.go`
via environment variable. This is unavoidable due to otherwise cyclic
dependencies.

```bash
# Run multi-agent unit tests
cd /path/to/exocomp/tools;

go clean -testcache;
go test -tags=agents -v ./
```

