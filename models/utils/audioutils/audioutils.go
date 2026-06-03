// Package audioutils provides small audio helpers for ONNX model adapters.
package audioutils

// PCM16ToFloat32 converts signed 16-bit PCM samples to normalized float32 audio.
func PCM16ToFloat32(pcm []int16) []float32 {
	audio := make([]float32, len(pcm))
	for i, s := range pcm {
		audio[i] = float32(s) / 32768.0
	}
	return audio
}
