package darwinarm64

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
		LibraryRel: "assets/libonnxruntime.dylib",
	}
	return bundles.Register(b)
}
