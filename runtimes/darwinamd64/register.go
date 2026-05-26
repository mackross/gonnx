// Package darwinamd64 registers the embedded ONNX Runtime bundle for darwin/amd64.
package darwinamd64

const platform = "darwin/amd64"

func init() {
	if err := register(); err != nil {
		panic(err)
	}
}
