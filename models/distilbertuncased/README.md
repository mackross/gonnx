# distilbertuncased

Embedded uncased English DistilBERT NER model package for `gonnx/models`.

## Model

- Open: `distilbertuncased.Open(gonnx.WithThreads(1))`
- Model adapter: `distilbertuncased.NER()`
- Source/export: <https://huggingface.co/onnx-community/distilbert-base-uncased-finetuned-conll03-english-ONNX>
- Source model: `elastic/distilbert-base-uncased-finetuned-conll03-english` per exported config
- Upstream ONNX revision packaged: `c6e55cbffbfaa63fe7f699ef5097cd16846f52df`
- License: Apache-2.0 per source model card at packaging time
- Labels: `O`, `B-PER`, `I-PER`, `B-ORG`, `I-ORG`, `B-LOC`, `I-LOC`, `B-MISC`, `I-MISC`
- Tokenizer: BERT uncased WordPiece (`vocab.txt`), lower-casing enabled
- Max sequence length: 512 tokens including special tokens
- Embedded ONNX asset: quantized ONNX (`assets/model_quantized.onnx`)

## Checksums

```text
a539542303839d8a421e655d35c865cf12b5ee048ee803b154d7fe473e11d9b6  assets/model_quantized.onnx
07eced375cec144d27c900241f3e339478dec958f92fddbc551f295c992038a3  assets/vocab.txt
5d20cd2ecea304cc23ed9372e3720a7ac72da8cf810b2953084b61e8a6172903  assets/config.json
```

## Reproduce assets

```sh
./prepare_assets.sh
```

## When to use

Use this model when you want DistilBERT speed/size with uncased English CoNLL labels.

Avoid it when you need case-sensitive distinctions, custom labels, or domain-specific labels.
