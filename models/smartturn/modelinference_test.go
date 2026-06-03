package smartturn_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/mackross/gonnx"
	"github.com/mackross/gonnx/models/smartturn"
	_ "github.com/mackross/gonnx/runtimes/linuxamd64"
)

func TestModelInference(t *testing.T) {
	session, err := smartturn.OpenSession(gonnx.WithThreads(1))
	if err != nil {
		t.Fatal(err)
	}
	if got, want := session.InputNames(), []string{"input_features"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("input names = %v, want %v", got, want)
	}
	if got, want := session.OutputNames(), []string{"logits"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("output names = %v, want %v", got, want)
	}
	session.Close()

	detector, err := smartturn.Open(gonnx.WithThreads(1))
	if err != nil {
		t.Fatal(err)
	}
	defer detector.Close()

	result, err := detector.PredictPCM16(context.Background(), make([]int16, smartturn.NSamples), smartturn.SampleRate)
	if err != nil {
		t.Fatal(err)
	}
	if result.Probability < 0 || result.Probability > 1 {
		t.Fatalf("turn probability out of range: %v", result.Probability)
	}
}
