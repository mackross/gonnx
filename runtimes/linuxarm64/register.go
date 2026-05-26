// Package linuxarm64 registers the embedded ONNX Runtime bundle for linux/arm64.
package linuxarm64

const platform = "linux/arm64"

func init() {
	if err := register(); err != nil {
		panic(err)
	}
}
