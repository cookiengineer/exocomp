#!/usr/bin/env bash

set -euo pipefail;

# ===== Config =====
# Check https://github.com/ggml-org/llama.cpp/releases
LLAMA_VERSION="b9010"
# ==================
BASE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)";
LLAMA_DIR="${BASE_DIR}/third_party/llama";
MODELS_DIR="${BASE_DIR}/third_party/models";
PROGRAMS_DIR="${BASE_DIR}/third_party/programs";
# ==================


install_program() {

	cmd_url="${1}";

	if [[ "${cmd_url}" == github.com/* ]]; then

		cd "${PROGRAMS_DIR}";
		env GOBIN="${PROGRAMS_DIR}" go install -v "${cmd_url}"

	fi;

}

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
			env)
				echo "   Debian/Ubuntu: sudo apt install coreutils";
				echo "   Arch:          sudo pacman -S coreutils";
				echo "   macOS:         brew install coreutils";
				;;
			go)
				echo "   Debian/Ubuntu: sudo apt install golang-go";
				echo "   Arch:          sudo pacman -S go";
				echo "   macOS:         brew install go";
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
require_command env;
require_command go;
require_command unzip;



mkdir -p "${LLAMA_DIR}";
mkdir -p "${MODELS_DIR}";
mkdir -p "${PROGRAMS_DIR}";



echo "";
echo "=> Installing amass";

if [[ ! -f "${PROGRAMS_DIR}/amass" ]]; then
	install_program "github.com/owasp-amass/amass/v5/cmd/amass@latest";
else
	echo "=> amass already exists";
fi;

echo "";
echo "=> Installing asnmap";

if [[ ! -f "${PROGRAMS_DIR}/asnmap" ]]; then
	install_program "github.com/projectdiscovery/asnmap/cmd/asnmap@latest";
else
	echo "=> asnmap already exists";
fi;

echo "";
echo "=> Installing httpx";

if [[ ! -f "${PROGRAMS_DIR}/httpx" ]]; then
	install_program "github.com/projectdiscovery/httpx/cmd/httpx@latest";
else
	echo "=> httpx already exists";
fi;

echo "";
echo "=> Installing nuclei";

if [[ ! -f "${PROGRAMS_DIR}/nuclei" ]]; then
	install_program "github.com/projectdiscovery/nuclei/v3/cmd/nuclei@latest";
else
	echo "=> nuclei already exists";
fi;

echo "";
echo "=> Installing katana";

if [[ ! -f "${PROGRAMS_DIR}/katana" ]]; then
	install_program "github.com/projectdiscovery/katana/cmd/katana@latest";
else
	echo "=> katana already exists";
fi;

echo "";
echo "=> Installing naabu";

if [[ ! -f "${PROGRAMS_DIR}/naabu" ]]; then
	install_program "github.com/projectdiscovery/naabu/v2/cmd/naabu@latest";
else
	echo "=> naabu already exists";
fi;

echo "";
echo "=> Installing subfinder";

if [[ ! -f "${PROGRAMS_DIR}/subfinder" ]]; then
	install_program "github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest";
else
	echo "=> subfinder already exists";
fi;

echo "";
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
	echo "=> llama.cpp already exists";
fi;

echo "";
echo "=> Downloading gemma4:31b (Q8_0) ...";

GEMMA_MODEL_URL="https://huggingface.co/mradermacher/gemma-4-31B-GGUF/resolve/main/gemma-4-31b.Q8_0.gguf?download=true";
# GEMMA_MODEL_FILE="gemma4-31b-q8_0.gguf";
GEMMA_MODEL_FILE="gemma4:31b.gguf"; # XXX: ollama compatible name
GEMMA_MODEL_PATH="${MODELS_DIR}/${GEMMA_MODEL_FILE}";

if [[ ! -f "${GEMMA_MODEL_PATH}" ]]; then

	curl -fL "${GEMMA_MODEL_URL}" -o "${GEMMA_MODEL_PATH}";

	if [ $? -ne 0 ]; then

		echo "!! Download failed !!";
		exit 1;

	else
		echo "=> Model installed to ${GEMMA_MODEL_PATH}";
	fi;

else
	echo "=> Model already exists";
fi

echo "";
echo "=> Downloading qwen3-coder:30b model (Q8_0) ...";

QWEN_MODEL_URL="https://huggingface.co/unsloth/Qwen3-Coder-30B-A3B-Instruct-GGUF/resolve/main/Qwen3-Coder-30B-A3B-Instruct-Q8_0.gguf?download=true";
# QWEN_MODEL_FILE="qwen3-coder-30b-a3b-instruct-q8_0.gguf";
QWEN_MODEL_FILE="qwen3-coder:30b.gguf"; # XXX: ollama compatible name
QWEN_MODEL_PATH="${MODELS_DIR}/${QWEN_MODEL_FILE}";

if [[ ! -f "${QWEN_MODEL_PATH}" ]]; then

	curl -fL "${QWEN_MODEL_URL}" -o "${QWEN_MODEL_PATH}";

	if [ $? -ne 0 ]; then

		echo "!! Download failed !!";
		exit 1;

	else
		echo "=> Model installed to ${QWEN_MODEL_PATH}";
	fi;

else
	echo "=> Model already exists";
fi;

echo "";
echo "=> Downloading qwen3.6:35b model (Q8_0) ...";

QWEN36_MODEL_URL="https://huggingface.co/llmfan46/Qwen3.6-35B-A3B-uncensored-heretic-GGUF/resolve/main/Qwen3.6-35B-A3B-uncensored-heretic-Q8_0.gguf?download=true";
#QWEN36_MODEL_FILE="qwen3.6-35b-a3b-heretic-q8_0.gguf";
QWEN36_MODEL_FILE="qwen3.6:35b.gguf"; # XXX: ollama compatible name
QWEN36_MODEL_PATH="${MODELS_DIR}/${QWEN_MODEL_FILE}";

if [[ ! -f "${QWEN36_MODEL_PATH}" ]]; then

	curl -fL "${QWEN36_MODEL_URL}" -o "${QWEN36_MODEL_PATH}";

	if [ $? -ne 0 ]; then

		echo "!! Download failed !!";
		exit 1;

	else
		echo "=> Model installed to ${QWEN36_MODEL_PATH}";
	fi;

else
	echo "=> Model already exists";
fi;



echo "";
echo "DONE.";
echo "";

