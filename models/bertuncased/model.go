// Package bertuncased provides a ready-to-use embedded BERT NER model.
//
// BaseUncased uses the ONNX Community export of dslim/bert-base-NER-uncased:
// https://huggingface.co/onnx-community/bert-base-NER-uncased-ONNX
//
// Use this model for English, uncased, general-purpose person/organization/
// location/miscellaneous entity recognition when a BERT-sized model is
// acceptable. Do not use it for case-sensitive text distinctions, non-English
// workloads, domains that require custom labels, or small binaries.
//
// Source model: dslim/bert-base-NER-uncased via ONNX Community export.
// Upstream URL: https://huggingface.co/onnx-community/bert-base-NER-uncased-ONNX
// License: MIT, according to the upstream model card at time of packaging.
// Labels: O, B-MISC, I-MISC, B-PER, I-PER, B-ORG, I-ORG, B-LOC, I-LOC.
// Tokenizer: BERT uncased WordPiece vocabulary, lower-casing enabled.
// Max sequence length: 512 WordPiece tokens including [CLS] and [SEP]; longer
// inputs are truncated by the current adapter.
// Known limitations: English-oriented CoNLL-style labels only; quantized ONNX
// weights may trade a small amount of accuracy for size and speed.
//
// Embedded asset SHA256 checksums:
//   - assets/model_quantized.onnx: 3036d3d0b80e30b2f99f4ca2bb34d93d3715430ad8b20d7b157146e388691bce
//   - assets/vocab.txt: 07eced375cec144d27c900241f3e339478dec958f92fddbc551f295c992038a3
//   - assets/config.json: 4f307e45513e88357bacb2ded56db364e3b01716f0fb8d41284c0a7363dba5a7
package bertuncased

import (
	"embed"
	"io/fs"

	"github.com/mackross/gonnx"
)

//go:embed assets/model_quantized.onnx assets/vocab.txt assets/config.json
var assets embed.FS

var labels = []string{"O", "B-MISC", "I-MISC", "B-PER", "I-PER", "B-ORG", "I-ORG", "B-LOC", "I-LOC"}

// FromBundle returns a dslim/bert-base-NER-uncased compatible model adapter using caller-provided assets.
func FromBundle(bundle gonnx.ModelBundle, vocabFS fs.FS, vocabRel string) Model {
	return Model{Bundle: bundle, Labels: labels, VocabFS: vocabFS, Vocab: vocabRel, Lower: true, MaxLen: 512}
}

// BaseUncased returns an embedded dslim/bert-base-NER-uncased ONNX model.
func BaseUncased() Model {
	return Model{
		Bundle: gonnx.ModelBundle{
			Name:     "bertuncased-base-uncased-quantized",
			FS:       assets,
			ModelRel: "assets/model_quantized.onnx",
			ExtraRels: []string{
				"assets/config.json",
				"assets/vocab.txt",
			},
		},
		Labels:  labels,
		VocabFS: assets,
		Vocab:   "assets/vocab.txt",
		Lower:   true,
		MaxLen:  512,
	}
}
