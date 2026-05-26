package gonnx

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	ort "github.com/shota3506/onnxruntime-purego/onnxruntime"
)

// ModelBundle describes an embedded model asset that can be extracted before
// creating an ONNX Runtime session.
type ModelBundle struct {
	// Name is used for the extraction directory. It defaults to the model file base name.
	Name string
	// FS contains ModelRel and ExtraRels.
	FS fs.FS
	// ModelRel is the path within FS to the .onnx model.
	ModelRel string
	// ExtraRels are optional adjacent model assets to extract, for example tokenizer files.
	ExtraRels []string
}

// PrepareModelBundle extracts a model bundle to a stable temp directory and
// returns the extracted model path. OpenBundle calls this automatically; use it
// directly only to prewarm the extraction cache or to pass the extracted path to
// lower-level APIs.
func PrepareModelBundle(bundle ModelBundle) (string, error) {
	if bundle.FS == nil || bundle.ModelRel == "" {
		return "", fmt.Errorf("gonnx: model bundle requires FS and ModelRel")
	}
	name := bundle.Name
	if name == "" {
		name = filepath.Base(bundle.ModelRel)
	}
	dir := filepath.Join(os.TempDir(), "gonnx", "models", sanitizePlatform(name))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	for _, rel := range append([]string{bundle.ModelRel}, bundle.ExtraRels...) {
		if err := extractFile(dir, bundle.FS, rel); err != nil {
			return "", err
		}
	}
	return filepath.Join(dir, filepath.Base(bundle.ModelRel)), nil
}

// Option configures NewRuntime, Open, OpenReader, and OpenBundle.
type Option func(*config)

type config struct {
	runtime    *ort.Runtime
	apiVersion uint32
	logID      string
	logLevel   *ort.LoggingLevel
	options    *ort.SessionOptions
}

// WithRuntime makes a session use an existing ONNX Runtime instance. Sessions
// opened with this option do not close the runtime.
func WithRuntime(runtime *ort.Runtime) Option {
	return func(c *config) { c.runtime = runtime }
}

// WithAPIVersion sets the ONNX Runtime C API version. It defaults to 23.
func WithAPIVersion(version uint32) Option {
	return func(c *config) { c.apiVersion = version }
}

// WithLogID sets the ONNX Runtime environment log ID. It defaults to "gonnx".
func WithLogID(id string) Option {
	return func(c *config) { c.logID = id }
}

// WithLogLevel sets the ONNX Runtime environment log level. It defaults to
// ort.LoggingLevelWarning.
func WithLogLevel(level ort.LoggingLevel) Option {
	return func(c *config) { c.logLevel = &level }
}

// WithSessionOptions passes ONNX Runtime session options to session creation.
func WithSessionOptions(options *ort.SessionOptions) Option {
	return func(c *config) { c.options = options }
}

// WithThreads sets SessionOptions.IntraOpNumThreads.
func WithThreads(threads int) Option {
	return func(c *config) {
		if c.options == nil {
			c.options = &ort.SessionOptions{}
		}
		c.options.IntraOpNumThreads = threads
	}
}

func newConfig(opts []Option) config {
	cfg := config{apiVersion: 23, logID: "gonnx"}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	if cfg.apiVersion == 0 {
		cfg.apiVersion = 23
	}
	if cfg.logID == "" {
		cfg.logID = "gonnx"
	}
	return cfg
}

// Session owns the ONNX Runtime handles commonly needed for one loaded model.
type Session struct {
	Runtime *ort.Runtime
	Env     *ort.Env
	Session *ort.Session

	ownsRuntime bool
}

// Open loads an ONNX model from a filesystem path.
func Open(modelPath string, opts ...Option) (*Session, error) {
	cfg := newConfig(opts)
	runtime, ownRuntime, err := runtimeForSession(cfg)
	if err != nil {
		return nil, err
	}
	s := &Session{Runtime: runtime, ownsRuntime: ownRuntime}
	defer func() {
		if err != nil {
			s.Close()
		}
	}()
	s.Env, err = runtime.NewEnv(cfg.logID, cfg.logLevelOrDefault())
	if err != nil {
		return nil, err
	}
	s.Session, err = runtime.NewSession(s.Env, modelPath, cfg.options)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// OpenReader loads an ONNX model from r.
func OpenReader(r io.Reader, opts ...Option) (*Session, error) {
	cfg := newConfig(opts)
	runtime, ownRuntime, err := runtimeForSession(cfg)
	if err != nil {
		return nil, err
	}
	s := &Session{Runtime: runtime, ownsRuntime: ownRuntime}
	defer func() {
		if err != nil {
			s.Close()
		}
	}()
	s.Env, err = runtime.NewEnv(cfg.logID, cfg.logLevelOrDefault())
	if err != nil {
		return nil, err
	}
	s.Session, err = runtime.NewSessionFromReader(s.Env, r, cfg.options)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// OpenBundle prepares bundle then opens a session for its model.
func OpenBundle(bundle ModelBundle, opts ...Option) (*Session, error) {
	modelPath, err := PrepareModelBundle(bundle)
	if err != nil {
		return nil, err
	}
	return Open(modelPath, opts...)
}

// Close releases session, environment, and runtime handles. It is safe to call multiple times.
func (s *Session) Close() {
	if s == nil {
		return
	}
	if s.Session != nil {
		s.Session.Close()
		s.Session = nil
	}
	if s.Env != nil {
		s.Env.Close()
		s.Env = nil
	}
	if s.ownsRuntime && s.Runtime != nil {
		_ = s.Runtime.Close()
	}
	s.Runtime = nil
	s.ownsRuntime = false
}

// InputNames returns a copy of the loaded model input names.
func (s *Session) InputNames() []string {
	if s == nil || s.Session == nil {
		return nil
	}
	return append([]string(nil), s.Session.InputNames()...)
}

// OutputNames returns a copy of the loaded model output names.
func (s *Session) OutputNames() []string {
	if s == nil || s.Session == nil {
		return nil
	}
	return append([]string(nil), s.Session.OutputNames()...)
}

// Run forwards to the underlying ONNX Runtime session.
func (s *Session) Run(ctx context.Context, inputs map[string]*ort.Value, opts ...ort.RunOption) (map[string]*ort.Value, error) {
	if s == nil || s.Session == nil {
		return nil, fmt.Errorf("gonnx: session is closed")
	}
	return s.Session.Run(ctx, inputs, opts...)
}

// Tensor creates an ONNX Runtime tensor value.
func Tensor[T ort.TensorData](runtime *ort.Runtime, data []T, shape ...int64) (*ort.Value, error) {
	return ort.NewTensorValue(runtime, data, shape)
}

// TensorData returns tensor data and shape.
func TensorData[T ort.TensorData](value *ort.Value) ([]T, []int64, error) {
	return ort.GetTensorData[T](value)
}

func (c config) logLevelOrDefault() ort.LoggingLevel {
	if c.logLevel != nil {
		return *c.logLevel
	}
	return ort.LoggingLevelWarning
}

func runtimeForSession(cfg config) (*ort.Runtime, bool, error) {
	if cfg.runtime != nil {
		return cfg.runtime, false, nil
	}
	runtime, err := NewRuntime(WithAPIVersion(cfg.apiVersion))
	return runtime, true, err
}
