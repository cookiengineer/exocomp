#!/bin/bash

./llama/llama-server \
	-m ./models/qwen3-coder-30b-a3b-instruct-q8_0.gguf \
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

# Needs 48GB VRAM GPU
# --threads 16;
# --batch-size 2048
# --ubatch-size 512
