package wordpiece

import (
	"reflect"
	"testing"
)

func TestEncodeOffsetsAndSubwords(t *testing.T) {
	tok := New(map[string]int64{
		"[UNK]": 0, "[CLS]": 1, "[SEP]": 2,
		"bar": 3, "##ack": 4, "obama": 5, ",": 6, "東京": 7,
	}, true)
	text := "Barack Obama, 東京"
	got := tok.Encode(text)
	want := []Token{
		{Text: "[CLS]", ID: 1, Start: -1, End: -1, Special: true},
		{Text: "bar", ID: 3, Start: 0, End: 3},
		{Text: "##ack", ID: 4, Start: 3, End: 6},
		{Text: "obama", ID: 5, Start: 7, End: 12},
		{Text: ",", ID: 6, Start: 12, End: 13},
		{Text: "東京", ID: 7, Start: 14, End: 20},
		{Text: "[SEP]", ID: 2, Start: -1, End: -1, Special: true},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v want %#v", got, want)
	}
	for _, token := range got {
		if !token.Special && text[token.Start:token.End] == "" {
			t.Fatalf("empty original span for %#v", token)
		}
	}
}

func TestEncodeUnknownKeepsWholeWordOffset(t *testing.T) {
	tok := New(map[string]int64{"[UNK]": 9, "[CLS]": 1, "[SEP]": 2}, true)
	got := tok.Encode("Mystery")
	want := []Token{{Text: "[CLS]", ID: 1, Start: -1, End: -1, Special: true}, {Text: "[UNK]", ID: 9, Start: 0, End: 7}, {Text: "[SEP]", ID: 2, Start: -1, End: -1, Special: true}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v want %#v", got, want)
	}
}

func TestEncodeUnicodeLowercaseExpansionOffsets(t *testing.T) {
	tok := New(map[string]int64{"[UNK]": 0, "[CLS]": 1, "[SEP]": 2, "i": 3, "##̇": 4}, true)
	text := "İ"
	got := tok.Encode(text)
	want := []Token{{Text: "[CLS]", ID: 1, Start: -1, End: -1, Special: true}, {Text: "i", ID: 3, Start: 0, End: 2}, {Text: "[SEP]", ID: 2, Start: -1, End: -1, Special: true}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v want %#v", got, want)
	}
}
