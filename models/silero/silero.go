// Package silero provides an embedded Silero VAD ONNX model bundle for gonnx.
package silero

import (
	"embed"

	"github.com/mackross/gonnx"
)

//go:embed models/silero_vad.onnx
var modelFS embed.FS

// Bundle returns the embedded Silero VAD model bundle.
func Bundle() gonnx.ModelBundle {
	return gonnx.ModelBundle{
		Name:     "silero-vad",
		FS:       modelFS,
		ModelRel: "models/silero_vad.onnx",
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
