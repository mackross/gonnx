# distilbert

Embedded English DistilBERT NER model package for `gonnx/models`.

## Model

- Open: `distilbertcased.Open(gonnx.WithThreads(1))`
- Model adapter: `distilbertcased.NER()`
- Source/export: <https://huggingface.co/onnx-community/distilbert-NER-ONNX>
- Upstream revision packaged: `3a19fe9404a4469d91aa3d551558a97f68872f67`
- Original model family: `dslim/distilbert-NER`
- License: MIT per upstream model card at packaging time
- Labels: `O`, `B-PER`, `I-PER`, `B-ORG`, `I-ORG`, `B-LOC`, `I-LOC`, `B-MISC`, `I-MISC`
- Tokenizer: DistilBERT cased WordPiece (`vocab.txt`), lower-casing disabled
- Max sequence length: 512 tokens including special tokens
- Embedded ONNX asset: quantized ONNX (`assets/model_quantized.onnx`)

## Checksums

```text
9419a876387ff2bbe5f21ab7429c7bef93eac86c50353390d4d8fca6e4a210d8  assets/model_quantized.onnx
eeaa9875b23b04b4c54ef759d03db9d1ba1554838f8fb26c5d96fa551df93d02  assets/vocab.txt
f109facddb205dac712adf5877e4315fae62041bc0916fe808d92abdb594d1fe  assets/config.json
```

## Reproduce assets

```sh
./prepare_assets.sh
```

## When to use

Use this model when you want an English cased NER model that is smaller and faster than the BERT-base package.

Avoid it when you need uncased normalization, custom labels, specialized domains, or non-English accuracy.
