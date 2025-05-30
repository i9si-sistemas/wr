package command

import (
	"io"
	"os/exec"
)

// OS implements the Executor interface using the os/exec package.
// It serves as a command builder that allows chaining configurations
// before running the command.
type OS struct {
	Runner
}

// New creates a new Executor instance based on the operating system.
// Useful to start building commands.
func New(r ...Runner) Executor {
	if len(r) > 0 {
		return &OS{r[0]}
	}
	return &OS{&Command{new(exec.Cmd)}}
}

// Convert transforms an existing Runner into an Executor,
// allowing reuse of previously configured commands.
func Convert(r Runner) Executor {
	return &OS{r}
}

// Execute initializes a new command with the given name and arguments.
// It overrides any previously configured command and allows chaining.
func (os *OS) Execute(name string, args ...string) Executor {
	os.Runner = &Command{exec.Command(name, args...)}
	return os
}

// Command wraps exec.Cmd and implements the Runner interface.
// It enables fluent configuration and execution of commands.
type Command struct {
	*exec.Cmd
}

// NewRunner creates a new Runner instance with specified stdout and stderr writers.
func NewRunner(stdout, stderr io.Writer) Runner {
	return &Command{&exec.Cmd{
		Stdout: stdout,
		Stderr: stderr,
	}}
}

// WithDir sets the working directory where the command will be executed.
// Returns the Runner to allow chaining.
func (c *Command) WithDir(dirPath string) Runner {
	c.Dir = dirPath
	return c
}

// WithPath sets the absolute path of the executable to be used.
// Useful when the binary is not in the system PATH.
func (c *Command) WithPath(path string) Runner {
	c.Cmd.Path = path
	return c
}

// Path returns the absolute path of the executable configured in the command.
func (c *Command) Path() string {
	return c.Cmd.Path
}

// CombinedOutput runs the command and returns its combined
// stdout and stderr output. Returns an error if it fails.
func (c *Command) CombinedOutput() ([]byte, error) {
	return c.Cmd.CombinedOutput()
}

// Run executes the command and blocks until it completes.
// Returns an error if the command fails.
func (c *Command) Run() error {
	return c.Cmd.Run()
}

// AppendArgs adds additional arguments to the existing command.
// Useful for dynamically building arguments before execution.
func (c *Command) AppendArgs(args ...string) Runner {
	c.Args = append(c.Args, args...)
	return c
}
