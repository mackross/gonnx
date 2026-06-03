package onnxner

import "testing"

func TestArgmaxSoftmax(t *testing.T) {
	id, score := argmaxSoftmax([]float32{1, 3, 2})
	if id != 1 {
		t.Fatalf("id=%d", id)
	}
	if score < 0.66 || score > 0.67 {
		t.Fatalf("score=%g", score)
	}
}

func TestLabelsFromConfigErrors(t *testing.T) {
	if _, err := LabelsFromConfig("testdata/missing.json"); err == nil {
		t.Fatal("expected missing file error")
	}
}

func TestLabelsFromConfig(t *testing.T) {
	got, err := LabelsFromConfig("testdata/config.json")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"O", "B-PER", "I-PER"}
	if len(got) != len(want) {
		t.Fatalf("got %#v want %#v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %#v want %#v", got, want)
		}
	}
}
