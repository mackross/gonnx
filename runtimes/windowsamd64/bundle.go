package windowsamd64

import (
	"embed"

	"github.com/mackross/gonnx/internal/bundles"
)

//go:embed assets/*
var bundleFS embed.FS

func register() error {
	b := bundles.Bundle{
		Platform:   platform,
		FS:         bundleFS,
		LibraryRel: "assets/onnxruntime.dll",
		ExtraRels:  []string{"assets/onnxruntime_providers_shared.dll"},
	}
	return bundles.Register(b)
}
