// Package darwinarm64 registers the embedded ONNX Runtime bundle for darwin/arm64.
package darwinarm64

const platform = "darwin/arm64"

func init() {
	if err := register(); err != nil {
		panic(err)
	}
}
