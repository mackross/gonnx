# NER model benchmarks

This document records a benchmark smoke result for the bundled named-entity
recognition models in `github.com/mackross/gonnx/models`.

Run the benchmarks from the models module:

```sh
cd models
go test -run '^$' -bench BenchmarkModelInference -benchmem ./bertuncased ./bertcased ./distilbertcased ./distilbertuncased ./multidistilbert ./neurobert
```

The model inference tests and benchmarks use the same direct package style that
callers use:

```go
recognizer, err := bertuncased.Open(gonnx.WithThreads(1))
```

Recent benchmark smoke result on an AMD Ryzen AI 9 HX 370 with
`GOMAXPROCS=18`, `WithThreads(1)`, and a 10-WordPiece input:

| Model package | ns/op | tokens/sec |
| --- | ---: | ---: |
| `neurobert.Open(...)` | 265,955 | 58,966 |
| `distilbertuncased.Open(...)` | 4,143,474 | 2,920 |
| `distilbertcased.Open(...)` | 3,846,573 | 2,872 |
| `multidistilbert.Open(...)` | 5,914,874 | 2,696 |
| `bertcased.Open(...)` | 8,221,905 | 1,378 |
| `bertuncased.Open(...)` | 8,407,055 | 1,317 |

Treat these numbers as relative guidance only. Hardware, operating system,
ONNX Runtime build, thread count, input length, model cache state, and Go version
all affect results.

For a quick size/performance rule of thumb:

- use `neurobert` when binary size and latency matter most;
- use `distilbertcased` or `distilbertuncased` for a balanced English model;
- use `bertcased` or `bertuncased` when BERT-base accuracy is worth the larger
  binary and slower CPU inference;
- use `multidistilbert` when multilingual coverage matters more than size.
