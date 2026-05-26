// Package windowsarm64 registers the embedded ONNX Runtime bundle for windows/arm64.
package windowsarm64

const platform = "windows/arm64"

func init() {
	if err := register(); err != nil {
		panic(err)
	}
}
