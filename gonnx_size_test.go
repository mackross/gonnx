package gonnx_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestSelectedRuntimePackagesAffectBuildSize(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping build-size integration test in short mode")
	}

	repoDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "go.mod"), fmt.Sprintf(`module gonnx-size-fixture

go 1.24.0

require github.com/mackross/gonnx v0.0.0

replace github.com/mackross/gonnx => %s
`, filepath.ToSlash(repoDir)))

	writeFile(t, filepath.Join(dir, "cmd", "one", "main.go"), `package main

import (
	"fmt"

	"github.com/mackross/gonnx"
	_ "github.com/mackross/gonnx/runtimes/linuxamd64"
)

func main() {
	rt, err := gonnx.NewRuntime()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rt.Close()
	fmt.Println("ok")
}
`)
	writeFile(t, filepath.Join(dir, "cmd", "two", "main.go"), `package main

import (
	"fmt"

	"github.com/mackross/gonnx"
	_ "github.com/mackross/gonnx/runtimes/linuxamd64"
	_ "github.com/mackross/gonnx/runtimes/linuxarm64"
)

func main() {
	rt, err := gonnx.NewRuntime()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rt.Close()
	fmt.Println("ok")
}
`)

	one := buildFixture(t, dir, "./cmd/one")
	two := buildFixture(t, dir, "./cmd/two")
	oneInfo, err := os.Stat(one)
	if err != nil {
		t.Fatal(err)
	}
	twoInfo, err := os.Stat(two)
	if err != nil {
		t.Fatal(err)
	}

	delta := twoInfo.Size() - oneInfo.Size()
	if delta < 15*1024*1024 {
		t.Fatalf("expected importing a second runtime package to add roughly its embedded asset size; one=%d two=%d delta=%d", oneInfo.Size(), twoInfo.Size(), delta)
	}
}

func buildFixture(t *testing.T, dir, pkg string) string {
	t.Helper()
	out := filepath.Join(dir, "bin", filepath.Base(pkg))
	if runtime.GOOS == "windows" {
		out += ".exe"
	}
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("go", "build", "-trimpath", "-mod=mod", "-o", out, pkg)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go build %s failed: %v\n%s", pkg, err, output)
	}
	return out
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
