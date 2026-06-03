# neurobert

Embedded tiny English BERT NER model package for `gonnx/models`.

## Model

- Open: `neurobert.Open(gonnx.WithThreads(1))`
- Model adapter: `neurobert.NER()`
- Source/export: <https://huggingface.co/onnx-community/NeuroBERT-NER-ONNX>
- Source model: <https://huggingface.co/boltuix/NeuroBERT-NER>
- Upstream ONNX revision packaged: `e3a6a290e438a506d2ba7bbb0f5875b7f034cebf`
- Source model revision observed: `fdc13cd31d4abc10b1b54df348caf10727a1ad93`
- Original/base model: `boltuix/NeuroBERT-Mini`
- License: Apache-2.0 per source model card at packaging time
- Labels: `O` plus BIO labels for `CARDINAL`, `DATE`, `EVENT`, `FAC`, `GPE`, `LANGUAGE`, `LAW`, `LOC`, `MONEY`, `NORP`, `ORDINAL`, `ORG`, `PERCENT`, `PERSON`, `PRODUCT`, `QUANTITY`, `TIME`, and `WORK_OF_ART`
- Tokenizer: BERT uncased WordPiece (`vocab.txt`), lower-casing enabled
- Max sequence length: 512 tokens including special tokens
- Embedded ONNX asset: quantized ONNX (`assets/model_quantized.onnx`)

## Checksums

```text
17dbf6ccfe500ee8a5da9177ebcdbb1c7394626e077e70571adebae8833c9709  assets/model_quantized.onnx
07eced375cec144d27c900241f3e339478dec958f92fddbc551f295c992038a3  assets/vocab.txt
75287d17bc1b625afc686332f16aaeddaa129deb1b5fca2d128153fa4b0b17ab  assets/config.json
```

## Reproduce assets

```sh
./prepare_assets.sh
```

## When to use

Use this model when you want an English NER model with a much smaller embedded ONNX asset than BERT-base or DistilBERT.

Avoid it when you need maximum accuracy, multilingual support, or domain-specific labels. Its small size means it can be less accurate and may confuse some organization/location classes.
