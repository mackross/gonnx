package gonnx

import (
	"runtime"

	"github.com/mackross/gonnx/internal/bundles"
	ort "github.com/shota3506/onnxruntime-purego/onnxruntime"
)

// NewRuntime loads the bundled ONNX Runtime for the current platform.
func NewRuntime(opts ...Option) (*ort.Runtime, error) {
	cfg := newConfig(opts)
	path, err := bundles.LibraryPath(runtime.GOOS + "/" + runtime.GOARCH)
	if err != nil {
		return nil, err
	}
	return ort.NewRuntime(path, cfg.apiVersion)
}
