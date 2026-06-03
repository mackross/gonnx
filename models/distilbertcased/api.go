package distilbertcased

import (
	"github.com/mackross/gonnx"
	"github.com/mackross/gonnx/models/utils/nerutils"
	"github.com/mackross/gonnx/models/utils/nerutils/onnxner"
)

type Entity = nerutils.Entity
type Recognizer = onnxner.Recognizer
type Model = onnxner.Model

// Open opens the embedded model for inference.
func Open(opts ...gonnx.Option) (*Recognizer, error) {
	return NER().Open(opts...)
}
