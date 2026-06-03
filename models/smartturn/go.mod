module github.com/mackross/gonnx/models/smartturn

go 1.25

require (
	github.com/mackross/gonnx v0.0.0
	github.com/mackross/gonnx/models v0.0.0
	github.com/shota3506/onnxruntime-purego v0.0.0-20260315223538-8db8bd7424b2
	gonum.org/v1/gonum v0.17.0
)

require github.com/ebitengine/purego v0.9.0 // indirect

replace github.com/mackross/gonnx => ../..

replace github.com/mackross/gonnx/models => ..

replace github.com/shota3506/onnxruntime-purego => ../../third_party/onnxruntime-purego
