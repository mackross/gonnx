package wordpiece

import (
	"bufio"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Token is one encoded token with byte offsets into the original text.
type Token struct {
	Text    string
	ID      int64
	Start   int
	End     int
	Special bool
}

// Tokenizer is a small WordPiece tokenizer for BERT-style NER models.
type Tokenizer struct {
	vocab map[string]int64
	unk   int64
	cls   int64
	sep   int64
	lower bool
}

// Load reads a newline-delimited WordPiece vocabulary.
func Load(r io.Reader, lower bool) (*Tokenizer, error) {
	s := bufio.NewScanner(r)
	vocab := map[string]int64{}
	var id int64
	for s.Scan() {
		vocab[s.Text()] = id
		id++
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return New(vocab, lower), nil
}

// New creates a tokenizer from vocab.
func New(vocab map[string]int64, lower bool) *Tokenizer {
	cp := make(map[string]int64, len(vocab))
	for k, v := range vocab {
		cp[k] = v
	}
	return &Tokenizer{vocab: cp, unk: cp["[UNK]"], cls: cp["[CLS]"], sep: cp["[SEP]"], lower: lower}
}

// Encode tokenizes text, including [CLS] and [SEP].
func (t *Tokenizer) Encode(text string) []Token {
	out := []Token{{Text: "[CLS]", ID: t.cls, Start: -1, End: -1, Special: true}}
	for _, bt := range basic(text) {
		out = append(out, t.pieces(text, bt)...)
	}
	out = append(out, Token{Text: "[SEP]", ID: t.sep, Start: -1, End: -1, Special: true})
	return out
}

type basicToken struct{ start, end int }

func basic(text string) []basicToken {
	var out []basicToken
	start := -1
	flush := func(end int) {
		if start >= 0 {
			out = append(out, basicToken{start, end})
			start = -1
		}
	}
	for i, r := range text {
		if unicode.IsSpace(r) {
			flush(i)
			continue
		}
		if unicode.IsPunct(r) {
			flush(i)
			out = append(out, basicToken{i, i + utf8.RuneLen(r)})
			continue
		}
		if start < 0 {
			start = i
		}
	}
	flush(len(text))
	return out
}

func (t *Tokenizer) pieces(original string, bt basicToken) []Token {
	word := original[bt.start:bt.end]
	lookup, offsets := normalize(word, t.lower)
	if id, ok := t.vocab[lookup]; ok {
		return []Token{{Text: lookup, ID: id, Start: bt.start, End: bt.end}}
	}
	var out []Token
	for start := 0; start < len(lookup); {
		end := len(lookup)
		var cur string
		var curStart, curEnd int
		for start < end {
			piece := lookup[start:end]
			if start > 0 {
				piece = "##" + piece
			}
			if _, ok := t.vocab[piece]; ok {
				cur = piece
				curStart, curEnd = start, end
				break
			}
			_, size := utf8.DecodeLastRuneInString(lookup[:end])
			if size == 0 {
				end--
			} else {
				end -= size
			}
		}
		if cur == "" {
			return []Token{{Text: "[UNK]", ID: t.unk, Start: bt.start, End: bt.end}}
		}
		out = append(out, Token{Text: cur, ID: t.vocab[cur], Start: bt.start + offsets[curStart], End: bt.start + offsets[curEnd]})
		start = curEnd
	}
	return out
}

func normalize(word string, lower bool) (string, []int) {
	if !lower {
		offsets := make([]int, len(word)+1)
		for i := range offsets {
			offsets[i] = i
		}
		return word, offsets
	}
	var b strings.Builder
	offsets := []int{0}
	for i, r := range word {
		lowered := string(unicode.ToLower(r))
		b.WriteString(lowered)
		end := i + utf8.RuneLen(r)
		for range len(lowered) {
			offsets = append(offsets, end)
		}
	}
	return b.String(), offsets
}
