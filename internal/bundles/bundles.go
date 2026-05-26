package bundles

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

// Bundle describes a platform-specific ONNX Runtime shared-library bundle.
type Bundle struct {
	// Platform is GOOS/GOARCH, for example "linux/amd64".
	Platform string
	// FS contains LibraryRel and ExtraRels.
	FS fs.FS
	// LibraryRel is the path within FS to the loadable ORT library.
	LibraryRel string
	// ExtraRels are sibling runtime files that must be present next to LibraryRel.
	ExtraRels []string
}

var (
	mu       sync.RWMutex
	registry = map[string]Bundle{}
)

// Register registers a bundled runtime. Platform defaults to runtime.GOOS+"/"+runtime.GOARCH.
func Register(bundle Bundle) error {
	if bundle.Platform == "" {
		bundle.Platform = runtime.GOOS + "/" + runtime.GOARCH
	}
	if bundle.FS == nil || bundle.LibraryRel == "" {
		return fmt.Errorf("gonnx: runtime bundle requires FS and LibraryRel")
	}
	for _, rel := range append([]string{bundle.LibraryRel}, bundle.ExtraRels...) {
		if _, err := fs.Stat(bundle.FS, rel); err != nil {
			return fmt.Errorf("gonnx: runtime bundle %s missing %s: %w", bundle.Platform, rel, err)
		}
	}
	mu.Lock()
	defer mu.Unlock()
	registry[bundle.Platform] = bundle
	return nil
}

// LibraryPath extracts the registered runtime for platform and returns its library path.
func LibraryPath(platform string) (string, error) {
	mu.RLock()
	bundle, ok := registry[platform]
	mu.RUnlock()
	if !ok {
		return "", fmt.Errorf("gonnx: no ONNX Runtime bundle registered for %s", platform)
	}
	return Extract(bundle)
}

// Extract extracts bundle files to a content-address-checked temp directory.
func Extract(bundle Bundle) (string, error) {
	dir := filepath.Join(os.TempDir(), "gonnx", sanitizePlatform(bundle.Platform))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	for _, rel := range append([]string{bundle.LibraryRel}, bundle.ExtraRels...) {
		if err := extractFile(dir, bundle.FS, rel); err != nil {
			return "", err
		}
	}
	return filepath.Join(dir, filepath.Base(bundle.LibraryRel)), nil
}

func extractFile(dir string, fsys fs.FS, rel string) error {
	data, err := fs.ReadFile(fsys, rel)
	if err != nil {
		return err
	}
	dst := filepath.Join(dir, filepath.Base(rel))
	sum := sha256.Sum256(data)
	sumHex := hex.EncodeToString(sum[:])
	if got, err := os.ReadFile(dst + ".sha256"); err == nil && string(got) == sumHex {
		if _, err := os.Stat(dst); err == nil {
			return nil
		}
	}
	if err := os.WriteFile(dst, data, 0o755); err != nil {
		return err
	}
	return os.WriteFile(dst+".sha256", []byte(sumHex), 0o644)
}

func sanitizePlatform(platform string) string {
	out := []byte(platform)
	for i, c := range out {
		if c == '/' || c == '\\' || c == ':' {
			out[i] = '-'
		}
	}
	return string(out)
}
