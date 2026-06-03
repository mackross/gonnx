// Package neurobert provides a tiny embedded English NER model.
//
// NER uses the ONNX Community export of boltuix/NeuroBERT-NER:
// https://huggingface.co/onnx-community/NeuroBERT-NER-ONNX
//
// Use this model when binary size and CPU latency matter more than maximum
// accuracy. It recognizes a richer OntoNotes-like English label set than the
// CoNLL-focused dslim models. Do not use it when you need the strongest general
// English accuracy, multilingual support, or custom/domain labels.
//
// Source model: boltuix/NeuroBERT-NER via ONNX Community export.
// Upstream ONNX revision packaged: e3a6a290e438a506d2ba7bbb0f5875b7f034cebf.
// Source model revision observed: fdc13cd31d4abc10b1b54df348caf10727a1ad93.
// Upstream URL: https://huggingface.co/onnx-community/NeuroBERT-NER-ONNX
// Source URL: https://huggingface.co/boltuix/NeuroBERT-NER
// License: Apache-2.0 according to the source model card at packaging time.
// Labels: O plus BIO labels for CARDINAL, DATE, EVENT, FAC, GPE, LANGUAGE, LAW,
// LOC, MONEY, NORP, ORDINAL, ORG, PERCENT, PERSON, PRODUCT, QUANTITY, TIME,
// and WORK_OF_ART.
// Tokenizer: BERT uncased WordPiece vocabulary, lower-casing enabled.
// Max sequence length: 512 WordPiece tokens including [CLS] and [SEP]; longer
// inputs are truncated by the current adapter.
// Known limitations: English-oriented labels only; very small model may trade
// accuracy for size and speed.
//
// Embedded asset SHA256 checksums:
//   - assets/model_quantized.onnx: 17dbf6ccfe500ee8a5da9177ebcdbb1c7394626e077e70571adebae8833c9709
//   - assets/vocab.txt: 07eced375cec144d27c900241f3e339478dec958f92fddbc551f295c992038a3
//   - assets/config.json: 75287d17bc1b625afc686332f16aaeddaa129deb1b5fca2d128153fa4b0b17ab
package neurobert

import (
	"embed"

	"github.com/mackross/gonnx"
)

//go:embed assets/model_quantized.onnx assets/vocab.txt assets/config.json
var assets embed.FS

var labels = []string{
	"B-CARDINAL", "B-DATE", "B-EVENT", "B-FAC", "B-GPE", "B-LANGUAGE", "B-LAW", "B-LOC", "B-MONEY", "B-NORP", "B-ORDINAL", "B-ORG", "B-PERCENT", "B-PERSON", "B-PRODUCT", "B-QUANTITY", "B-TIME", "B-WORK_OF_ART",
	"I-CARDINAL", "I-DATE", "I-EVENT", "I-FAC", "I-GPE", "I-LANGUAGE", "I-LAW", "I-LOC", "I-MONEY", "I-NORP", "I-ORDINAL", "I-ORG", "I-PERCENT", "I-PERSON", "I-PRODUCT", "I-QUANTITY", "I-TIME", "I-WORK_OF_ART",
	"O",
}

// NER returns an embedded boltuix/NeuroBERT-NER quantized ONNX model.
func NER() Model {
	return Model{Bundle: gonnx.ModelBundle{Name: "neurobert-ner-quantized", FS: assets, ModelRel: "assets/model_quantized.onnx", ExtraRels: []string{"assets/config.json", "assets/vocab.txt"}}, Labels: labels, VocabFS: assets, Vocab: "assets/vocab.txt", Lower: true, MaxLen: 512}
}
