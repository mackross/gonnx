package chunk

import "strings"

// Token is an input token prediction for BIO chunking.
type Token struct {
	Text    string
	Start   int
	End     int
	Label   string
	Score   float32
	Special bool
}

// Entity is a chunked entity span.
type Entity struct {
	Text  string
	Label string
	Start int
	End   int
	Score float32
}

type active struct {
	label string
	start int
	end   int
	sum   float32
	count int
}

// BIO chunks BIO-tagged token predictions. Orphan I-X starts a new X entity;
// I-Y while X is active closes X and starts Y.
func BIO(tokens []Token, original string) ([]Entity, error) {
	var out []Entity
	var cur *active
	flush := func() {
		if cur == nil {
			return
		}
		text := ""
		if cur.start >= 0 && cur.end >= cur.start && cur.end <= len(original) {
			text = original[cur.start:cur.end]
		}
		out = append(out, Entity{Text: text, Label: cur.label, Start: cur.start, End: cur.end, Score: cur.sum / float32(cur.count)})
		cur = nil
	}
	startEntity := func(label string, t Token) {
		cur = &active{label: label, start: t.Start, end: t.End, sum: t.Score, count: 1}
	}
	for _, t := range tokens {
		if t.Special || t.Start < 0 || t.End <= t.Start || t.End > len(original) {
			continue
		}
		prefix, label := split(t.Label)
		switch prefix {
		case "O", "":
			flush()
		case "B":
			flush()
			startEntity(label, t)
		case "I":
			if cur == nil || cur.label != label {
				flush()
				startEntity(label, t)
			} else {
				cur.end = t.End
				cur.sum += t.Score
				cur.count++
			}
		default:
			flush()
		}
	}
	flush()
	return out, nil
}

func split(label string) (prefix, typ string) {
	if label == "" || label == "O" {
		return label, ""
	}
	p, rest, ok := strings.Cut(label, "-")
	if !ok || rest == "" {
		return "B", label
	}
	return p, rest
}
