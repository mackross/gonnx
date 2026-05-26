package gonnx

import (
	"crypto/sha256"
	"encoding/hex"
	"io/fs"
	"os"
	"path/filepath"
)

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
