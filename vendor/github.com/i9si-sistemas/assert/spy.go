package assert

import (
	"bytes"
	"io"
	"log"
	"os"
)

type params struct {
	Args   []any
	Format string
}

type spy struct {
	failed,
	skipped bool
	calls     map[string]bool
	params    map[string]params
	setOutput func(io.Writer)
	t         T
}

func newSpy(t T) *spy {
	return &spy{
		calls:     make(map[string]bool),
		params:    make(map[string]params),
		t:         t,
		setOutput: log.SetOutput,
	}
}

func (s *spy) Cleanup(fn func()) {
	s.calls["Cleanup"] = true
	fn()
}

func (s *spy) Fail() {
	s.calls["Fail"] = true
	s.failed = true
}

func (s *spy) Failed() bool {
	s.calls["Failed"] = true
	return s.failed
}

func (s *spy) Fatal(args ...any) {
	name := "Fatal"
	s.calls[name] = true
	s.params[name] = params{Args: args}
}

func (s *spy) Helper() {
	s.calls["Helper"] = true
}

func (s *spy) Setenv(key, value string) {
	name := "Setenv"
	s.calls[name] = true
	s.params[name] = params{Args: []any{key, value}}
}

func (s *spy) Skip(args ...any) {
	name := "Skip"
	s.calls[name] = true
	s.params[name] = params{Args: args}
	s.skipped = true
}

func (s *spy) Skipped() bool {
	s.calls["Skipped"] = true
	return s.skipped
}

func (s *spy) TempDir() string {
	s.calls["TempDir"] = true
	return "./spy"
}

func (s *spy) LogOutput() (buf *bytes.Buffer, cancel func()) {
	buf = new(bytes.Buffer)
	s.setOutput(buf)
	cancel = func() {
		s.setOutput(os.Stdout)
	}
	return
}
