
# Inference Servers

Exocomp supports the following inference servers:

- [llama.cpp](../utils/api/llamacpp)
- [ollama](../utils/api/ollama)
- [vllm](../utils/api/vllm)

Support for a new inference server requires the following API endpoints:

- `http://server:port/v1/chat/completions` support
- `http://server:port/v1/models` support
- [types.Config](../types/Config.go) `GetContextLength(model)` support via custom `utils/api` adapter

## ollama Usage

[ollama](https://ollama.com) has been tested and can be used as an inference server.

Known limitations:

- `32k` context length and unavoidable input pruning, despite environment variable and request options field.

```bash
# Start and setup ollama
ollama serve;

# Pull models
ollama pull qwen3-coder:30b;
ollama pull gemma4:31b;
```

```bash
cd /path/to/exocomp;

# Run exocomp with external server
go run ./cmds/exocomp/main.go web planner --url="http://ollama_server_ip:11434/v1";
```

## llama.cpp Usage

[llama.cpp](https://github.com/ggml-org/llama.cpp) has been tested and can be used as an inference server.

Recommended hardware requirements:

- `48GB VRAM GPU` with at least `Vulkan` support
- `80GB HDD space` for models

```bash
cd /path/to/exocomp;

# Install llama.ccp and LLM models
bash install-deps.sh;

# Run exocomp with internal server
go run -tags=with_llamacpp ./cmds/exocomp/main.go web planner;
```

Full llama.cpp server parameters:

```bash
cd /path/to/exocomp/third_party;

# On systems with less than 48GB of VRAM
./llama/llama-server \
	--models-dir ./models \
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
    --port=11434;

# On systems with at least 48GB of VRAM
./llama/llama-server \
	--models-dir ./models \
	--gpu-layers all \
	--ctx-size 262144 \
	--batch-size 2048 \
	--ubatch-size 512 \
	--cache-type-k q8_0 \
	--cache-type-v q8_0 \
	--flash-attn auto \
	--no-slots \
	--no-webui \
	--no-webui-mcp-proxy \
	--jinja \
    --port=11434;
```
