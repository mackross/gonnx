package silero_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/mackross/gonnx"
	"github.com/mackross/gonnx/models/silero"
	_ "github.com/mackross/gonnx/runtimes/linuxamd64"
)

func TestModelInference(t *testing.T) {
	session, err := silero.OpenSession(gonnx.WithThreads(1))
	if err != nil {
		t.Fatal(err)
	}
	if got, want := session.InputNames(), []string{"input", "state", "sr"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("input names = %v, want %v", got, want)
	}
	if got, want := session.OutputNames(), []string{"output", "stateN"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("output names = %v, want %v", got, want)
	}
	session.Close()

	vad, err := silero.Open(gonnx.WithThreads(1))
	if err != nil {
		t.Fatal(err)
	}
	defer vad.Close()

	probability, _, err := vad.ProbabilityPCM16(context.Background(), make([]int16, silero.ChunkSamples), silero.SampleRate)
	if err != nil {
		t.Fatal(err)
	}
	if probability < 0 || probability > 1 {
		t.Fatalf("speech probability out of range: %v", probability)
	}
}
