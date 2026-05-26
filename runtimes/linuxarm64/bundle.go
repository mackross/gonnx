package linuxarm64

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
		LibraryRel: "assets/libonnxruntime.so",
	}
	b.ExtraRels = []string{"assets/libonnxruntime_providers_shared.so"}
	return bundles.Register(b)
}
