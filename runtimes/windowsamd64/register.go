// Package windowsamd64 registers the embedded ONNX Runtime bundle for windows/amd64.
package windowsamd64

const platform = "windows/amd64"

func init() {
	if err := register(); err != nil {
		panic(err)
	}
}
