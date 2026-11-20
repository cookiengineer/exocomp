# Servers

Exocomp supports these inference servers through the OpenAI-compatible API:

- [llama.cpp](../source/utils/api/llamacpp)
- [ollama](../source/utils/api/ollama)
- [vLLM](../source/utils/api/vllm)

A server must provide:

- `http://server:port/v1/chat/completions`
- `http://server:port/v1/models`
- context-length discovery through the matching `utils/api` adapter

## ollama

Start ollama and pull the models:

```bash
ollama serve;
ollama pull deepseek-v4-pro:cloud;
ollama pull huihui_ai/Qwen3.8-abliterated:27b;
```

`quickstart-ollama.sh` sets `OLLAMA_CONTEXT_LENGTH` from available memory. Start
Exocomp against the ollama endpoint:

```bash
cd /path/to/exocomp;
make quick;

cd /path/to/project;
/path/to/exocomp/build/<operating-system>/exocomp-web web --role=planner --url="http://ollama_server_ip:11434/v1";
```

## llama.cpp

Install `llama-server` and the model separately. Exocomp does not download
either. See the [llama.cpp repository](https://github.com/ggml-org/llama.cpp).

Build Exocomp and start it with the default local endpoint:

```bash
cd /path/to/exocomp;
make quick;

cd /path/to/project;
/path/to/exocomp/build/<operating-system>/exocomp-web web --role=planner;
```

Server parameters for a machine with less than 48 GB of VRAM:

```bash
llama-server \
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
	--port 11434;
```

For at least 48 GB of VRAM, raise `--ctx-size` to `262144`, `--batch-size` to
`2048` and `--ubatch-size` to `512`.

## vLLM

vLLM is supported through the same OpenAI-compatible endpoints. The example
container definition and Python requirements are in [`docs/docker`](./docker/).
The example container listens on `5000`; a default vLLM server uses `8000`.

```bash
cd /path/to/project;
/path/to/exocomp/build/<operating-system>/exocomp-web web --role=planner --url="http://vllm_server_ip:5000/v1";
```
