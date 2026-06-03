// Package multidistilbert provides an embedded multilingual DistilBERT NER model.
//
// NER uses the Xenova ONNX export of Davlan/distilbert-base-multilingual-cased-ner-hrl:
// https://huggingface.co/Xenova/distilbert-base-multilingual-cased-ner-hrl
//
// Use this model for multilingual cased NER over high-resource languages when
// the generic labels DATE, PER, ORG, and LOC are sufficient. Do not use it when
// you need English-only maximum accuracy, PII-specific labels, domain labels, or
// a small binary.
//
// Source model: Davlan/distilbert-base-multilingual-cased-ner-hrl via Xenova ONNX export.
// Upstream ONNX revision packaged: c2a4dbf593c57f47004c5bc2d3770d311aee9c43.
// Source model revision observed: d421f57d5b1d36b375408588669e9340f9b11a89.
// Upstream URL: https://huggingface.co/Xenova/distilbert-base-multilingual-cased-ner-hrl
// Source URL: https://huggingface.co/Davlan/distilbert-base-multilingual-cased-ner-hrl
// License: AFL-3.0 according to the source model card at packaging time.
// Labels: O, B-DATE, I-DATE, B-PER, I-PER, B-ORG, I-ORG, B-LOC, I-LOC.
// Tokenizer: multilingual cased WordPiece vocabulary, lower-casing disabled.
// Max sequence length: 512 WordPiece tokens including [CLS] and [SEP]; longer
// inputs are truncated by the current adapter.
// Known limitations: high-resource-language training scope from upstream model;
// quantized ONNX weights may trade a small amount of accuracy for size/speed.
//
// Embedded asset SHA256 checksums:
//   - assets/model_quantized.onnx: 24a0b98f4dd4cd92842f5a541272f86f760225a64a29928eddef14bdb2edb986
//   - assets/vocab.txt: fe0fda7c425b48c516fc8f160d594c8022a0808447475c1a7c6d6479763f310c
//   - assets/config.json: 38847be4dc6699b1218a749ed69f888c2ccc7b4deba98e3c4a1cac8cb34d54c8
package multidistilbert

import (
	"embed"

	"github.com/mackross/gonnx"
)

//go:embed assets/model_quantized.onnx assets/vocab.txt assets/config.json
var assets embed.FS

var labels = []string{"O", "B-DATE", "I-DATE", "B-PER", "I-PER", "B-ORG", "I-ORG", "B-LOC", "I-LOC"}

// NER returns an embedded multilingual DistilBERT quantized ONNX NER model.
func NER() Model {
	return Model{Bundle: gonnx.ModelBundle{Name: "multidistilbert-ner-quantized", FS: assets, ModelRel: "assets/model_quantized.onnx", ExtraRels: []string{"assets/config.json", "assets/vocab.txt"}}, Labels: labels, VocabFS: assets, Vocab: "assets/vocab.txt", Lower: false, MaxLen: 512}
}
