package onnxner

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"math"
	"os"
	"sort"
	"strconv"
	"sync"

	"github.com/mackross/gonnx"
	"github.com/mackross/gonnx/models/utils/nerutils"
	"github.com/mackross/gonnx/models/utils/nerutils/chunk"
	"github.com/mackross/gonnx/models/utils/nerutils/wordpiece"
	ort "github.com/shota3506/onnxruntime-purego/onnxruntime"
)

// Model is a BERT token-classification ONNX NER model.
type Model struct {
	Bundle  gonnx.ModelBundle
	Labels  []string
	VocabFS fs.FS
	Vocab   string
	Lower   bool
	MaxLen  int
}

func (m Model) Open(opts ...gonnx.Option) (*Recognizer, error) {
	if len(m.Labels) == 0 {
		return nil, fmt.Errorf("labels are required")
	}
	for i, l := range m.Labels {
		if l == "" {
			return nil, fmt.Errorf("label %d is empty", i)
		}
	}
	if m.VocabFS == nil || m.Vocab == "" {
		return nil, fmt.Errorf("vocabulary is required")
	}
	vf, err := m.VocabFS.Open(m.Vocab)
	if err != nil {
		return nil, fmt.Errorf("open vocabulary: %w", err)
	}
	tok, err := wordpiece.Load(vf, m.Lower)
	cerr := vf.Close()
	if err != nil {
		return nil, fmt.Errorf("load vocabulary: %w", err)
	}
	if cerr != nil {
		return nil, fmt.Errorf("close vocabulary: %w", cerr)
	}
	sess, err := gonnx.OpenBundle(m.Bundle, opts...)
	if err != nil {
		return nil, err
	}
	max := m.MaxLen
	if max == 0 {
		max = 512
	}
	return &Recognizer{session: sess, tokenizer: tok, labels: append([]string(nil), m.Labels...), maxLen: max}, nil
}

// LabelsFromConfig loads Hugging Face id2label labels from config.json.
func LabelsFromConfig(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg struct {
		ID2Label map[string]string `json:"id2label"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	labels := make([]string, len(cfg.ID2Label))
	for k, v := range cfg.ID2Label {
		i, err := strconv.Atoi(k)
		if err != nil {
			return nil, fmt.Errorf("invalid label id %q: %w", k, err)
		}
		if i < 0 || i >= len(labels) {
			return nil, fmt.Errorf("label id %d out of range", i)
		}
		labels[i] = v
	}
	return labels, nil
}

// Recognizer runs tokenization, ONNX inference, and BIO chunking.
type Recognizer struct {
	session   *gonnx.Session
	tokenizer *wordpiece.Tokenizer
	labels    []string
	maxLen    int
	mu        sync.Mutex
}

func (r *Recognizer) Recognize(ctx context.Context, text string) ([]nerutils.Entity, error) {
	tokens := r.tokenizer.Encode(text)
	if len(tokens) > r.maxLen {
		sep := tokens[len(tokens)-1]
		tokens = tokens[:r.maxLen]
		tokens[len(tokens)-1] = sep
	}
	ids, mask, types := make([]int64, len(tokens)), make([]int64, len(tokens)), make([]int64, len(tokens))
	for i, t := range tokens {
		ids[i], mask[i] = t.ID, 1
	}
	shape := []int64{1, int64(len(tokens))}
	inputIDs, err := gonnx.Tensor(r.session.Runtime, ids, shape...)
	if err != nil {
		return nil, fmt.Errorf("create input_ids tensor: %w", err)
	}
	defer inputIDs.Close()
	attentionMask, err := gonnx.Tensor(r.session.Runtime, mask, shape...)
	if err != nil {
		return nil, fmt.Errorf("create attention_mask tensor: %w", err)
	}
	defer attentionMask.Close()
	tokenTypeIDs, err := gonnx.Tensor(r.session.Runtime, types, shape...)
	if err != nil {
		return nil, fmt.Errorf("create token_type_ids tensor: %w", err)
	}
	defer tokenTypeIDs.Close()
	r.mu.Lock()
	outputs, err := r.session.Run(ctx, map[string]*ort.Value{"input_ids": inputIDs, "attention_mask": attentionMask, "token_type_ids": tokenTypeIDs})
	r.mu.Unlock()
	if err != nil {
		return nil, fmt.Errorf("run inference: %w", err)
	}
	defer closeOutputs(outputs)
	logitValue := outputs["logits"]
	if logitValue == nil {
		return nil, fmt.Errorf("model output missing logits")
	}
	logits, dims, err := gonnx.TensorData[float32](logitValue)
	if err != nil {
		return nil, fmt.Errorf("read logits: %w", err)
	}
	if len(dims) != 3 || dims[0] != 1 || dims[1] != int64(len(tokens)) || dims[2] != int64(len(r.labels)) {
		return nil, fmt.Errorf("unexpected logits shape %v for %d tokens and %d labels", dims, len(tokens), len(r.labels))
	}
	preds := make([]chunk.Token, 0, len(tokens))
	for i, tok := range tokens {
		start := i * len(r.labels)
		id, score := argmaxSoftmax(logits[start : start+len(r.labels)])
		preds = append(preds, chunk.Token{Text: tok.Text, Start: tok.Start, End: tok.End, Label: r.labels[id], Score: score, Special: tok.Special})
	}
	chunks, err := chunk.BIO(preds, text)
	if err != nil {
		return nil, err
	}
	out := make([]nerutils.Entity, len(chunks))
	for i, e := range chunks {
		out[i] = nerutils.Entity(e)
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Start < out[j].Start })
	return out, nil
}

func (r *Recognizer) Close() error {
	if r == nil || r.session == nil {
		return nil
	}
	r.session.Close()
	r.session = nil
	return nil
}

func closeOutputs(outputs map[string]*ort.Value) {
	for _, v := range outputs {
		if v != nil {
			v.Close()
		}
	}
}

func argmaxSoftmax(xs []float32) (int, float32) {
	best := 0
	for i := 1; i < len(xs); i++ {
		if xs[i] > xs[best] {
			best = i
		}
	}
	max := xs[best]
	var sum float64
	for _, x := range xs {
		sum += math.Exp(float64(x - max))
	}
	return best, float32(1 / sum)
}
