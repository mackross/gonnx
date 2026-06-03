# gonnx

`gonnx` is a small Go package for shipping and using ONNX Runtime from Go
without cgo in your application code. It builds on
[`onnxruntime-purego`](https://github.com/shota3506/onnxruntime-purego), which loads libonnx using purego (an alternative to cgo - amazing), and adds
packaging helpers for applications that want a batteries included fast ONNX Runtime that embeds in the Go
binary. While there are other non-cgo, non purego options, I found the inference to be orders
of magnitudes slower. This library also has packaging helpers for easy embedding of models 
into your binary.

Benefits:

- platform-specific ONNX Runtime bundles for Linux, macOS, and Windows
- blank-import runtime packages, so binaries only include the platforms you choose
- extraction of embedded shared libraries to a SHA-256 checked temp cache
- simple `Open`, `OpenReader`, and `OpenBundle` helpers for models
- small tensor helpers for common input/output handling
- optional bundled model packages under `github.com/mackross/gonnx/models`

## Quick start

```go
import (
    "github.com/mackross/gonnx"
    _ "github.com/mackross/gonnx/runtimes/linuxamd64"
)

sess, err := gonnx.Open("model.onnx", gonnx.WithThreads(1))
if err != nil {
    return err
}
defer sess.Close()
```

Importing a runtime package embeds that ONNX Runtime build in your binary. At
runtime, `gonnx` extracts the registered shared libraries into a checked cache
under the system temp directory, then loads ONNX Runtime via
`onnxruntime-purego`.

If you do not want bundled runtime assets, use `onnxruntime-purego` directly or
pass your own runtime with `gonnx.WithRuntime(rt)`.

## Bundled runtimes

The repository includes ONNX Runtime `1.23.x` assets, matching the supported
version range of `onnxruntime-purego`.

Available runtime packages:

- `github.com/mackross/gonnx/runtimes/linuxamd64`
- `github.com/mackross/gonnx/runtimes/linuxarm64`
- `github.com/mackross/gonnx/runtimes/darwinamd64`
- `github.com/mackross/gonnx/runtimes/darwinarm64`
- `github.com/mackross/gonnx/runtimes/windowsamd64`
- `github.com/mackross/gonnx/runtimes/windowsarm64`

## Bundled models

Bundled models live in the nested module `github.com/mackross/gonnx/models`.
They are optional: importing the root `gonnx` package does not embed model
weights. Import only the model packages you want to ship, plus one runtime
package for your target platform.

Each model package exposes a high-level `Open(opts ...gonnx.Option)` helper for
normal use. Voice packages also expose `OpenSession` for raw ONNX Runtime access.

| Package | Task | Embedded assets | Public entry point | Notes |
| --- | --- | ---: | --- | --- |
| `models/silero` | Voice activity detection | ~2 MiB ONNX | `silero.Open(...)` | Stateful Silero VAD for 16 kHz PCM. Accepts 512-sample chunks and returns speech probability. |
| `models/smartturn` | Voice turn-completion detection | ~25 MiB ONNX | `smartturn.Open(...)` | Smart Turn v3.2 CPU model. Converts the last 8 seconds of 16 kHz PCM to log-mel features and predicts whether the user turn is complete. |
| `models/neurobert` | Named entity recognition | ~9.3 MiB ONNX + vocab/config | `neurobert.Open(...)` | Tiny English NER model. Fastest/smallest bundled NER option, with lower accuracy tradeoff and richer OntoNotes-like labels such as `PERSON`, `ORG`, `GPE`, `DATE`. |
| `models/distilbertuncased` | Named entity recognition | ~64 MiB ONNX + vocab/config | `distilbertuncased.Open(...)` | English uncased DistilBERT NER. Smaller/faster than BERT-base; ignores case. |
| `models/distilbertcased` | Named entity recognition | ~63 MiB ONNX + vocab/config | `distilbertcased.Open(...)` | English cased DistilBERT NER. Smaller/faster than BERT-base while preserving casing signal. |
| `models/bertuncased` | Named entity recognition | ~105 MiB ONNX + vocab/config | `bertuncased.Open(...)` | English uncased BERT-base NER with CoNLL labels: `PER`, `ORG`, `LOC`, `MISC`. Good default when binary size is acceptable. |
| `models/bertcased` | Named entity recognition | ~104 MiB ONNX + vocab/config | `bertcased.Open(...)` | English cased BERT-base NER with CoNLL labels. Prefer when casing carries useful signal. |
| `models/multidistilbert` | Named entity recognition | ~130 MiB ONNX + vocab/config | `multidistilbert.Open(...)` | Multilingual cased DistilBERT NER for high-resource languages. Larger binary; useful when language coverage matters. |

See [`models/NER_BENCHMARK.md`](models/NER_BENCHMARK.md) for NER model
benchmark notes and recent relative throughput numbers.

Example NER use:

```go
import (
    "context"

    "github.com/mackross/gonnx"
    "github.com/mackross/gonnx/models/bertuncased"
    _ "github.com/mackross/gonnx/runtimes/linuxamd64"
)

recognizer, err := bertuncased.Open(gonnx.WithThreads(1))
if err != nil {
    return err
}
defer recognizer.Close()

entities, err := recognizer.Recognize(context.Background(), "Barack Obama worked at Microsoft in Seattle.")
```

Example voice use:

```go
import (
    "context"

    "github.com/mackross/gonnx"
    "github.com/mackross/gonnx/models/silero"
    _ "github.com/mackross/gonnx/runtimes/linuxamd64"
)

vad, err := silero.Open(gonnx.WithThreads(1))
if err != nil {
    return err
}
defer vad.Close()

probability, _, err := vad.ProbabilityPCM16(context.Background(), chunk, silero.SampleRate)
```

## Options

`Open`, `OpenReader`, `OpenBundle`, and `NewRuntime` use the same option style:

```go
sess, err := gonnx.Open("model.onnx",
    gonnx.WithThreads(1),
    gonnx.WithLogLevel(onnxruntime.LoggingLevelWarning),
)
```

Useful options:

- `WithRuntime(rt)` uses an existing `*onnxruntime.Runtime`
- `WithAPIVersion(version)` overrides the default ONNX Runtime C API version, currently `23`
- `WithLogID(id)` sets the ONNX Runtime environment log ID
- `WithLogLevel(level)` sets the ONNX Runtime log level, including verbose
- `WithSessionOptions(options)` passes raw `onnxruntime.SessionOptions`
- `WithThreads(n)` sets `SessionOptions.IntraOpNumThreads`

## Running inference

```go
input, err := gonnx.Tensor(sess.Runtime, []float32{1, 2, 3, 4}, 1, 4)
if err != nil {
    return err
}
defer input.Close()

outputs, err := sess.Run(ctx, map[string]*onnxruntime.Value{
    sess.InputNames()[0]: input,
})
if err != nil {
    return err
}
defer outputs[sess.OutputNames()[0]].Close()

data, shape, err := gonnx.TensorData[float32](outputs[sess.OutputNames()[0]])
```

## Embedded model bundles

Models and sidecar files can use the same extract-and-cache pattern as runtime
libraries:

```go
//go:embed models/model.onnx models/vocab.txt
var modelFS embed.FS

sess, err := gonnx.OpenBundle(gonnx.ModelBundle{
    Name:      "my-model",
    FS:        modelFS,
    ModelRel:  "models/model.onnx",
    ExtraRels: []string{"models/vocab.txt"},
}, gonnx.WithThreads(1))
```

`OpenBundle` automatically prepares the model bundle. `PrepareModelBundle` is
available if you want to prewarm the extraction cache or pass the extracted path
to lower-level APIs.

## Example

See [`examples/bert_ner`](examples/bert_ner) for a live named-entity recognition
example using an ONNX export of
[`dslim/bert-base-NER-uncased`](https://huggingface.co/dslim/bert-base-NER-uncased).

Run the live test with:

```bash
GONNX_LIVE_BERT_NER=1 go test -v ./examples/bert_ner -run TestLiveBertNERRecognizesEntities
```

## Testing

```bash
go test ./...
cd models && go test ./...
```
