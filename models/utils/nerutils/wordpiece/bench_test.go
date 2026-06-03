package wordpiece

import "testing"

func BenchmarkEncode(b *testing.B) {
	tok := New(map[string]int64{
		"[UNK]": 0, "[CLS]": 1, "[SEP]": 2,
		"barack": 3, "obama": 4, "was": 5, "born": 6, "in": 7, "hawaii": 8, ".": 9,
	}, true)
	text := "Barack Obama was born in Hawaii. Barack Obama was born in Hawaii."
	b.ReportAllocs()
	for b.Loop() {
		_ = tok.Encode(text)
	}
}
