package linuxamd64

import (
	"embed"

	"github.com/mackross/gonnx/internal/bundles"
)

//go:embed assets/*
var bundleFS embed.FS

func register() error {
	return bundles.Register(bundles.Bundle{
		Platform:   platform,
		FS:         bundleFS,
		LibraryRel: "assets/libonnxruntime.so",
		ExtraRels:  []string{"assets/libonnxruntime_providers_shared.so"},
	})
}
