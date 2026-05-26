package windowsamd64

import "testing"

func TestBundledRuntimeAssetManifest(t *testing.T) {
	if err := register(); err != nil {
		t.Fatal(err)
	}
}
