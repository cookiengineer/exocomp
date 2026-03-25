FROM python:3.12-slim

# Environment variables
ENV DEBIAN_FRONTEND=noninteractive
ENV LANG="C.UTF-8"
ENV LC_ALL="C.UTF-8"
ENV MODEL_NAME="codestral:22b"
ENV MODEL_REPO="mistralai/Codestral-22B-v0.1"

# Install system dependencies
RUN apt-get update && apt-get install -y git curl wget ca-certificates build-essential libssl-dev libffi-dev && apt-get clean && rm -rf /var/lib/apt/lists/*

# Upgrade pip
# RUN python3 -m pip install --upgrade pip

# RUN python3 -m pip install torch torchvision --index-url https://download.pytorch.org/whl/cpu
RUN python3 -m pip install --upgrade vllm transformers
RUN python3 -m pip install accelerate

WORKDIR /workspace

EXPOSE 5000

# Optional: pre-cache model for faster startup (can take a while for 22B)
# RUN python3 -c "from transformers import AutoTokenizer, AutoModelForCausalLM; AutoTokenizer.from_pretrained('${MODEL_REPO}'); AutoModelForCausalLM.from_pretrained('${MODEL_REPO}', device_map='auto')"

CMD ["bash", "-c", "python3 -m vllm.entrypoints.http_server --model ${MODEL_REPO} --port 5000 --dtype int8"]
