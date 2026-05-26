package gonnx

import (
	"testing"

	ort "github.com/shota3506/onnxruntime-purego/onnxruntime"
)

func TestDefaultLogLevel(t *testing.T) {
	if got := newConfig(nil).logLevelOrDefault(); got != ort.LoggingLevelWarning {
		t.Fatalf("default log level = %v, want %v", got, ort.LoggingLevelWarning)
	}
}

func TestWithLogLevelVerbose(t *testing.T) {
	if got := newConfig([]Option{WithLogLevel(ort.LoggingLevelVerbose)}).logLevelOrDefault(); got != ort.LoggingLevelVerbose {
		t.Fatalf("explicit verbose log level = %v, want %v", got, ort.LoggingLevelVerbose)
	}
}
