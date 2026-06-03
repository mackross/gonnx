package neurobert_test

import (
	"context"
	"testing"

	"github.com/mackross/gonnx"
	"github.com/mackross/gonnx/models/neurobert"
	_ "github.com/mackross/gonnx/runtimes/linuxamd64"
)

func TestModelInference(t *testing.T) {
	recognizer, err := neurobert.Open(gonnx.WithThreads(1))
	if err != nil {
		t.Fatal(err)
	}
	defer recognizer.Close()

	entities, err := recognizer.Recognize(context.Background(), "Barack Obama worked at Microsoft in Seattle.")
	if err != nil {
		t.Fatal(err)
	}
	for _, entity := range entities {
		if entity.Text == "" || entity.Label == "" || entity.Start < 0 || entity.End <= entity.Start {
			t.Fatalf("invalid entity: %#v", entity)
		}
	}
}

func BenchmarkModelInference(b *testing.B) {
	ctx := context.Background()
	const wordPieceTokens = 10
	recognizer, err := neurobert.Open(gonnx.WithThreads(1))
	if err != nil {
		b.Fatal(err)
	}
	defer recognizer.Close()
	text := "Barack Obama worked at Microsoft in Seattle."
	if _, err := recognizer.Recognize(ctx, text); err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := recognizer.Recognize(ctx, text); err != nil {
			b.Fatal(err)
		}
	}
	b.ReportMetric(float64(b.N*wordPieceTokens)/b.Elapsed().Seconds(), "tokens/s")
}
