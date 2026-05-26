// Package linuxamd64 registers the embedded ONNX Runtime bundle for linux/amd64.
package linuxamd64

const platform = "linux/amd64"

func init() {
	if err := register(); err != nil {
		panic(err)
	}
}
