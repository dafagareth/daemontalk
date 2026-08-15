---
title: "Running LLMs Locally with llama.cpp and Ollama"
slug: d2c5e837
aliases: ai-inference-local-llm
date: 2026-07-01
tags: [tools, linux, performance]
lang: en
draft: false
---

Large language models are no longer exclusive to cloud APIs. With llama.cpp and Ollama, you can run capable quantized models on a consumer laptop or a small home server. The performance is practical for interactive use: a 7B parameter model running on a modern CPU can produce 20-40 tokens per second, which is faster than reading speed. This post covers the technical details of the GGUF format, quantization options, hardware planning, and running a model as a persistent service.

## Fun Facts

**Fact 1.** llama.cpp was created by Georgi Gerganov and first published in March 2023, just weeks after Meta released the original LLaMA weights. It is written in C and C++ with no mandatory Python dependency and runs on CPU, CUDA, Metal, Vulkan, and ROCm backends.

**Fact 2.** The GGUF format replaced the earlier GGML format in August 2023. GGUF is a self-contained binary format that stores model weights, tokenizer vocabulary, and metadata in a single file with a versioned magic number, making it forward and backward compatible across llama.cpp versions.

**Fact 3.** Quantization reduces model weights from 16-bit or 32-bit floats to 4-bit or 8-bit integers using calibration data. The quality loss at Q4_K_M is small enough that most benchmarks show less than a 2-point drop on common reasoning benchmarks compared to the full FP16 model.

---

## Tips and Tricks

### 1. Understand GGUF Quantization Levels

GGUF files are available in several quantization levels. Choosing the right one depends on your available RAM and the quality you need.

| Quantization | Bits/Weight | 7B Model Size | Notes |
|---|---|---|---|
| F16 | 16 | ~14 GB | Full precision, highest quality |
| Q8_0 | 8 | ~7.7 GB | Near-lossless, good for testing |
| Q4_K_M | 4 | ~4.1 GB | Best quality-to-size ratio |
| Q4_0 | 4 | ~3.8 GB | Faster than K_M, slightly lower quality |
| Q2_K | 2 | ~2.7 GB | Noticeable quality degradation |

```bash
# Download a model from HuggingFace using huggingface-cli
pip install huggingface_hub
huggingface-cli download \
  bartowski/Meta-Llama-3.1-8B-Instruct-GGUF \
  Meta-Llama-3.1-8B-Instruct-Q4_K_M.gguf \
  --local-dir ~/models/

# Check file size
ls -lh ~/models/*.gguf

# Inspect GGUF metadata
python3 - << 'EOF'
import struct, sys
with open(sys.argv[1], "rb") as f:
    magic = f.read(4)
    version = struct.unpack("<I", f.read(4))[0]
    print(f"GGUF version: {version}")
EOF ~/models/Meta-Llama-3.1-8B-Instruct-Q4_K_M.gguf
```

### 2. Build and Run llama.cpp Directly

Building llama.cpp from source gives you the most control and access to the latest optimizations.

```bash
# Clone and build with CPU support
git clone https://github.com/ggerganov/llama.cpp
cd llama.cpp
cmake -B build -DLLAMA_CURL=ON
cmake --build build --config Release -j$(nproc)

# Build with CUDA support (requires CUDA toolkit)
cmake -B build -DGGML_CUDA=ON
cmake --build build --config Release -j$(nproc)

# Run inference
./build/bin/llama-cli \
  --model ~/models/Meta-Llama-3.1-8B-Instruct-Q4_K_M.gguf \
  --ctx-size 4096 \
  --n-predict 200 \
  --prompt "Explain the difference between TCP and UDP in two sentences."

# Run the OpenAI-compatible server
./build/bin/llama-server \
  --model ~/models/Meta-Llama-3.1-8B-Instruct-Q4_K_M.gguf \
  --ctx-size 4096 \
  --port 8080 \
  --host 0.0.0.0
```

### 3. Use Ollama for a Simpler Workflow

Ollama wraps llama.cpp with a model registry, automatic downloads, and a REST API compatible with the OpenAI client library.

```bash
# Install Ollama
curl -fsSL https://ollama.com/install.sh | sh

# Pull and run a model
ollama pull llama3.1:8b
ollama run llama3.1:8b "What is the capital of Vietnam?"

# List downloaded models
ollama list

# Run the Ollama API server (starts automatically with service)
curl http://localhost:11434/api/generate -d '{
  "model": "llama3.1:8b",
  "prompt": "What is 17 * 23?",
  "stream": false
}' | python3 -m json.tool

# Use with the OpenAI Python client
python3 - << 'EOF'
from openai import OpenAI
client = OpenAI(base_url="http://localhost:11434/v1", api_key="ollama")
resp = client.chat.completions.create(
    model="llama3.1:8b",
    messages=[{"role": "user", "content": "Explain GGUF in one sentence."}]
)
print(resp.choices[0].message.content)
EOF
```

### 4. Plan Hardware: RAM vs VRAM

The model must fit entirely in memory to run without swapping. GPU inference is significantly faster than CPU inference because of memory bandwidth.

```bash
# Calculate required RAM for a model
# Formula: (parameters in billions * bits_per_weight / 8) + ~10% overhead
# Example: 8B model at Q4_K_M = (8 * 4.5 / 8) * 1.1 ≈ 4.95 GB

# Check available system RAM
free -h

# Check VRAM on NVIDIA GPU
nvidia-smi --query-gpu=memory.total,memory.free --format=csv

# Check VRAM on AMD GPU
rocm-smi --showmeminfo vram

# With Ollama: force CPU-only (useful for testing without GPU)
OLLAMA_NUM_GPU=0 ollama run llama3.1:8b

# With llama.cpp: split model across GPU and CPU (for partial VRAM fit)
./build/bin/llama-cli \
  --model ~/models/Meta-Llama-3.1-8B-Instruct-Q4_K_M.gguf \
  --n-gpu-layers 20 \
  --ctx-size 4096 \
  --prompt "Hello"
# --n-gpu-layers controls how many transformer layers go on GPU
```

### 5. Run Ollama as a systemd Service and Benchmark

For a persistent local inference endpoint, run Ollama under systemd. The installer already creates a service unit, but you can customize it.

```bash
# Check the installed systemd unit
cat /etc/systemd/system/ollama.service

# Override environment variables (e.g., bind to a specific interface)
sudo systemctl edit ollama
# Add under [Service]:
# Environment="OLLAMA_HOST=127.0.0.1:11434"
# Environment="OLLAMA_NUM_PARALLEL=2"

# Reload and restart
sudo systemctl daemon-reload
sudo systemctl restart ollama
sudo systemctl status ollama

# Benchmark tokens per second using llama.cpp bench tool
./build/bin/llama-bench \
  --model ~/models/Meta-Llama-3.1-8B-Instruct-Q4_K_M.gguf \
  --n-prompt 512 \
  --n-gen 128 \
  --threads $(nproc)

# Quick manual benchmark with Ollama
time echo "Write a 100-word paragraph about the Linux kernel." \
  | ollama run llama3.1:8b
```

The `llama-bench` output shows prompt processing speed (pp) and generation speed (tg) separately. On a modern 8-core CPU with AVX2, a Q4_K_M 8B model typically achieves 15-25 tokens/sec generation. A consumer NVIDIA GPU with 8 GB VRAM brings this to 60-100 tokens/sec for the same model.
