#!/usr/bin/env bash

set -euo pipefail;

OLLAMA_MODELS="${OLLAMA_MODELS:-${HOME}/.ollama/models}";
REQUIRED_MODELS=("huihui_ai/Qwen3.6-abliterated:35b" "deepseek-v4-pro:cloud");

echo "Checking Ollama models ...";

ollama_model_exists() {

	local model="$1";
	local repo="${model%%:*}";
	local tag="${model##*:}";

	# organization/model:tag
	if [[ -f "${OLLAMA_MODELS}/manifests/registry.ollama.ai/${repo}/${tag}" ]]; then
		return 0;
	fi;

	# default "library" organization
	if [[ -f "${OLLAMA_MODELS}/manifests/registry.ollama.ai/library/${repo}/${tag}" ]]; then
		return 0;
	fi;

	return 1;

}

for model in "${REQUIRED_MODELS[@]}"; do

	if ollama_model_exists "${model}"; then
		echo "> ${model} [OK]";
	else
		echo "> ${model} [NOT FOUND]";
		echo "Please execute: ollama pull \"${model}\"";
		exit 1;
	fi;

done;

AVAILABLE_KB=0;

while IFS=': ' read -r KEY VALUE _; do
	if [[ "$KEY" == "MemAvailable" ]]; then
		AVAILABLE_KB="$VALUE";
		break;
	fi;
done < /proc/meminfo;

AVAILABLE_MB=$((AVAILABLE_KB / 1024));

CONTEXT=0;

#
# ~24GB for qwen3.6/qwen3.8 model
#  ~8GB prompt cache
#
#  ~8GB for KV @ 64k
# ~16GB for KV @ 128k
#
if (( AVAILABLE_MB >= 98304)); then
	CONTEXT=131072;
elif (( AVAILABLE_MB >= 49152 )); then
	CONTEXT=65536;
elif (( AVAILABLE_MB >= 32768 )); then
	CONTEXT=32768;
elif (( AVAILABLE_MB >= 16384 )); then
	CONTEXT=16384;
elif (( AVAILABLE_MB >= 8192 )); then
	CONTEXT=8192;
else
	CONTEXT=4096;
fi

export OLLAMA_CONTEXT_LENGTH="${CONTEXT}";

echo "";
echo "Available memory:         ${AVAILABLE_MB} MiB";
echo "Available context length: ${OLLAMA_CONTEXT_LENGTH} tokens";
echo "";
echo "Starting Ollama ...";

exit 1;

exec ollama serve;

