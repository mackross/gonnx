package main

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mackross/gonnx"
	ort "github.com/shota3506/onnxruntime-purego/onnxruntime"
)

func TestLiveBertNERRecognizesEntities(t *testing.T) {
	if os.Getenv("GONNX_LIVE_BERT_NER") == "" {
		t.Skip("set GONNX_LIVE_BERT_NER=1 to download and run the live BERT NER model")
	}

	cacheDir := filepath.Join(os.TempDir(), "gonnx-bert-ner-live")
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		t.Fatal(err)
	}

	modelPath := filepath.Join(cacheDir, "model_quantized.onnx")
	vocabPath := filepath.Join(cacheDir, "vocab.txt")
	configPath := filepath.Join(cacheDir, "config.json")
	mustDownload(t, "https://huggingface.co/onnx-community/bert-base-NER-uncased-ONNX/resolve/main/onnx/model_quantized.onnx", modelPath)
	mustDownload(t, "https://huggingface.co/onnx-community/bert-base-NER-uncased-ONNX/resolve/main/vocab.txt", vocabPath)
	mustDownload(t, "https://huggingface.co/onnx-community/bert-base-NER-uncased-ONNX/resolve/main/config.json", configPath)

	labels, err := loadLabels(configPath)
	if err != nil {
		t.Fatal(err)
	}
	tok, err := loadTokenizer(vocabPath)
	if err != nil {
		t.Fatal(err)
	}

	model, err := os.Open(modelPath)
	if err != nil {
		t.Fatal(err)
	}
	defer model.Close()

	sess, err := gonnx.OpenReader(model, gonnx.WithLogLevel(ort.LoggingLevelWarning))
	if err != nil {
		t.Fatal(err)
	}
	defer sess.Close()

	entities, err := recognize(context.Background(), sess, tok, labels, "Hugging Face Inc. is based in New York City and Sarah lives in London.")
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]string{}
	for _, e := range entities {
		t.Logf("%s %s %.3f", e.Token, e.Label, e.Score)
		got[e.Token] = e.Label
	}
	for token, prefix := range map[string]string{"hugging": "B-ORG", "york": "I-LOC", "sarah": "B-PER", "london": "B-LOC"} {
		if !strings.HasPrefix(got[token], prefix) {
			t.Fatalf("entity %q = %q, want %q; all entities: %+v", token, got[token], prefix, entities)
		}
	}
}

func mustDownload(t testing.TB, url, path string) {
	t.Helper()
	if info, err := os.Stat(path); err == nil && info.Size() > 0 {
		return
	}
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET %s: %s", url, resp.Status)
	}
	tmp := path + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.Copy(f, resp.Body); err != nil {
		f.Close()
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(tmp, path); err != nil {
		t.Fatal(err)
	}
}
