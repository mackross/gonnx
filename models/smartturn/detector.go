package smartturn

import (
	"context"
	"fmt"
	"time"

	"github.com/mackross/gonnx"
	ort "github.com/shota3506/onnxruntime-purego/onnxruntime"
)

// Result is a Smart Turn completion prediction.
type Result struct {
	Complete    bool
	Probability float32
	Duration    time.Duration
}

// Detector runs the Smart Turn model using a gonnx session.
type Detector struct {
	session   *gonnx.Session
	extractor *FeatureExtractor
}

// Open opens the bundled Smart Turn detector.
func Open(opts ...gonnx.Option) (*Detector, error) {
	sess, err := OpenSession(opts...)
	if err != nil {
		return nil, err
	}
	d := &Detector{session: sess, extractor: NewFeatureExtractor()}
	if got := d.session.InputNames(); len(got) != 1 || got[0] != "input_features" {
		d.Close()
		return nil, fmt.Errorf("smartturn: unexpected inputs: %v", got)
	}
	if got := d.session.OutputNames(); len(got) != 1 || got[0] != "logits" {
		d.Close()
		return nil, fmt.Errorf("smartturn: unexpected outputs: %v", got)
	}
	return d, nil
}

// Close releases ONNX Runtime resources.
func (d *Detector) Close() {
	if d != nil && d.session != nil {
		d.session.Close()
		d.session = nil
	}
}

// PredictPCM16 predicts turn completion from mono 16 kHz signed PCM.
func (d *Detector) PredictPCM16(ctx context.Context, pcm []int16, sampleRate int) (Result, error) {
	if sampleRate != SampleRate {
		return Result{}, fmt.Errorf("smartturn: expected %d Hz audio, got %d", SampleRate, sampleRate)
	}
	audio := make([]float32, len(pcm))
	for i, s := range pcm {
		audio[i] = float32(s) / 32768.0
	}
	return d.PredictFloat32(ctx, audio)
}

// PredictFloat32 predicts turn completion from normalized mono 16 kHz PCM.
func (d *Detector) PredictFloat32(ctx context.Context, audio []float32) (Result, error) {
	return d.PredictFeatures(ctx, d.extractor.Extract(audio))
}

// PredictFeatures predicts turn completion from flattened [80,800] log-mel features.
func (d *Detector) PredictFeatures(ctx context.Context, features []float32) (Result, error) {
	if d == nil || d.session == nil {
		return Result{}, fmt.Errorf("smartturn: detector is closed")
	}
	if len(features) != FeatureSize*NumFrames {
		return Result{}, fmt.Errorf("smartturn: expected %d features, got %d", FeatureSize*NumFrames, len(features))
	}
	input, err := gonnx.Tensor(d.session.Runtime, features, 1, FeatureSize, NumFrames)
	if err != nil {
		return Result{}, err
	}
	defer input.Close()
	start := time.Now()
	outputs, err := d.session.Run(ctx, map[string]*ort.Value{"input_features": input}, ort.WithOutputNames("logits"))
	if err != nil {
		return Result{}, err
	}
	out := outputs["logits"]
	defer out.Close()
	data, _, err := gonnx.TensorData[float32](out)
	if err != nil {
		return Result{}, err
	}
	if len(data) != 1 {
		return Result{}, fmt.Errorf("smartturn: expected one output logit/probability, got %d", len(data))
	}
	p := data[0]
	return Result{Complete: p > 0.5, Probability: p, Duration: time.Since(start)}, nil
}
