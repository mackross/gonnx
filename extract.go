package gonnx

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
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

func extractJoinedFile(dir string, fsys fs.FS, rel string, parts []string) error {
	sumHex, err := joinedSHA256(fsys, parts)
	if err != nil {
		return err
	}
	dst := filepath.Join(dir, filepath.Base(rel))
	if got, err := os.ReadFile(dst + ".sha256"); err == nil && string(got) == sumHex {
		if _, err := os.Stat(dst); err == nil {
			return nil
		}
	}
	tmp := dst + ".tmp"
	out, err := os.OpenFile(tmp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o755)
	if err != nil {
		return err
	}
	for _, part := range parts {
		if err := copyPart(out, fsys, part); err != nil {
			_ = out.Close()
			_ = os.Remove(tmp)
			return err
		}
	}
	if err := out.Close(); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	if err := os.Rename(tmp, dst); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return os.WriteFile(dst+".sha256", []byte(sumHex), 0o644)
}

func joinedSHA256(fsys fs.FS, parts []string) (string, error) {
	h := sha256.New()
	for _, part := range parts {
		if err := copyPart(h, fsys, part); err != nil {
			return "", err
		}
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func copyPart(dst io.Writer, fsys fs.FS, rel string) error {
	f, err := fsys.Open(rel)
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(dst, f)
	closeErr := f.Close()
	if copyErr != nil {
		return copyErr
	}
	return closeErr
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
