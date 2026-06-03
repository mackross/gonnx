// Package smartturn provides an embedded Smart Turn ONNX model bundle for gonnx.
package smartturn

import (
	"embed"

	"github.com/mackross/gonnx"
)

//go:embed models/smart-turn-v3.2-cpu.onnx
var modelFS embed.FS

// Bundle returns the embedded Smart Turn model bundle.
func Bundle() gonnx.ModelBundle {
	return gonnx.ModelBundle{
		Name:     "smartturn-v3.2-cpu",
		FS:       modelFS,
		ModelRel: "models/smart-turn-v3.2-cpu.onnx",
	}
}

// Prepare extracts the embedded model to gonnx's checked temp cache and returns
// the extracted ONNX model path.
func Prepare() (string, error) {
	return gonnx.PrepareModelBundle(Bundle())
}

// OpenSession extracts and opens the embedded model as a raw gonnx session.
func OpenSession(opts ...gonnx.Option) (*gonnx.Session, error) {
	return gonnx.OpenBundle(Bundle(), opts...)
}
