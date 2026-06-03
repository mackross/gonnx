# bertuncased

Embedded English BERT NER model package for `gonnx/models`.

## Model

- Open: `bertuncased.Open(gonnx.WithThreads(1))`
- Model adapter: `bertuncased.BaseUncased()`
- Source/export: <https://huggingface.co/onnx-community/bert-base-NER-uncased-ONNX>
- Original model family: `dslim/bert-base-NER-uncased`
- License: MIT per upstream model card at packaging time
- Labels: `O`, `B-MISC`, `I-MISC`, `B-PER`, `I-PER`, `B-ORG`, `I-ORG`, `B-LOC`, `I-LOC`
- Tokenizer: BERT uncased WordPiece (`vocab.txt`), lower-casing enabled
- Max sequence length: 512 tokens including special tokens
- Embedded ONNX asset: quantized ONNX (`assets/model_quantized.onnx`)

## Checksums

```text
3036d3d0b80e30b2f99f4ca2bb34d93d3715430ad8b20d7b157146e388691bce  assets/model_quantized.onnx
07eced375cec144d27c900241f3e339478dec958f92fddbc551f295c992038a3  assets/vocab.txt
4f307e45513e88357bacb2ded56db364e3b01716f0fb8d41284c0a7363dba5a7  assets/config.json
```

Verify with:

```sh
cd models/bertuncased && sha256sum -c checksums.txt
```

## Reproduce assets

```sh
./prepare_assets.sh
```

The script downloads the ONNX model, vocabulary, and config from Hugging Face and verifies SHA256 checksums.

## When to use

Use this model for simple English CoNLL-style NER when a BERT-sized embedded model is acceptable and uncased behavior is fine.

Avoid it when you need a small binary, case-sensitive distinctions, custom labels, specialized domains, or non-English accuracy.
