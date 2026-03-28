FROM python:3.12-slim

# Environment variables
ENV DEBIAN_FRONTEND=noninteractive
ENV LANG="C.UTF-8"
ENV LC_ALL="C.UTF-8"

# Target devices can be cpu, cuda, rocm

ENV PYTHONUNBUFFERED=1
ENV PIP_NO_CACHE_DIR=1

WORKDIR /setup

# Install system dependencies
RUN apt-get update && apt-get install -y git curl wget ca-certificates build-essential libssl-dev libffi-dev && apt-get clean && rm -rf /var/lib/apt/lists/*

COPY requirements.cpu.txt /setup/requirements.txt
COPY common.txt /setup/common.txt

RUN pip install --upgrade pip && pip install --extra-index-url https://download.pytorch.org/whl/cpu -r requirements.txt
RUN pip install --extra-index-url https://download.pytorch.org/whl/cpu vllm transformers accelerate

WORKDIR /app

ENV MODEL_NAME="codestral:22b"
ENV MODEL_REPO="mistralai/Codestral-22B-v0.1"
ENV VLLM_TARGET_DEVICE="cpu"

EXPOSE 5000
CMD ["vllm", "serve", "${MODEL_REPO}"]
