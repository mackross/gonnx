// Package distilbert provides a ready-to-use embedded DistilBERT NER model.
//
// NER uses the ONNX Community export of dslim/distilbert-NER:
// https://huggingface.co/onnx-community/distilbert-NER-ONNX
//
// Use this model for English, cased, general-purpose person/organization/
// location/miscellaneous entity recognition when lower runtime and binary cost
// than BERT-base is preferred. Do not use it for uncased text normalization,
// non-English workloads, domains that require custom labels, or maximum BERT
// accuracy.
//
// Source model: dslim/distilbert-NER via ONNX Community export.
// Upstream revision packaged: 3a19fe9404a4469d91aa3d551558a97f68872f67.
// Upstream URL: https://huggingface.co/onnx-community/distilbert-NER-ONNX
// License: MIT, according to the upstream model card at time of packaging.
// Labels: O, B-PER, I-PER, B-ORG, I-ORG, B-LOC, I-LOC, B-MISC, I-MISC.
// Tokenizer: DistilBERT cased WordPiece vocabulary, lower-casing disabled.
// Max sequence length: 512 WordPiece tokens including special tokens; longer
// inputs are truncated by the current adapter.
// Known limitations: English-oriented CoNLL-style labels only; quantized ONNX
// weights may trade a small amount of accuracy for size and speed.
//
// Embedded asset SHA256 checksums:
//   - assets/model_quantized.onnx: 9419a876387ff2bbe5f21ab7429c7bef93eac86c50353390d4d8fca6e4a210d8
//   - assets/vocab.txt: eeaa9875b23b04b4c54ef759d03db9d1ba1554838f8fb26c5d96fa551df93d02
//   - assets/config.json: f109facddb205dac712adf5877e4315fae62041bc0916fe808d92abdb594d1fe
package distilbertcased

import (
	"embed"

	"github.com/mackross/gonnx"
)

//go:embed assets/model_quantized.onnx assets/vocab.txt assets/config.json
var assets embed.FS

var labels = []string{"O", "B-PER", "I-PER", "B-ORG", "I-ORG", "B-LOC", "I-LOC", "B-MISC", "I-MISC"}

// NER returns an embedded dslim/distilbert-NER quantized ONNX model.
func NER() Model {
	return Model{Bundle: gonnx.ModelBundle{Name: "distilbert-ner-quantized", FS: assets, ModelRel: "assets/model_quantized.onnx", ExtraRels: []string{"assets/config.json", "assets/vocab.txt"}}, Labels: labels, VocabFS: assets, Vocab: "assets/vocab.txt", Lower: false, MaxLen: 512}
}
