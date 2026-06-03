package chunk

import "testing"

func BenchmarkBIO(b *testing.B) {
	text := "Barack Obama was born in Hawaii."
	tokens := []Token{{Start: 0, End: 6, Label: "B-PER", Score: .9}, {Start: 7, End: 12, Label: "I-PER", Score: .8}, {Start: 13, End: 16, Label: "O", Score: 1}, {Start: 17, End: 21, Label: "O", Score: 1}, {Start: 25, End: 31, Label: "B-LOC", Score: .7}}
	b.ReportAllocs()
	for b.Loop() {
		_, _ = BIO(tokens, text)
	}
}
