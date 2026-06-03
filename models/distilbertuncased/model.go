// Package distilbertuncased provides an embedded uncased DistilBERT English NER model.
//
// License: Apache-2.0 according to the source model card at packaging time.
package distilbertuncased

import (
	"embed"

	"github.com/mackross/gonnx"
)

//go:embed assets/model_quantized.onnx assets/vocab.txt assets/config.json
var assets embed.FS

var labels = []string{"O", "B-PER", "I-PER", "B-ORG", "I-ORG", "B-LOC", "I-LOC", "B-MISC", "I-MISC"}

// NER returns an embedded uncased DistilBERT CoNLL03 quantized ONNX NER model.
func NER() Model {
	return Model{Bundle: gonnx.ModelBundle{Name: "distilbert-ner-uncased-quantized", FS: assets, ModelRel: "assets/model_quantized.onnx", ExtraRels: []string{"assets/config.json", "assets/vocab.txt"}}, Labels: labels, VocabFS: assets, Vocab: "assets/vocab.txt", Lower: true, MaxLen: 512}
}
