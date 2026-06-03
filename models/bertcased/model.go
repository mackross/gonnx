// Package bertcased provides an embedded cased BERT-base English NER model.
package bertcased

import (
	"embed"

	"github.com/mackross/gonnx"
)

//go:embed assets/model_quantized.onnx assets/vocab.txt assets/config.json
var assets embed.FS

var labels = []string{"O", "B-MISC", "I-MISC", "B-PER", "I-PER", "B-ORG", "I-ORG", "B-LOC", "I-LOC"}

// BaseCased returns an embedded dslim/bert-base-NER cased quantized ONNX model.
func BaseCased() Model {
	return Model{Bundle: gonnx.ModelBundle{Name: "bert-base-ner-cased-quantized", FS: assets, ModelRel: "assets/model_quantized.onnx", ExtraRels: []string{"assets/config.json", "assets/vocab.txt"}}, Labels: labels, VocabFS: assets, Vocab: "assets/vocab.txt", Lower: false, MaxLen: 512}
}
