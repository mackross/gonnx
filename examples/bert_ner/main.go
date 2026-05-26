package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mackross/gonnx"
	ort "github.com/shota3506/onnxruntime-purego/onnxruntime"
)

func main() {
	modelPath := flag.String("model", "examples/bert_ner/model.onnx", "path to an ONNX export of dslim/bert-base-NER-uncased")
	vocabPath := flag.String("vocab", "examples/bert_ner/vocab.txt", "path to BERT vocab.txt")
	configPath := flag.String("config", "examples/bert_ner/config.json", "path to model config.json")
	text := flag.String("text", "Hugging Face Inc. is based in New York City and Sarah lives in London.", "text to tag")
	runtimePath := flag.String("runtime", os.Getenv("ONNXRUNTIME_LIBRARY"), "optional path to libonnxruntime shared library; defaults to bundled runtime")
	flag.Parse()

	for _, path := range []string{*modelPath, *vocabPath, *configPath} {
		if _, err := os.Stat(path); err != nil {
			log.Fatalf("required file not found: %s", path)
		}
	}

	labels, err := loadLabels(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	tok, err := loadTokenizer(*vocabPath)
	if err != nil {
		log.Fatal(err)
	}

	opts := []gonnx.Option{gonnx.WithLogLevel(ort.LoggingLevelWarning)}
	if *runtimePath != "" {
		rt, err := ort.NewRuntime(*runtimePath, 23)
		if err != nil {
			log.Fatal(err)
		}
		defer rt.Close()
		opts = append(opts, gonnx.WithRuntime(rt))
	}

	model, err := os.Open(*modelPath)
	if err != nil {
		log.Fatal(err)
	}
	defer model.Close()

	sess, err := gonnx.OpenReader(model, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer sess.Close()

	entities, err := recognize(context.Background(), sess, tok, labels, *text)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("text: %s\n", *text)
	for _, e := range entities {
		fmt.Printf("%s\t%s\t%.3f\n", e.Token, e.Label, e.Score)
	}
}
