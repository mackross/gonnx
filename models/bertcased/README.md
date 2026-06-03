# bertcased

Embedded cased English BERT-base NER model package for `gonnx/models`.

## Model

- Open: `bertcased.Open(gonnx.WithThreads(1))`
- Model adapter: `bertcased.BaseCased()`
- Source/export: <https://huggingface.co/onnx-community/bert-base-NER-ONNX>
- Source model: <https://huggingface.co/dslim/bert-base-NER>
- Upstream ONNX revision packaged: `9faa2f4a2d59b396888b318f596ff719cc893f1e`
- Source model revision observed: `d1a3e8f13f8c3566299d95fcfc9a8d2382a9affc`
- License: MIT per source/export model card at packaging time
- Labels: `O`, `B-MISC`, `I-MISC`, `B-PER`, `I-PER`, `B-ORG`, `I-ORG`, `B-LOC`, `I-LOC`
- Tokenizer: BERT cased WordPiece (`vocab.txt`), lower-casing disabled
- Max sequence length: 512 tokens including special tokens
- Embedded ONNX asset: quantized ONNX (`assets/model_quantized.onnx`)

## Checksums

```text
b324e829f1fad3b897f926d1a1d1372803c6d04546831a3a2ed103652b916adf  assets/model_quantized.onnx
eeaa9875b23b04b4c54ef759d03db9d1ba1554838f8fb26c5d96fa551df93d02  assets/vocab.txt
e5be6d68d7b6a7af53b759a4b7a3f57332100009ab1e798d70a7c8d8b90dcf3d  assets/config.json
```

## Reproduce assets

```sh
./prepare_assets.sh
```

## When to use

Use this model when cased English entity recognition matters and BERT-base runtime/binary cost is acceptable.

Avoid it when you need uncased normalization, smaller binaries, or multilingual/domain-specific labels.
