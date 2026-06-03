// Package nerutils contains shared implementation helpers for bundled NER models.
//
// Model packages should wrap or alias these types so callers do not need to
// import utils packages directly.
package nerutils

import "context"

// Recognizer recognizes named entities in text.
type Recognizer interface {
	// Recognize returns named entities found in text. Implementations return byte
	// offsets into the original input.
	Recognize(ctx context.Context, text string) ([]Entity, error)
}

// Entity is a named entity span in the original input text.
type Entity struct {
	Text  string
	Label string
	Start int
	End   int
	Score float32
}

// ModelConfig contains model-opening settings.
type ModelConfig struct {
	// Threads is the requested inference thread count. Zero means model default.
	Threads int
}
