package bertuncased

import (
	"github.com/mackross/gonnx"
	"github.com/mackross/gonnx/models/utils/nerutils"
	"github.com/mackross/gonnx/models/utils/nerutils/onnxner"
)

// Entity is a named entity span in the original input text.
type Entity = nerutils.Entity

// Recognizer runs tokenization, ONNX inference, and BIO chunking for this model.
type Recognizer = onnxner.Recognizer
type Model = onnxner.Model

// NewRecognizer opens the embedded model for inference.

// Open opens the embedded model for inference.
func Open(opts ...gonnx.Option) (*Recognizer, error) {
	return BaseUncased().Open(opts...)
}
