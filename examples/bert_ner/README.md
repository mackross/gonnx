# BERT NER example

This example loads an ONNX export of [`dslim/bert-base-NER-uncased`](https://huggingface.co/dslim/bert-base-NER-uncased) with `gonnx`.

The upstream Hugging Face repo publishes PyTorch weights. The live test and
example use the ONNX Community export:
[`onnx-community/bert-base-NER-uncased-ONNX`](https://huggingface.co/onnx-community/bert-base-NER-uncased-ONNX).

Download or place these files:

```text
examples/bert_ner/model.onnx
examples/bert_ner/vocab.txt
examples/bert_ner/config.json
```

Then run with an installed ONNX Runtime shared library:

```bash
go run ./examples/bert_ner \
  -runtime /path/to/libonnxruntime.so \
  -model ./examples/bert_ner/model.onnx
```

You can also set `ONNXRUNTIME_LIBRARY` instead of `-runtime`.

For a bundled app, import one runtime package and run:

```bash
go run ./examples/bert_ner -model ./examples/bert_ner/model.onnx
```

Live test:

```bash
GONNX_LIVE_BERT_NER=1 go test -v ./examples/bert_ner -run TestLiveBertNERRecognizesEntities
```
