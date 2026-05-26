package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"unicode"

	"github.com/mackross/gonnx"
	ort "github.com/shota3506/onnxruntime-purego/onnxruntime"
)

type Entity struct {
	Token string
	Label string
	Score float32
}

type bertConfig struct {
	ID2Label map[string]string `json:"id2label"`
}

type wordPieceTokenizer struct {
	vocab map[string]int64
}

func loadLabels(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg bertConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	labels := make([]string, len(cfg.ID2Label))
	for id, label := range cfg.ID2Label {
		var idx int
		if _, err := fmt.Sscanf(id, "%d", &idx); err != nil {
			return nil, err
		}
		labels[idx] = label
	}
	return labels, nil
}

func loadTokenizer(vocabPath string) (*wordPieceTokenizer, error) {
	f, err := os.Open(vocabPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	vocab := map[string]int64{}
	s := bufio.NewScanner(f)
	var id int64
	for s.Scan() {
		vocab[s.Text()] = id
		id++
	}
	return &wordPieceTokenizer{vocab: vocab}, s.Err()
}

func (t *wordPieceTokenizer) Encode(text string) (ids, mask, types []int64, tokens []string) {
	tokens = append(tokens, "[CLS]")
	ids = append(ids, t.vocab["[CLS]"])
	for _, word := range basicTokens(text) {
		pieces := t.wordPieces(word)
		for _, p := range pieces {
			tokens = append(tokens, p)
			ids = append(ids, t.vocab[p])
		}
	}
	tokens = append(tokens, "[SEP]")
	ids = append(ids, t.vocab["[SEP]"])
	mask = make([]int64, len(ids))
	types = make([]int64, len(ids))
	for i := range mask {
		mask[i] = 1
	}
	return ids, mask, types, tokens
}

func (t *wordPieceTokenizer) wordPieces(word string) []string {
	word = strings.ToLower(word)
	if _, ok := t.vocab[word]; ok {
		return []string{word}
	}
	var out []string
	for start := 0; start < len(word); {
		end := len(word)
		var cur string
		for start < end {
			piece := word[start:end]
			if start > 0 {
				piece = "##" + piece
			}
			if _, ok := t.vocab[piece]; ok {
				cur = piece
				break
			}
			end--
		}
		if cur == "" {
			return []string{"[UNK]"}
		}
		out = append(out, cur)
		start = end
	}
	return out
}

func basicTokens(text string) []string {
	var tokens []string
	var b strings.Builder
	flush := func() {
		if b.Len() > 0 {
			tokens = append(tokens, b.String())
			b.Reset()
		}
	}
	for _, r := range text {
		switch {
		case unicode.IsSpace(r):
			flush()
		case unicode.IsPunct(r):
			flush()
			tokens = append(tokens, string(r))
		default:
			b.WriteRune(unicode.ToLower(r))
		}
	}
	flush()
	return tokens
}

func recognize(ctx context.Context, sess *gonnx.Session, tok *wordPieceTokenizer, labels []string, text string) ([]Entity, error) {
	ids, mask, types, tokens := tok.Encode(text)
	shape := []int64{1, int64(len(ids))}

	inputIDs, err := gonnx.Tensor(sess.Runtime, ids, shape...)
	if err != nil {
		return nil, err
	}
	defer inputIDs.Close()
	attentionMask, err := gonnx.Tensor(sess.Runtime, mask, shape...)
	if err != nil {
		return nil, err
	}
	defer attentionMask.Close()
	tokenTypeIDs, err := gonnx.Tensor(sess.Runtime, types, shape...)
	if err != nil {
		return nil, err
	}
	defer tokenTypeIDs.Close()

	outputs, err := sess.Run(ctx, map[string]*ort.Value{
		"input_ids":      inputIDs,
		"attention_mask": attentionMask,
		"token_type_ids": tokenTypeIDs,
	})
	if err != nil {
		return nil, err
	}
	defer outputs["logits"].Close()
	logits, _, err := gonnx.TensorData[float32](outputs["logits"])
	if err != nil {
		return nil, err
	}

	var ents []Entity
	for i, tok := range tokens {
		if tok == "[CLS]" || tok == "[SEP]" {
			continue
		}
		start := i * len(labels)
		id, score := argmaxSoftmax(logits[start : start+len(labels)])
		label := labels[id]
		if label != "O" {
			ents = append(ents, Entity{Token: tok, Label: label, Score: score})
		}
	}
	sort.SliceStable(ents, func(i, j int) bool { return ents[i].Score > ents[j].Score })
	return ents, nil
}

func argmaxSoftmax(xs []float32) (int, float32) {
	best := 0
	for i := 1; i < len(xs); i++ {
		if xs[i] > xs[best] {
			best = i
		}
	}
	// Good enough for display/tests: return max logit when full softmax is not needed.
	return best, xs[best]
}
