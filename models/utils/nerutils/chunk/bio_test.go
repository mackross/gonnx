package chunk

import (
	"reflect"
	"testing"
)

func TestBIOValidationCases(t *testing.T) {
	text := "Barack Obama IBM Alice"
	cases := []struct {
		name string
		in   []Token
		want []Entity
	}{
		{
			name: "B PER I PER becomes one entity",
			in: []Token{
				{Start: 0, End: 6, Label: "B-PER", Score: 0.8},
				{Start: 7, End: 12, Label: "I-PER", Score: 1.0},
			},
			want: []Entity{{Text: "Barack Obama", Label: "PER", Start: 0, End: 12, Score: 0.9}},
		},
		{
			name: "B PER B PER becomes two entities",
			in: []Token{
				{Start: 0, End: 6, Label: "B-PER", Score: 0.8},
				{Start: 7, End: 12, Label: "B-PER", Score: 1.0},
			},
			want: []Entity{{Text: "Barack", Label: "PER", Start: 0, End: 6, Score: 0.8}, {Text: "Obama", Label: "PER", Start: 7, End: 12, Score: 1.0}},
		},
		{
			name: "orphan I starts deterministic entity",
			in:   []Token{{Start: 0, End: 6, Label: "I-PER", Score: 0.7}},
			want: []Entity{{Text: "Barack", Label: "PER", Start: 0, End: 6, Score: 0.7}},
		},
		{
			name: "mismatched I closes and starts",
			in: []Token{
				{Start: 0, End: 6, Label: "B-ORG", Score: 0.8},
				{Start: 7, End: 12, Label: "I-PER", Score: 0.9},
			},
			want: []Entity{{Text: "Barack", Label: "ORG", Start: 0, End: 6, Score: 0.8}, {Text: "Obama", Label: "PER", Start: 7, End: 12, Score: 0.9}},
		},
		{
			name: "O terminates active entity",
			in: []Token{
				{Start: 0, End: 6, Label: "B-PER", Score: 0.8},
				{Start: 7, End: 12, Label: "O", Score: 1.0},
				{Start: 13, End: 16, Label: "B-ORG", Score: 0.6},
			},
			want: []Entity{{Text: "Barack", Label: "PER", Start: 0, End: 6, Score: 0.8}, {Text: "IBM", Label: "ORG", Start: 13, End: 16, Score: 0.6}},
		},
		{
			name: "special padding ignored",
			in: []Token{
				{Start: 0, End: 6, Label: "B-PER", Score: 0.8, Special: true},
				{Start: 17, End: 22, Label: "B-PER", Score: 0.9},
				{Start: 22, End: 22, Label: "I-PER", Score: 1.0},
			},
			want: []Entity{{Text: "Alice", Label: "PER", Start: 17, End: 22, Score: 0.9}},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := BIO(tc.in, text)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got %#v want %#v", got, tc.want)
			}
		})
	}
}
