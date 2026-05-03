#!/usr/bin/env bash

set -euo pipefail;

# ===== Config =====
# Check https://github.com/ggml-org/llama.cpp/releases
LLAMA_VERSION="b9010"
# ==================

require_command() {

	if ! command -v "$1" >/dev/null 2>&1; then

		echo "Required command not found: $1";
		echo "Please install it first.";

		case "$1" in
			curl)
				echo "   Debian/Ubuntu: sudo apt install curl";
				echo "   Arch:          sudo pacman -S curl";
				echo "   macOS:         brew install curl";
				;;
			unzip)
				echo "   Debian/Ubuntu: sudo apt install unzip";
				echo "   Arch:          sudo pacman -S unzip";
				echo "   macOS:         brew install unzip";
				;;
		esac;

		exit 1;

	fi;

}

require_command curl;
require_command unzip;

BASE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)";
LLAMA_DIR="${BASE_DIR}/third_party/llama";
MODEL_DIR="${BASE_DIR}/third_party/models";

mkdir -p "${LLAMA_DIR}";
mkdir -p "${MODEL_DIR}";

echo "=> Installing llama.cpp version: $LLAMA_VERSION";

if [[ ! -f "${LLAMA_DIR}/llama-server" ]]; then

	LLAMA_URL="https://github.com/ggml-org/llama.cpp/releases/download/${LLAMA_VERSION}/llama-${LLAMA_VERSION}-bin-linux-vulkan.zip";
	TMP_ZIP="$(mktemp)";

	echo "=> Downloading llama.cpp ...";
	curl -fL "${LLAMA_URL}" -o "${TMP_ZIP}";
	if [ $? -ne 0 ]; then

		echo "Download of llama.cpp release $LLAMA_VERSION failed!";
		exit 1;

	else

		echo "=> Extracting llama.cpp ..."

		TMP_DIR="$(mktemp -d)";
		unzip -q "${TMP_ZIP}" -d "${TMP_DIR}";

		EXTRACTED_DIR="$(find "${TMP_DIR}" -mindepth 1 -maxdepth 1 -type d | head -n 1)";

		find "${EXTRACTED_DIR}" -type f \( \
			-name "LICENSE" -o \
			-name "*.so*" -o \
			-name "llama-server" \
			\) -exec cp {} "${LLAMA_DIR}" \;

		chmod +x "${LLAMA_DIR}/llama-server";

		rm -rf "${TMP_DIR}";
		rm -rf "${TMP_ZIP}";

		echo "=> llama.cpp installed to ${LLAMA_DIR}";

	fi;

else
	echo "=> llama.cpp already exists, skipping download";
fi;


echo "";
echo "=> Downloading gemma4:31b (Q8_0) ...";

GEMMA_MODEL_URL="https://huggingface.co/mradermacher/gemma-4-31B-GGUF/resolve/main/gemma-4-31b.Q8_0.gguf?download=true";
GEMMA_MODEL_PATH="${MODEL_DIR}/gemma4-31b-q8_0.gguf";

if [[ ! -f "${GEMMA_MODEL_PATH}" ]]; then

	curl -fL "${GEMMA_MODEL_URL}" -o "${GEMMA_MODEL_PATH}";

	if [ $? -ne 0 ]; then

		echo "!! Download failed !!";
		exit 1;

	else
		echo "=> Model installed to ${GEMMA_MODEL_PATH}";
	fi;

else
	echo "=> Model already exists, skipping download";
fi


echo "";
echo "=> Downloading qwen3-coder:30b model (Q8_0) ...";

QWEN_MODEL_URL="https://huggingface.co/unsloth/Qwen3-Coder-30B-A3B-Instruct-GGUF/resolve/main/Qwen3-Coder-30B-A3B-Instruct-Q8_0.gguf?download=true";
QWEN_MODEL_PATH="${MODEL_DIR}/qwen3-coder-30b-a3b-instruct-q8_0.gguf";

if [[ ! -f "${QWEN_MODEL_PATH}" ]]; then

	curl -fL "${QWEN_MODEL_URL}" -o "${QWEN_MODEL_PATH}";

	if [ $? -ne 0 ]; then

		echo "!! Download failed !!";
		exit 1;

	else
		echo "=> Model installed to ${QWEN_MODEL_PATH}";
	fi;

else
	echo "=> Model already exists, skipping download";
fi;


echo "";
echo "DONE.";
echo "";

