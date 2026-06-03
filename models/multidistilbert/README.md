# multidistilbert

Embedded multilingual DistilBERT NER model package for `gonnx/models`.

## Model

- Open: `multidistilbert.Open(gonnx.WithThreads(1))`
- Model adapter: `multidistilbert.NER()`
- Source/export: <https://huggingface.co/Xenova/distilbert-base-multilingual-cased-ner-hrl>
- Source model: <https://huggingface.co/Davlan/distilbert-base-multilingual-cased-ner-hrl>
- Upstream ONNX revision packaged: `c2a4dbf593c57f47004c5bc2d3770d311aee9c43`
- Source model revision observed: `d421f57d5b1d36b375408588669e9340f9b11a89`
- License: AFL-3.0 per source model card at packaging time
- Labels: `O`, `B-DATE`, `I-DATE`, `B-PER`, `I-PER`, `B-ORG`, `I-ORG`, `B-LOC`, `I-LOC`
- Tokenizer: multilingual cased WordPiece (`vocab.txt`), lower-casing disabled
- Max sequence length: 512 tokens including special tokens
- Embedded ONNX asset: quantized ONNX (`assets/model_quantized.onnx`)

## Checksums

```text
24a0b98f4dd4cd92842f5a541272f86f760225a64a29928eddef14bdb2edb986  assets/model_quantized.onnx
fe0fda7c425b48c516fc8f160d594c8022a0808447475c1a7c6d6479763f310c  assets/vocab.txt
38847be4dc6699b1218a749ed69f888c2ccc7b4deba98e3c4a1cac8cb34d54c8  assets/config.json
```

## Reproduce assets

```sh
./prepare_assets.sh
```

## When to use

Use this model when you need multilingual cased NER with generic `DATE`, `PER`, `ORG`, and `LOC` labels.

Avoid it when you need a small binary, English-only best accuracy, PII-specific labels, or domain-specific labels.
