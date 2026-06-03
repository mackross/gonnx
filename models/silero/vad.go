package silero

import (
	"context"
	"fmt"
	"time"

	"github.com/mackross/gonnx"
	ort "github.com/shota3506/onnxruntime-purego/onnxruntime"
)

const (
	// SampleRate is the PCM sample rate expected by the Silero VAD model.
	SampleRate = 16000
	// ChunkSamples is the exact number of PCM samples accepted per inference.
	ChunkSamples = 512
)

// VAD runs the Silero voice activity detector with recurrent state.
type VAD struct {
	session *gonnx.Session
	state   []float32
	context []float32
}

// Open opens the bundled Silero VAD model.
func Open(opts ...gonnx.Option) (*VAD, error) {
	sess, err := OpenSession(opts...)
	if err != nil {
		return nil, err
	}
	return &VAD{session: sess, state: make([]float32, 2*1*128), context: make([]float32, 64)}, nil
}

// Close releases ONNX Runtime resources.
func (v *VAD) Close() {
	if v != nil && v.session != nil {
		v.session.Close()
		v.session = nil
	}
}

// Reset clears recurrent model state.
func (v *VAD) Reset() {
	if v == nil {
		return
	}
	clear(v.state)
	clear(v.context)
}

// ProbabilityPCM16 returns speech probability for one 512-sample mono 16 kHz chunk.
func (v *VAD) ProbabilityPCM16(ctx context.Context, chunk []int16, sampleRate int) (float32, time.Duration, error) {
	if v == nil || v.session == nil {
		return 0, 0, fmt.Errorf("silero: VAD is closed")
	}
	if sampleRate != SampleRate {
		return 0, 0, fmt.Errorf("silero: expected %d Hz, got %d", SampleRate, sampleRate)
	}
	if len(chunk) != ChunkSamples {
		return 0, 0, fmt.Errorf("silero: expected %d samples, got %d", ChunkSamples, len(chunk))
	}
	inputData := make([]float32, 64+ChunkSamples)
	copy(inputData, v.context)
	for i, s := range chunk {
		inputData[64+i] = float32(s) / 32768.0
	}
	input, err := gonnx.Tensor(v.session.Runtime, inputData, 1, int64(len(inputData)))
	if err != nil {
		return 0, 0, err
	}
	defer input.Close()
	state, err := gonnx.Tensor(v.session.Runtime, v.state, 2, 1, 128)
	if err != nil {
		return 0, 0, err
	}
	defer state.Close()
	sr, err := gonnx.Tensor(v.session.Runtime, []int64{int64(sampleRate)})
	if err != nil {
		return 0, 0, err
	}
	defer sr.Close()
	start := time.Now()
	outputs, err := v.session.Run(ctx, map[string]*ort.Value{"input": input, "state": state, "sr": sr})
	if err != nil {
		return 0, 0, err
	}
	dur := time.Since(start)
	out := outputs["output"]
	stateOut := outputs["stateN"]
	defer out.Close()
	defer stateOut.Close()
	prob, _, err := gonnx.TensorData[float32](out)
	if err != nil {
		return 0, 0, err
	}
	newState, _, err := gonnx.TensorData[float32](stateOut)
	if err != nil {
		return 0, 0, err
	}
	copy(v.state, newState)
	copy(v.context, inputData[len(inputData)-64:])
	if len(prob) == 0 {
		return 0, 0, fmt.Errorf("silero: empty output")
	}
	return prob[0], dur, nil
}
