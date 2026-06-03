#!/usr/bin/env bash
set -euo pipefail
root="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
out="$root/assets"
base="https://huggingface.co/onnx-community/NeuroBERT-NER-ONNX/resolve/main"
mkdir -p "$out"
curl -L --fail "$base/onnx/model_quantized.onnx" -o "$out/model_quantized.onnx"
curl -L --fail "$base/vocab.txt" -o "$out/vocab.txt"
curl -L --fail "$base/config.json" -o "$out/config.json"
(cd "$root" && sha256sum -c checksums.txt)
