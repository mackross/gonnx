// Package gonnx opens ONNX models using bundled ONNX Runtime binaries from pure Go.
//
// Import the runtime package you want to ship, usually as a blank import, then
// call Open:
//
//	import _ "github.com/mackross/gonnx/runtimes/linuxamd64"
//
//	sess, err := gonnx.Open("model.onnx")
//
// Runtime packages under runtimes/* embed native ONNX Runtime binaries. The
// root package extracts the selected runtime to a checked temp cache and loads
// it through onnxruntime-purego. Use OpenBundle for embedded model files and
// PrepareModelBundle only when you want to prewarm that extraction cache.
package gonnx
