package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mackross/gonnx"
	ort "github.com/shota3506/onnxruntime-purego/onnxruntime"
)

func BenchmarkLiveBertNERExtraction(b *testing.B) {
	if os.Getenv("GONNX_LIVE_BERT_NER") == "" {
		b.Skip("set GONNX_LIVE_BERT_NER=1 to download and benchmark the live BERT NER model")
	}

	cacheDir := filepath.Join(os.TempDir(), "gonnx-bert-ner-live")
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		b.Fatal(err)
	}

	modelPath := filepath.Join(cacheDir, "model_quantized.onnx")
	vocabPath := filepath.Join(cacheDir, "vocab.txt")
	configPath := filepath.Join(cacheDir, "config.json")
	mustDownload(b, "https://huggingface.co/onnx-community/bert-base-NER-uncased-ONNX/resolve/main/onnx/model_quantized.onnx", modelPath)
	mustDownload(b, "https://huggingface.co/onnx-community/bert-base-NER-uncased-ONNX/resolve/main/vocab.txt", vocabPath)
	mustDownload(b, "https://huggingface.co/onnx-community/bert-base-NER-uncased-ONNX/resolve/main/config.json", configPath)

	labels, err := loadLabels(configPath)
	if err != nil {
		b.Fatal(err)
	}
	tok, err := loadTokenizer(vocabPath)
	if err != nil {
		b.Fatal(err)
	}

	model, err := os.Open(modelPath)
	if err != nil {
		b.Fatal(err)
	}
	defer model.Close()

	sess, err := gonnx.OpenReader(model, gonnx.WithLogLevel(ort.LoggingLevelWarning), gonnx.WithThreads(1))
	if err != nil {
		b.Fatal(err)
	}
	defer sess.Close()

	seed := strings.Join([]string{
		"Hugging Face Inc. is based in New York City and Sarah lives in London.",
		"Microsoft and OpenAI announced new research from Seattle while Alice visited Berlin.",
		"Google Cloud works with Deutsche Bank in Frankfurt and Carlos moved to Madrid.",
		"The United Nations met in Geneva before Emma flew from Paris to Tokyo.",
	}, " ")

	for _, targetTokens := range []int{64, 128, 256} {
		text := textForTokens(tok, seed, targetTokens)
		ids, _, _, _ := tok.Encode(text)
		b.Run(fmt.Sprintf("tokens_%d", len(ids)), func(b *testing.B) {
			ctx := context.Background()
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				entities, err := recognize(ctx, sess, tok, labels, text)
				if err != nil {
					b.Fatal(err)
				}
				if len(entities) == 0 {
					b.Fatal("no entities returned")
				}
			}
		})
	}
}

func textForTokens(tok *wordPieceTokenizer, seed string, target int) string {
	parts := []string{seed}
	for {
		text := strings.Join(parts, " ")
		ids, _, _, _ := tok.Encode(text)
		if len(ids) >= target {
			return text
		}
		parts = append(parts, seed)
	}
}
